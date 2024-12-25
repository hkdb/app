package redhat

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var sudo = [3]string{"/usr/bin/sudo", "/bin/sh", "-c"}
var cmd = "/usr/bin/dnf"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "dnf", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" is already installed...", "")
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

	install := exec.Command(sudo[0], sudo[1], sudo[2], command+pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "dnf", pname)
	if rerr != nil {
		utils.PrintErrorExit("Record Error:", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "dnf", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == false {
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
	}

	action := " remove "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	remove := exec.Command(sudo[0], sudo[1], sudo[2], command+pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "dnf", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error:", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -a remove -p " + pkg + "...")

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

	err := syscall.Exec("/usr/bin/rpm", []string{"/usr/bin/rpm", "-qa"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func ListSystemSearch(pkg string) {

	err := syscall.Exec("/usr/bin/rpm", []string{"/usr/bin/rpm", "-qa", "|", "grep", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Search Error:", err)
	}

}

func Update() {

	action := " check-update"
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	update := exec.Command(sudo[0], sudo[1], sudo[2], command)

	update.Stdin = os.Stdin
	update.Stdout = os.Stdout
	update.Stderr = os.Stderr
	update.Env = os.Environ()
	if err := update.Run(); err != nil && err.ExitCode() != 100 {
		utils.PrintErrorExit(msg, err)
	}

	fmt.Println()

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

	fmt.Println("This is an apt only command. Just use app -a upgrade -p...")

}

func Search(pkg string) {

	err := syscall.Exec(cmd, []string{cmd, "search", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error:", err)
	}

}

func InstallAll() {

	// apt
	fmt.Println("DNF:\n")
	pkgs, aperr := db.ReadPkgs("", "packages", "dnf")
	if aperr != nil {
		utils.PrintErrorExit("DNF - Read ERROR:", aperr)
	}
	action := " install "
	command := cmd + action
	if env.AutoYes == true {
		command = cmd + " -y" + action
	}

	install := exec.Command(sudo[0], sudo[1], sudo[2], command+pkgs)
	utils.RunCmd(install, "Installation Error:")

}
