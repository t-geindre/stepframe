package midi

import (
	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

func CloseDriver(logger zerolog.Logger) {
	midi.CloseDriver()
	logger.Info().Str("component", "midi driver").Msg("midi driver closed")
}
