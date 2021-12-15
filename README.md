portexporter
========

A tunnel server and reverse proxy client that are powered by [rancher/remotedialer](https://github.com/rancher/remotedialer).

The client is placed in the hostNetwork of a host and makes an outbound connection with the tunnel server.

The tunnel server is placed in its own network and can now do `net.Dial` on the client and pipe all bytes back and forth.

The client can also initialize multiple `httputil.ReverseProxy` that `net.Dial` requests will be forwarded to; this allows HTTP requests proxied via the tunnel server to optionally utilize client certificates to access a HTTPS service on the hostNetwork.

## Building

`make`

## Running

`./bin/portexporter`

## License
Copyright (c) 2019 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
