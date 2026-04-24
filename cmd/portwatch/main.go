package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

var version = "dev"

func main() {
	var (
		cfgPath  = flag.String("config", "", "path to config file (optional)")
		printVer = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *printVer {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	cfgFile, err := config.FindConfigFile(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error locating config: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	scanner, err := ports.NewScanner()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating scanner: %v\n", err)
		os.Exit(1)
	}

	ruleSet, err := rules.Compile(cfg.Rules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error compiling rules: %v\n", err)
		os.Exit(1)
	}

	notifier := alert.NewLogNotifier(cfg.LogLevel)
	dispatcher := alert.NewDispatcher(ruleSet, notifier)

	mon := monitor.New(scanner, dispatcher, cfg.Interval)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Printf("portwatch %s starting (interval: %s)\n", version, cfg.Interval)

	if err := mon.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "monitor exited with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("portwatch stopped")
}
