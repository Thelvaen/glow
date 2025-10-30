package main

import (
	"flag"
	"glow/config"
	"glow/ui"
	"os"

	"k8s.io/klog/v2"
)

var (
	cfgPath string
)

func main() {
	defer klog.Flush()

	flag.StringVar(&cfgPath, "config", "", "Path to the config file")
	flag.Parse()

	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	cfg := config.Load(cfgPath)
	ui.Run(cfg)
}
