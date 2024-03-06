package tools

import (
	"github.com/rs/zerolog/log"
	"math/rand"
	"strconv"
	"strings"
)

func RandomCode(len int) int {
	lowStr := "1" + strings.Repeat("0", len-1)
	highStr := strings.Repeat("9", len)
	low, err := strconv.Atoi(lowStr)
	if err != nil {
		log.Error().Err(err).Msg("RandomCode error, parse string")
	}
	high, err := strconv.Atoi(highStr)
	if err != nil {
		log.Error().Err(err).Msg("RandomCode error, parse string")
	}
	return low + rand.Intn(high-low)
}
