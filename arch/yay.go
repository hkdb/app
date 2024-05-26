package arch

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"syscall"
	"os"
	"os/exec"
	"fmt"
)

var cmd_yay = "/usr/bin/yay"

func YayInstall(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "yay", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	install := exec.Command(cmd_yay, "-S", pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "yay", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func YayRemove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "yay", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}

	remove := exec.Command(cmd_yay, "-R", pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")	
	derr := db.RemovePkg("", "packages", "yay", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func YayPurge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m yay remove " + pkg + "...")

}

func YayAutoRemove() {

	out, err := exec.Command("/bin/bash", "-c", cmd_yay + " -Qtdq").Output()
	rmList := string(out)
	if err != nil  && rmList != "" {
		utils.PrintErrorExit("Read Auto Remove Package List Error:", err)
	}
	
	if rmList == "" {
		fmt.Println("No packages need to be automatically removed...\n")
		os.Exit(1)
	}
	aRemove := exec.Command(cmd_yay, "-Rns", rmList)
	utils.RunCmd(aRemove, "Auto Remove Error:")

}

func YayListSystem() {
	
	err := syscall.Exec(cmd_yay, []string{cmd_yay, "-Q"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func YayListSystemSearch(pkg string) {

	listSys := exec.Command("/bin/sh", "-c", cmd_yay + " -Q |grep " + pkg)
	utils.RunCmd(listSys, "List Yay Packages Error:")

}

func YayUpdate() {

	upgrade := exec.Command(cmd_yay, "-Syy")
	utils.RunCmd(upgrade, "Update Error")

}

func YayUpgrade() {

	upgrade := exec.Command(cmd_yay, "-Syyu")
	utils.RunCmd(upgrade, "Upgrade Error")

}

func YayDistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m yay upgrade...")

}

func YaySearch(pkg string) {

	err := syscall.Exec(cmd_yay, []string{cmd_yay, "-Ss", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error:", err)
	}

}

func YayInstallAll() {
	
	// yay
	fmt.Println("YAY:\n")
	pkgs, aperr := db.ReadPkgs("", "packages", "yay")
	if aperr != nil {
		utils.PrintErrorExit("yay - Read ERROR:", aperr)
		os.Exit(1)
	}
	installAll := exec.Command(cmd_yay, "-S", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}
