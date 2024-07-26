package http_proxy

import "math"

type IRetryOptions interface {
	GetNextTimer(retries int) int
}

type LinearRetry struct {
	IRetryOptions
	Milliseconds int
}

func (options *LinearRetry) GetNextTimer(retries int) int {
	return options.Milliseconds * retries
}

type ExponentialRetry struct {
	IRetryOptions
	Base           float64
	InitialDelayMS int
	MaxDelay       int
}

func (options *ExponentialRetry) GetNextTimer(retries int) int {
	delay := options.InitialDelayMS * int(math.Pow(options.Base, float64(retries)))
	if options.MaxDelay > 0 && delay > options.MaxDelay {
		return options.MaxDelay
	}
	return delay
}
