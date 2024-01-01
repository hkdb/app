package arch

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"syscall"
	"os"
	"os/exec"
	"fmt"
)

func YayInstall(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "yay", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	install := exec.Command("/usr/bin/yay", "-S", pkg)
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

	remove := exec.Command("/usr/bin/yay", "-R", pkg)
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

	out, err := exec.Command("/bin/bash", "-c", "/usr/bin/yay -Qtdq").Output()
	rmList := string(out)
	if err != nil  && rmList != "" {
		utils.PrintErrorExit("Read Auto Remove Package List Error:", err)
	}
	
	if rmList == "" {
		fmt.Println("No packages need to be automatically removed...\n")
		os.Exit(1)
	}
	aRemove := exec.Command("/usr/bin/yay", "-Rns", rmList)
	utils.RunCmd(aRemove, "Auto Remove Error:")

}

func YayListSystem() {
	
	err := syscall.Exec("/usr/bin/yay", []string{"/usr/bin/yay", "-Q"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func YayListSystemSearch(pkg string) {

	listSys := exec.Command("/bin/sh", "-c", "/usr/bin/yay -Q |grep " + pkg)
	utils.RunCmd(listSys, "List Yay Packages Error:")

}

func YayUpdate() {

	fmt.Println("This is an apt only command. Just use app -m yay upgrade...")

}

func YayUpgrade() {

	upgrade := exec.Command("/usr/bin/yay", "-Syyu")
	utils.RunCmd(upgrade, "Upgrade Error")

}

func YayDistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m yay upgrade...")

}

func YaySearch(pkg string) {

	err := syscall.Exec("/usr/bin/yay", []string{"yay", "-Ss", pkg}, os.Environ())
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
	installAll := exec.Command("/usr/bin/yay", "-S", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}
