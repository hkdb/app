package arch

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var cmd_paru = "/usr/bin/paru"

func ParuInstall(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "paru", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" is already installed...", "")
	}

	install := exec.Command(cmd_paru, "-S", pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "paru", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func ParuRemove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "paru", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == false {
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
	}

	remove := exec.Command(cmd_paru, "-R", pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "paru", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func ParuPurge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m paru remove " + pkg + "...")

}

func ParuAutoRemove() {

	out, err := exec.Command("/bin/bash", "-c", cmd_paru+" -Qtdq").Output()
	rmList := string(out)
	if err != nil && rmList != "" {
		utils.PrintErrorExit("Read Auto Remove Package List Error:", err)
	}

	if rmList == "" {
		fmt.Println("No packages need to be automatically removed...\n")
		os.Exit(1)
	}
	aRemove := exec.Command(cmd_paru, "-Rns", rmList)
	utils.RunCmd(aRemove, "Auto Remove Error:")

}

func ParuListSystem() {

	err := syscall.Exec(cmd_paru, []string{cmd_paru, "-Q"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func ParuListSystemSearch(pkg string) {

	listSys := exec.Command("/bin/sh", "-c", cmd_paru+" -Q |grep "+pkg)
	utils.RunCmd(listSys, "List paru Packages Error:")

}

func ParuUpdate() {

	upgrade := exec.Command(cmd_paru, "-Syy")
	utils.RunCmd(upgrade, "Update Error")

}

func ParuUpgrade() {

	upgrade := exec.Command(cmd_paru, "-Syyu")
	utils.RunCmd(upgrade, "Upgrade Error")

}

func ParuDistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m paru upgrade...")

}

func ParuSearch(pkg string) {

	err := syscall.Exec(cmd_paru, []string{cmd_paru, "-Ss", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error:", err)
	}

}

func ParuInstallAll() {

	// paru
	fmt.Println("paru:\n")
	pkgs, aperr := db.ReadPkgs("", "packages", "paru")
	if aperr != nil {
		utils.PrintErrorExit("paru - Read ERROR:", aperr)
		os.Exit(1)
	}
	installAll := exec.Command(cmd_paru, "-S", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}
