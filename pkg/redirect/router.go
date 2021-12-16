package redirect

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	routePath = "/{scheme}/{host}{path:.*}"
)

// router redirects incoming HTTP requests based on configured redirects
type router struct {
	*mux.Router
	redirectHandlers map[string]http.Handler
	redirectLock     sync.RWMutex
}

// Router returns a router that can add or remove redirects
func Router() *router {
	r := &router{
		redirectHandlers: make(map[string]http.Handler),
	}
	r.Router = mux.NewRouter()
	r.HandleFunc(routePath, func(rw http.ResponseWriter, req *http.Request) {
		// Figure out target address
		var err error
		vars := mux.Vars(req)
		address := fmt.Sprintf("%s://%s", vars["scheme"], vars["host"])
		addressWithPath := fmt.Sprintf("%s%s", address, vars["path"])
		req.URL, err = url.Parse(addressWithPath)
		if err != nil {
			http.Error(rw, fmt.Sprintf("could not parse URL from %s", addressWithPath), http.StatusBadRequest)
			return
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		// ensure that the handler for a redirect does not change while processing a request
		r.redirectLock.RLock()
		defer r.redirectLock.RUnlock()
		// grab the redirect handler
		handler, ok := r.redirectHandlers[address]
		if !ok {
			http.Error(rw, fmt.Sprintf("redirect address %s has not been registered", address), http.StatusBadRequest)
			return
		}
		// pass request and response to redirect for processing
		handler.ServeHTTP(rw, req)
	})
	return r
}

// RegisterHandler configures a redirect to the provided address
func (r *router) RegisterHandler(address string, redirect Redirect) error {
	if _, ok := r.redirectHandlers[address]; ok {
		return fmt.Errorf("cannot register multiple redirects for address %s", address)
	}

	// add handler
	r.redirectLock.Lock()
	r.redirectHandlers[address] = redirect.ToHandler()
	r.redirectLock.Unlock()

	// configure watcher for upgrading handler on file changes
	go func() {
		w, err := redirect.RestartWatcher()
		if err != nil {
			logrus.Errorf("unable to set up watcher for redirect %s", redirect)
		}
		defer w.Close()
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					r.redirectLock.Lock()
					r.redirectHandlers[address] = redirect.ToHandler()
					r.redirectLock.Unlock()
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()
	return nil
}
