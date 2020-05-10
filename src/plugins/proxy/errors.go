package proxy

import "errors"

var (
	noAvailableIps = errors.New("unable to find available proxy ip")
)
