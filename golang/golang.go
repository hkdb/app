package golang

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"os"
	"os/exec"
	"fmt"
)

var mgr = "go"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "go", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	install := exec.Command(mgr, "install", pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "go", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "go", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}

	remove := exec.Command("rm", env.HomeDir + "/go/bin/" + pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "go", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m go remove " + pkg + "...")

}

func AutoRemove() {

	fmt.Println("This is an apt only command. It does not apply to go...")

}

func ListSystem() {

	list := exec.Command("ls", "-lah", env.HomeDir + "/go/bin")
	utils.RunCmd(list, "List Package Error:")

}

func ListSystemSearch(pkg string) {

	listSearch := exec.Command("ls", "-lah", env.HomeDir + "/go/bin", "|", "grep", pkg)
	utils.RunCmd(listSearch, "List Package Search Error:")
}

func Update() {

	fmt.Println("This is an apt only command. Just use app -m go upgrade...")

}

func Upgrade() {

	chkDep := exec.Command("/bin/bash", "-c", "go-global-update")
	err := utils.ChkIfCmdRuns(chkDep)
	if err != nil {
		fmt.Print("The \"go-global-update\" command isn't installed... Do you want to install it? (Y/n) ")
		resp := utils.Confirm()
		if resp == true {
			installDep := exec.Command("go", "install", "github.com/Gelio/go-global-update@latest")
			utils.RunCmd(installDep, "Upgrade Dependency Installation Error:")
		}
	}

	upgrade := exec.Command("go-global-update")
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m go upgrade...")

}

func Search(pkg string) {

	utils.PrintErrorMsgExit("Errors:", "This is not a supported action...")

}

func InstallAll() {
	
	// go
	fmt.Println("go:\n")
	pkgs, fperr := db.ReadPkgs("", "packages", "go")
	if fperr != nil {
		utils.PrintErrorExit("go - Read ERROR:", fperr)
		os.Exit(1)
	}
	installAll := exec.Command(mgr, "install", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}

