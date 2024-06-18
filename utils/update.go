package utils

import (
	"os"
	"os/exec"
)

func UpdateApp() {
	
	download := exec.Command("curl", "-L", "-o", "/tmp/updateapp.sh", "https://hkdb.github.io/app/updateapp.sh")
	RunCmd(download, "Failed to download update script...")

	/* 
	err := Copy("$HOME/Development/app/dist/updateapp.sh", "/tmp/updateapp.sh")
	if err != nil {
		PrintErrorExit("Copy Error:", err)
	}
	*/

	update := exec.Command("bash", "-c", "/tmp/updateapp.sh")
	os.Chmod("/tmp/updateapp.sh", 0711)
	RunCmd(update, "Failed to update app...")

}
