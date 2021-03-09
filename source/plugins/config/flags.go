package config

import "flag"

var (
	configsPath = flag.String("config", "/dev/null", "Path to configs file")
	mode = Mode{
		value: flag.String("mode", "cli", "Application mode (i.e. cli, tg)"),
	}
)
