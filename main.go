package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/lowk3v/takeaddr/internal"
	"github.com/lowk3v/takeaddr/internal/enum"
	"github.com/lowk3v/takeaddr/utils"
	"os"
)
import global "github.com/lowk3v/takeaddr/config"

func __existArg(arg string) bool {
	args := os.Args[1:]
	return len(args) > 0 && args[0] == arg
}

type ArgList []string

func (a *ArgList) String() string {
	return ""
}

func (a *ArgList) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func _banner() {
	// https://patorjk.com/software/taag/#p=display&f=ANSI%20Shadow&t=%20dumpsc
	_, _ = fmt.Fprintf(os.Stderr, "%s %s by %s\n%s\nCredits: https://github.com/%s/%s\nTwitter: https://twitter.com/%s\n\n",
		color.HiBlueString(`    
	████████╗ █████╗ ██╗  ██╗███████╗ █████╗ ██████╗ ██████╗ ██████╗ 
    ╚══██╔══╝██╔══██╗██║ ██╔╝██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔══██╗
       ██║   ███████║█████╔╝ █████╗  ███████║██║  ██║██║  ██║██████╔╝
       ██║   ██╔══██║██╔═██╗ ██╔══╝  ██╔══██║██║  ██║██║  ██║██╔══██╗
       ██║   ██║  ██║██║  ██╗███████╗██║  ██║██████╔╝██████╔╝██║  ██║
       ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═╝  ╚═╝
                                                                     `),
		color.BlueString("v"+global.Version),
		color.BlueString("@LowK"),
		"Get all smart contract address from urls",
		"LowK",
		"takeaddr",
		"LowK3v_",
	)
	_, _ = fmt.Fprintf(os.Stderr, "Usage of: %s <options> <args>\n", os.Args[0])
	flag.PrintDefaults()
}

func _parseFlags() (*internal.Options, error) {
	var configPath string
	var verbose bool
	var version bool
	var output string
	var resolve bool
	var complete bool
	var noColors bool
	var insecure bool

	// global configurations
	options := &internal.Options{
		Action:     enum.NONE,
		Verbose:    verbose,
		Version:    version,
		Output:     output,
		Url:        "",
		Method:     "GET",
		OutputFile: "",
		InputFile:  "",
		Resolve:    false,
		Complete:   false,
		NoColors:   false,
		Headers:    []string{},
		Insecure:   true,
		Timeout:    10,
	}

	// global arguments
	flag.StringVar(&configPath, "c", "", "optional. Path to config file")
	flag.BoolVar(&version, "version", false, "print version and exit")
	// module configurations, implement if needed
	flag.StringVar(&options.Url, "u", "", "url to crawl")
	flag.StringVar(&options.Method, "m", "GET", "http method")
	flag.StringVar(&options.OutputFile, "o", "", "output file")
	flag.StringVar(&options.InputFile, "i", "", "input file")
	flag.BoolVar(&resolve, "resolve", true, "Resolve the output and filter out the non existing files (Can only be used in combination with --complete)")
	flag.BoolVar(&complete, "complete", true, "Complete the urls. e.g. /js/index.js -> https://example.com/js/index.js")
	flag.BoolVar(&noColors, "no-colors", false, "no colors")
	flag.BoolVar(&insecure, "insecure", true, "insecure")
	flag.IntVar(&options.Timeout, "timeout", 10, "timeout")
	flag.Var((*ArgList)(&options.Headers), "header", "http header")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()

	if version {
		options.Action = enum.SHOWVERSION
		return options, nil
	}

	if verbose {
		options.Verbose = true
	}

	if noColors {
		options.NoColors = true
	}

	if resolve {
		options.Resolve = true
	}

	if complete {
		options.Complete = true
	}

	if insecure {
		options.Insecure = true
	}

	if options.Url != "" {
		options.Action = enum.URL
	} else {
		options.Action = enum.NONE
	}

	// Custom config file
	if len(configPath) > 0 {
		if err := utils.FileExists(configPath, false); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return options, err
		}
		if err := global.CustomConfig(configPath); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return options, err
		}
	}

	// Action is required
	if options.Action == enum.NONE {
		_banner()
		os.Exit(0)
	}

	return options, nil
}

func main() {
	options, err := _parseFlags()
	if err != nil {
		os.Exit(0)
	}
	addresses, err := internal.Run(*options)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	internal.Explorer(addresses)
}
