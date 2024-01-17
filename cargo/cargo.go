package cargo

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"os"
	"os/exec"
	"fmt"
)

var mgr = "cargo"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", mgr, pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	install := exec.Command("/bin/bash", "-c", mgr + " install " + pkg)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", mgr, pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", mgr, pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}

	remove := exec.Command(mgr, "uninstall", pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", mgr, pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m cargo remove " + pkg + "...")

}

func AutoRemove() {

	fmt.Println("This is an apt only command. It does not apply to Cargo...")

}

func ListSystem() {

	list := exec.Command("ls", "-lah", env.HomeDir + "/.cargo/bin")
	utils.RunCmd(list, "List Package Error:")

}

func ListSystemSearch(pkg string) {

	listSearch := exec.Command("/bin/bash", "-c", "ls -lah " + "/.cargo/bin |grep " + pkg)
	utils.RunCmd(listSearch, "List Package Search Error:")
}

func Update() {

	fmt.Println("This is an apt only command. Just use app -m cargo upgrade...")

}

func Upgrade() {
	
	chkDep := exec.Command("/bin/bash", "-c", "cargo --list |grep install-update")
	err := utils.ChkIfCmdRuns(chkDep)
	if err != nil {
		fmt.Print("The \"install-update\" Cargo command isn't installed... Do you want to install it? (Y/n) ")
		resp := utils.Confirm()
		if resp == true {
			installDep := exec.Command("cargo", "install", "cargo-update")
			utils.RunCmd(installDep, "Upgrade Dependency Installation Error:")
		}
	}

	upgrade := exec.Command(mgr, "install-update", "-a")
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m cargo upgrade...")

}

func Search(pkg string) {

	search := exec.Command(mgr, "search", pkg)
	utils.RunCmd(search, "Search Error:")

}

func InstallAll() {
	
	// cargo
	fmt.Println("Cargo:\n")
	pkgs, fperr := db.ReadPkgs("", "packages", mgr)
	if fperr != nil {
		utils.PrintErrorExit("Cargo - Read ERROR:", fperr)
		os.Exit(1)
	}
	installAll := exec.Command(mgr, "install", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}

