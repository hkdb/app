package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func RunCmdReturn(c *exec.Cmd) (string, error) {
	
	cmd := c
	var outbuf strings.Builder
	cmd.Stdout = &outbuf
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return "", err
	}
	stdout := outbuf.String()
	return stdout, nil

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
