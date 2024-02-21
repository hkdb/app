package debian

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"
	"github.com/hkdb/app/env"

	"syscall"
	"os"
  "os/exec"
	"fmt"
)

var sudo = [3]string{"/usr/bin/sudo", "/bin/sh", "-c"}
var cmd = "/usr/bin/apt"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "apt", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	pname := pkg
	local, name, pfile := utils.IsLocalInstall(pkg)
	if local == true {
		pname = name
		pkg = pfile
	}

	action := " install "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	install := exec.Command(sudo[0], sudo[1], sudo[2], command + pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "apt", pname)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "apt", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}
	
	action := " remove "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	remove := exec.Command(sudo[0], sudo[1], sudo[2], command + pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "apt", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

  // Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "apt", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}

	action := " purge "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	purge := exec.Command(sudo[0], sudo[1], sudo[2], command + pkg)
	utils.RunCmd(purge, "Purge Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "apt", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func AutoRemove() {

	action := " autoremove"
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	err := syscall.Exec(sudo[0], []string{sudo[0], sudo[1], sudo[2], command}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Auto Remove Error:", err)
	}

}

func ListSystem() {
	
	err := syscall.Exec("/usr/bin/dpkg", []string{"/usr/bin/dpkg", "-l"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func ListSystemSearch(pkg string) {

	err := syscall.Exec("/bin/bash", []string{"/bin/bash", "-c", "/usr/bin/dpkg -l |grep " + pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Search Error:", err)
	}

}

func Update() {

	action := " update"
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	update := exec.Command(sudo[0], sudo[1], sudo[2], command)
	utils.RunCmd(update, "Update Error:")

}

func Upgrade() {

	action := " upgrade"
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}
	
	upgrade := exec.Command(sudo[0], sudo[1], sudo[2], command)
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	action := " dist-upgrade"
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	distUpgrade := exec.Command(sudo[0], sudo[1], sudo[2], command)
	utils.RunCmd(distUpgrade, "Dist Upgrade Error:")

}

func Search(pkg string) {

	err := syscall.Exec("/usr/bin/apt-cache", []string{"apt-cache", "search", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error:", err)
	}

}

func InstallAll() {
	
	// apt
	fmt.Println("APT:\n")
	pkgs, aperr := db.ReadPkgs("", "packages", "apt")
	if aperr != nil {
		utils.PrintErrorExit("APT - Read ERROR:", aperr)
	}
	action := " install "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	install := exec.Command(sudo[0], sudo[1], sudo[2], command + pkgs)
	utils.RunCmd(install, "Installation Error:")

}
