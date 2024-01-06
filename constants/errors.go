package constants

import "errors"

var (
	ErrUserNotInVoiceChannel = errors.New("couldn't find voice channel with user in it")
)
