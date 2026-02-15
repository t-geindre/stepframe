package midi

import "errors"

var ErrPortNotFound = errors.New("midi port not found")
var ErrUnknownCommand = errors.New("unknown command")
var ErrAlreadyRunning = errors.New("already running")
