package interfaces

import "models"

type (
	Controller interface {
		Run(processing models.ProcessingChannels)
	}
)
