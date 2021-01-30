package util

import (
	"time"

	"github.com/ory/viper"
)

// NewTimeOutContext is timeout duration
func NewTimeOutContext() time.Duration {
	timeoutContext := time.Duration(viper.GetInt("contextTimeout")) * time.Second
	return timeoutContext
}
