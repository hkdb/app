package macos

import (
	"strings"

	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"
	"github.com/hkdb/app/env"

	"fmt"
	"os"
	"os/exec"
)

var mgr = env.BrewCmd

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "brew", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" is already installed...", "")
	}

	install := exec.Command(mgr, "install", pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "brew", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "brew", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == false {
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
	}

	remove := exec.Command(mgr, "uninstall", pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "brew", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m brew remove " + pkg + "...")

}

func AutoRemove() {

	fmt.Println("This action does not apply to Homebrew...")

}

func ListSystem() {

	list := exec.Command(mgr, "list")
	utils.RunCmd(list, "List Package Error:")

}

func ListSystemSearch(pkg string) {

	listSearch := exec.Command(mgr, "list", "|", "grep", pkg)
	utils.RunCmd(listSearch, "List Package Search Error:")
}

func Update() {

	update := exec.Command(mgr, "update")
	utils.RunCmd(update, "Update Error:")

}

func Upgrade() {

	upgrade := exec.Command(mgr, "update")
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m brew upgrade...")

}

func Search(pkg string) {

	search := exec.Command(mgr, "search", pkg)
	utils.RunCmd(search, "Search Error:")

}

func InstallAll() {

	// brew
	fmt.Println("Brew:\n")
	pkgs, fperr := db.ReadPkgs("", "packages", "brew")
	if fperr != nil {
		utils.PrintErrorExit("Homebrew - Read ERROR:", fperr)
		os.Exit(1)
	}
	args := append([]string{"install"}, strings.Split(pkgs, " ")...)
	installAll := exec.Command(mgr, args...)
	utils.RunCmd(installAll, "Installation Error:")

}
