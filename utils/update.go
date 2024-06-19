package utils

import (
	"github.com/hkdb/app/env"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func UpdateApp() {
	
	url := "https://hkdb.github.io/app/version.txt"
	
	resp, err := http.Get(url)
	if err != nil {
		PrintErrorExit("Site Access Error:", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		PrintErrorExit("Site Data Read Error:", err)
	}

	latest := string(data)
	latest = strings.TrimSuffix(latest, "\n")
	latest = latest[1:]
	current := env.Version
	current = current[1:]

	latest_num, err := strconv.ParseFloat(latest, 64)
	if err != nil {
		PrintErrorExit("Latest Parse Error:", err)
	}
	current_num, err := strconv.ParseFloat(current, 64)
	if err != nil {
		PrintErrorExit("Current Parse Error:", err)
	}


	if current_num >= latest_num {
		fmt.Println("")
		fmt.Println("✅️ app on your system is already the latest release... Nothing to do...")
		fmt.Println("")
		os.Exit(0)
	}

	download := exec.Command("curl", "-L", "-o", "/tmp/updateapp.sh", "https://hkdb.github.io/app/updateapp.sh")
	RunCmd(download, "Failed to download update script...")

	/*
	// For testing
	err := Copy("$HOME/Development/app/dist/updateapp.sh", "/tmp/updateapp.sh")
	if err != nil {
		PrintErrorExit("Copy Error:", err)
	}
	*/

	update := exec.Command("bash", "-c", "/tmp/updateapp.sh")
	os.Chmod("/tmp/updateapp.sh", 0711)
	RunCmd(update, "Failed to update app...")

}
