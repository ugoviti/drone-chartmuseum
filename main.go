package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "drone-chartmuseum-plugin"
	app.Usage = "drone plugin to upload charts to chartmuseum server"
	app.Version = fmt.Sprintf("1.0.0")

	mainFlag := []cli.Flag{
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key `AWS_ACCESS_KEY`",
			EnvVar: "PLUGIN_ACCESS_KEY,AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Key `AWS_SECRET_KEY`",
			EnvVar: "PLUGIN_SECRET_KEY,AWS_SECRET_ACCESS_KEY,AWS_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "aws-region",
			Value:  "ap-southeast-1",
			Usage:  "AWS Region `AWS_REGION`",
			EnvVar: "PLUGIN_REGION, AWS_REGION",
		},
		cli.StringFlag{
			Name:   "repo-url",
			Value:  "",
			Usage:  "chartmuseum server endpoint",
			EnvVar: "PLUGIN_REPO_URL",
		},
		cli.StringFlag{
			Name:   "mode",
			Value:  "",
			Usage:  "which mode to run (all|diff|single)",
			EnvVar: "PLUGIN_MODE",
		},
		cli.StringFlag{
			Name:   "chart-path",
			Usage:  "chart path (required if `mode` is `single`)",
			Value:  "",
			EnvVar: "PLUGIN_CHART_PATH",
		},
		cli.StringFlag{
			Name:   "chart-dir",
			Value:  "",
			Usage:  "chart directory (required if `mode` is `diff` or `all`)",
			EnvVar: "PLUGIN_CHART_DIR",
		},
		cli.StringFlag{
			Name:   "previous-commit",
			Usage:  "previous commit id (required if `mode` is `diff`)",
			EnvVar: "PLUGIN_PREVIOUS_COMMIT",
		},
		cli.StringFlag{
			Name:   "current-commit",
			Usage:  "current commit id (required if `mode` is `diff`)",
			EnvVar: "PLUGIN_CURRENT_COMMIT",
		},
	}

	app.Action = cli.ActionFunc(defaultAction)
	app.Flags = mainFlag

	return app
}

func defaultAction(c *cli.Context) error {
	action := c.String("mode")
	switch action {
	case "all":
		allMode(c)
	case "diff":
		diffMode(c)
	case "single":
		singleMode(c)
	default:
		log.Fatal("mode not valid!")
	}
	return nil
}

func allMode(c *cli.Context) error {
	return nil
}

func diffMode(c *cli.Context) error {
	return nil
}

func singleMode(c *cli.Context) error {
	return nil
}

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
