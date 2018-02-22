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
			Usage:  "chart path (required if mode is single)",
			Value:  "",
			EnvVar: "PLUGIN_CHART_PATH",
		},
		cli.StringFlag{
			Name:   "chart-dir",
			Value:  "",
			Usage:  "chart directory (required if mode is diff or all)",
			EnvVar: "PLUGIN_CHART_DIR",
		},
		cli.StringFlag{
			Name:   "save-dir",
			Value:  "uploads/",
			Usage:  "directory to save chart packages",
			EnvVar: "PLUGIN_SAVE_DIR",
		},
		cli.StringFlag{
			Name:   "previous-commit",
			Usage:  "previous commit id (`COMMIT_SHA`, required if mode is diff)",
			EnvVar: "PLUGIN_PREVIOUS_COMMIT",
		},
		cli.StringFlag{
			Name:   "current-commit",
			Usage:  "current commit id (`COMMIT_SHA`, required if mode is diff)",
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
	repoURL := c.String("repo-url")
	chartDir := c.String("chart-dir")
	previousCommitID := c.String("previous-commit")
	commitID := c.String("current-commit")
	saveDir := c.String("save-dir")
	files := getDiffFiles(chartDir, previousCommitID, commitID)

	files = getUniqueParentFolders(filterExtFiles(files))
	var resultList []string
	for _, file := range files {
		resultList = append(resultList, saveChartToPackage(file, saveDir))
	}
	fmt.Println(resultList)
	uploadToServer(resultList, repoURL)
	fmt.Println("how about this")
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
