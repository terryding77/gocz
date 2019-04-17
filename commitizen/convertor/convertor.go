package convertor

import "github.com/terryding77/gocz/commitizen/config"

// Convertor is common interface of all specific Convertor
type Convertor interface {
	Setup(config.Config) bool
	Convert(msg string) interface{}
}
