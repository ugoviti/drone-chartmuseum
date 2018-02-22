package main

import (
	"fmt"
	"io/ioutil"
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
			Value:  "./",
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

func initAction(c *cli.Context) config {
	var conf config
	conf.RepoURL = c.String("repo-url")
	conf.ChartDir = c.String("chart-dir")
	conf.ChartPath = c.String("chart-path")
	conf.PreviousCommitID = c.String("previous-commit")
	conf.CurrentCommitID = c.String("current-commit")
	conf.SaveDir = c.String("save-dir")

	return conf
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
	conf := initAction(c)
	dirs, err := ioutil.ReadDir(conf.ChartDir)
	var resultList []string
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
		chart, err := saveChartToPackage(dir.Name(), conf.SaveDir)
		if err == nil {
			resultList = append(resultList, chart)
		}
	}

	uploadToServer(resultList, conf.RepoURL)
	return nil
}

func diffMode(c *cli.Context) error {

	conf := initAction(c)
	files := getDiffFiles(conf.ChartDir, conf.PreviousCommitID, conf.CurrentCommitID)

	files = getUniqueParentFolders(filterExtFiles(files))
	if len(files) == 0 {
		fmt.Print("No chart needs to be updated! Exit ... ")
		os.Exit(0)
	}
	var resultList []string
	for _, file := range files {
		chart, err := saveChartToPackage(file, conf.SaveDir)
		if err == nil {
			resultList = append(resultList, chart)
		}
	}
	uploadToServer(resultList, conf.RepoURL)
	return nil
}

func singleMode(c *cli.Context) error {
	conf := initAction(c)

	var resultList []string
	chart, err := saveChartToPackage(conf.ChartPath, conf.SaveDir)
	if err == nil {
		resultList = append(resultList, chart)
	}
	uploadToServer(resultList, conf.RepoURL)
	return nil
}

func main() {
	app := initApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
