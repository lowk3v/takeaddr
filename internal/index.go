package internal

import (
	global "github.com/lowk3v/micro-tool-template/config"
	"github.com/lowk3v/micro-tool-template/internal/enum"
)

type Options struct {
	Action  enum.ACTION
	Verbose bool
	Version bool
	Output  string
}

func New(opt *Options) *Options {
	return opt
}

func (opt *Options) Run() {
	// TODO
	if opt.Verbose {
		global.Log.SetLevel("debug")
	}

	switch opt.Action {
	case enum.SHOWVERSION:
		global.Log.Infof("Version: %s", global.Version)
		break
	default:
		global.Log.Errorf("Action not found")
	}
}
