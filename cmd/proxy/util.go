package proxy

import "fmt"

type port uint16
type address string

func allHosts(p port) address {
	return address(fmt.Sprintf(":%d", p))
}

func localhost(p port) address {
	return address(fmt.Sprintf("localhost:%d", p))
}
