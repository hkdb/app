package main

import (
	"github.com/hkdb/app/cli"
	"github.com/hkdb/app/utils"
	"github.com/hkdb/app/mgr"
)

func main() {

	utils.LogLaunchBanner("v0.01")
	
	// Detect environment
	cli.GetEnv()

	// Get user input based on flags:
	flags := cli.ParseFlags()
	
	// Process user input
	mgr.Process(flags)

}
