package models

type (
	ProcessingChannels struct {
		Done  chan bool
		Error chan error
	}
)
