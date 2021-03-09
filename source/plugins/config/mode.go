package config

import "flag"

type Mode struct {
	value *string
}

var mode = Mode{
	value: flag.String("mode", "cli", "Application mode (i.e. cli, tg)"),
}

func (m Mode) IsCLI() bool {
	return *m.value == "cli"
}

func (m Mode) IsTelegram() bool {
	return *m.value == "tg"
}

func (m Mode) Value() string {
	return *m.value
}
