package utils

import (
	"os"
	"os/exec"
	"fmt"
)

func RunCmd(c *exec.Cmd, msg string) {

	cmd := c
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
    PrintErrorExit(msg, err)
  }

	fmt.Println()

}

func RunCmdQuiet(c *exec.Cmd, msg string) {
	
	cmd := c
	if err := cmd.Run(); err != nil {
		PrintErrorExit(msg, err)
	}

}

func ChkIfCmdRuns(c *exec.Cmd) error {

	err := c.Run()
	return err

}
