//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"fmt"
	"os"

	"github.com/aiyengar2/portexporter/cmd/gateway"
	"github.com/aiyengar2/portexporter/cmd/proxy"
	"github.com/aiyengar2/portexporter/cmd/redirector"
	"github.com/aiyengar2/portexporter/cmd/test"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Version    = "v0.0.0-dev"
	GitCommit  = "HEAD"
	KubeConfig string

	// go:embed description.txt
	AppDescription string
)

func main() {
	app := cli.NewApp()
	app.Name = "portexporter"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Usage = "A tunnel server and reverse proxy that are powered by rancher/remotedialer"
	app.Description = AppDescription
	app.CommandNotFound = func(cliCtx *cli.Context, s string) {
		fmt.Fprintf(cliCtx.App.Writer, "Invalid Command: %s \n\n", s)
		if pcliCtx := cliCtx.Parent(); pcliCtx == nil {
			cli.ShowAppHelpAndExit(cliCtx, 1)
		} else {
			cli.ShowCommandHelpAndExit(cliCtx, pcliCtx.Command.Name, 1)
		}
	}
	app.OnUsageError = func(cliCtx *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(cliCtx.App.Writer, "Incorrect Usage: %s \n\n", err.Error())
		if isSubcommand {
			cli.ShowSubcommandHelp(cliCtx)
		} else {
			cli.ShowAppHelp(cliCtx)
		}
		return nil
	}
	app.Before = func(cliCtx *cli.Context) error {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
		logrus.SetOutput(cliCtx.App.Writer)
		return nil
	}

	app.Commands = []cli.Command{
		gateway.NewCommand(),
		proxy.NewCommand(),
		redirector.NewCommand(),
		test.NewCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
