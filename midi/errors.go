package midi

import "errors"

var ErrUnknownEventType = errors.New("Unknown MIDI event type")
var ErrPortNotFound = errors.New("MIDI port not found")
var ErrUnknownCommand = errors.New("Unknown command")
