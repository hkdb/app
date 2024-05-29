package pip

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
)

var mgr = "pip"

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", mgr, pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" is already installed...", "")
	}

	install := exec.Command(mgr, "install", pkg)
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
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
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

	fmt.Println("This is an apt only command. Just use app -m pip remove " + pkg + "...")

}

func AutoRemove() {

	fmt.Println("This is an apt only command. It does not apply to pip...")

}

func ListSystem() {

	list := exec.Command(mgr, "list")
	utils.RunCmd(list, "List Package Error:")

}

func ListSystemSearch(pkg string) {

	listSearch := exec.Command(env.Bash, "-c", mgr+" list |grep "+pkg)
	utils.RunCmd(listSearch, "List Package Search Error:")
}

func Update() {

	fmt.Println("This is an apt only command. Just use app -m pip upgrade...")

}

func Upgrade() {

	chkDep := exec.Command(env.Bash, "-c", "pip-review")
	err := utils.ChkIfCmdRuns(chkDep)
	if err != nil {
		fmt.Print("The pip-review command isn't installed... Do you want to install it? (Y/n) ")
		resp := utils.Confirm()
		if resp == true {
			installDep := exec.Command("pip", "install", "pip-review")
			utils.RunCmd(installDep, "Upgrade Dependency Installation Error:")
		}
	}

	upgrade := exec.Command("pip-review", "--local", "--interactive")
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m pip upgrade...")

}

func Search(pkg string) {

	fmt.Println("This is not a supported action. ... Search for packages at https://pypi.org/search")

}

func InstallAll() {

	// pip
	fmt.Println("Pip:\n")
	pkgs, fperr := db.ReadPkgs("", "packages", mgr)
	if fperr != nil {
		utils.PrintErrorExit("Pip - Read ERROR:", fperr)
		os.Exit(1)
	}
	installAll := exec.Command(mgr, "install", pkgs)
	utils.RunCmd(installAll, "Installation Error:")

}
