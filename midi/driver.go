package midi

import (
	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

func Open(logger zerolog.Logger) {
	logger = logger.With().Str("component", "midi_driver").Logger()

	inPorts := midi.GetInPorts()
	for i, port := range inPorts {
		logger.Info().Int("port", i).Str("name", port.String()).Msg("midi in found")
	}

	outPorts := midi.GetOutPorts()
	for i, port := range outPorts {
		logger.Info().Int("port", i).Str("name", port.String()).Msg("midi out found")
	}
}

func Close(logger zerolog.Logger) {
	midi.CloseDriver()
	logger.Info().Str("component", "midi_driver").Msg("midi driver closed")
}
