package models

import (
	"github.com/phuslu/log"
)

type Log struct {
	Level          log.Level
	HttpStatusCode int
	Error          error
	Message        string
}
