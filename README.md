portexporter
========

An HTTP(S) Proxy that registers one or more Gateways across a network boundary and proxies requests to those Gateways, powered by [rancher/remotedialer](https://github.com/rancher/remotedialer).

It is expected that the network on which a Gateway is deployed is capable of starting an **outbound** connection with the Proxy; however, the Gateway itself can be behind a network boundary.

The Gateway also supports a `/proxy/{address}` HTTPS endpoint that will interpret incoming HTTPS requests, parse the provided address, and see if a proxy has been configured for it. If so, it will forward the request using a pre-configured `http.Transport` (optionally containing a `tls.Config`). On forwarding a request, it can also modify the request in flight (e.g. adding or removing headers, such as Bearer Auth tokens) if necessary.

### Use Cases

The primary use case for this project is to enable querying the loopback address networks of multiple hosts behind a network boundary via a proxy.

An example of such a use case would be to enable Prometheus to securely scrape metrics from hosts whose ports are cordoned off by network firewalls.

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
