package utils

import "github.com/lowk3v/micro-tool-template/config"

func HandleError(err error, customMsg string) bool {
	if err != nil {
		config.Log.
			WithField("msg", customMsg).
			Error(err)
		return true
	}
	return false
}
