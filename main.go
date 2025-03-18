package main

import (
	"github.com/hkdb/app/cli"
	"github.com/hkdb/app/mgr"
	"github.com/hkdb/app/utils"
)

func main() {

	utils.LogLaunchBanner("v0.24")

	// Detect environment
	cli.GetEnv()

	// Get user input based on flags:
	flags := cli.ParseFlags()

	// Process user input
	mgr.Process(flags)

}
