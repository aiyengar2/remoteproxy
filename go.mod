module github.com/aiyengar2/portexporter

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/rancher/remotedialer v0.2.6-0.20201012155453-8b1b7bb7d05f
	github.com/rancher/wrangler v0.8.0
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	inet.af/tcpproxy v0.0.0-20200125044825-b6bb9b5b8252
)
