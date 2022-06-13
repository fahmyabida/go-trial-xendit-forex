package utils

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func LogIN(traceId string, s ...string) {
	logging(traceId, "IN", s...)
}

func LogOUT(traceId string, s ...string) {
	logging(traceId, "OUT", s...)
}

func logging(traceId, flow string, s ...string) {
	log.Printf("[%v][%v][%v]", traceId, flow, strings.Join(s, "]["))
}
