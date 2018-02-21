package main

import (
	"os"

	"github.com/mbtproject/mbt/lib"
	"github.com/sirupsen/logrus"
)

func buildStageCB(a *lib.Module, s lib.BuildStage) {
	switch s {
	case lib.BuildStageBeforeBuild:
		logrus.Infof("BUILD %s in %s for %s", a.Name(), a.Path(), a.Version())
	case lib.BuildStageSkipBuild:
		logrus.Infof("SKIP %s in %s for %s", a.Name(), a.Path(), a.Version())
	}
}

func main() {
	var in string
	var debug bool
	debug = true
	if in == "" {
		cwd, _ := os.Getwd()
		in = cwd
	}

	level := lib.LogLevelNormal
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		level = lib.LogLevelDebug
	}

	system, _ := lib.NewSystem(in, level)

	system.BuildDiff("d3dd7708b1ae54f5790e172322b59b8e882b9b44", "b4330ef6c5c8a8c706f0eaa9b20d37288319da80", os.Stdin, os.Stdout, os.Stderr, buildStageCB)
	// fmt.Println(system)
}
