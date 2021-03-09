package interfaces

import "habr-bot/source/models"

type (
	Controller interface {
		Run(processing models.ProcessingChannels)
	}
)
