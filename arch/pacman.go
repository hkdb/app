package arch

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var sudo = [4]string{"/usr/bin/sudo", "-S", "/bin/sh", "-c"}
var cmd = "/usr/bin/pacman"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "pacman", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" is already installed...", "")
	}

	install := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], cmd+" -S "+pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "pacman", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "pacman", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == false {
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
	}

	remove := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], cmd+" -R "+pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "pacman", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app remove " + pkg + "...")

}

func AutoRemove() {

	out, err := exec.Command("/bin/bash", "-c", cmd+" -Qtdq").Output()
	rmList := string(out)
	if err != nil && rmList != "" {
		utils.PrintErrorExit("Read Auto Remove Package List Error:", err)
	}

	if rmList == "" {
		fmt.Println("No packages need to be automatically removed...\n")
		os.Exit(1)
	}
	aRemove := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], "/usr/bin/pacman -Rns "+rmList)
	utils.RunCmd(aRemove, "Auto Remove Error:")

}

func ListSystem() {

	err := syscall.Exec("/usr/bin/pacman", []string{cmd, "-Q"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func ListSystemSearch(pkg string) {

	listSys := exec.Command("/bin/sh", "-c", cmd+" -Q |grep "+pkg)
	utils.RunCmd(listSys, "List System Packages Error:")

}

func Update() {

	action := " -Syy"
	command := cmd + action

	update := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], command)
	utils.RunCmd(update, "Update Error:")

}

func Upgrade() {

	switch env.Distro {
	case "garuda":
		upgrade := exec.Command("/usr/bin/garuda-update")
		utils.RunCmd(upgrade, "Upgrade Error:")
	default:
		upgrade := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], cmd+" -Syyu")
		utils.RunCmd(upgrade, "Upgrade Error:")
	}

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app upgrade...")

}

func Search(pkg string) {

	err := syscall.Exec("/usr/bin/pacman", []string{cmd, "-Ss", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error:", err)
	}

}

func InstallAll() {

	// pacman
	fmt.Println("PACMAN:\n")
	pkgs, aperr := db.ReadPkgs("", "packages", "pacman")
	if aperr != nil {
		utils.PrintErrorExit("PACMAN - Read ERROR:", aperr)
	}

	command := cmd + " -S "
	install := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], command+pkgs)
	utils.RunCmd(install, "Installation Error:")

}
