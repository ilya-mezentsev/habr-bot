package config

type Mode struct {
	value *string
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
