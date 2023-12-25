package main

import (
	"github.com/FirinKinuo/capyback/cli"
	"github.com/FirinKinuo/configpath"
)

var (
	version           = "unversioned"
	defaultConfigFile = "config.yml"
)

func main() {
	configPath := configpath.ConfigPath{
		Application: "capyback",
	}

	userConfigPath := configPath.UserFile(defaultConfigFile)

	capybackCli := cli.NewCapyback(version, userConfigPath)
	_ = capybackCli.Execute()
}
