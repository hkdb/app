package snap

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
var mgr = "/usr/bin/snap"

func Install(pkg, c string, classic bool) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "snap", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == true {
		utils.PrintErrorMsgExit(pkg+" has already been installed...", "")
	}

	cmd := "install "
	end := ""
	if c != "" {
		switch c {
		case "edge", "stable", "candidate", "beta":
			cmd = cmd + "--channel=" + c + " "
		default:
			utils.PrintErrorMsgExit("Input Error:", "Not a valid channel for snap...")
		}
	}
	if classic == true {
		end = " --classic"
	}
	install := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], "/usr/bin/snap "+cmd+pkg+end)
	utils.RunCmd(install, "Installation Error:")

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	if c != "" {
		rcerr := db.RecordChan("snap", pkg, c)
		if rcerr != nil {
			utils.PrintErrorExit("Channel Record Error:", rcerr)
		}
	}
	clerr := db.RecordClassic("snap", pkg, classic)
	if clerr != nil {
		utils.PrintErrorExit("Confinement Record Error:", clerr)
	}
	rerr := db.RecordPkg("", "packages", "snap", pkg)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "snap", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}

	if inst == false {
		utils.PrintErrorMsgExit(pkg+" was not installed by app...", "")
	}

	remove := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], "/usr/bin/snap remove "+pkg)
	utils.RunCmd(remove, "Remove Error:")

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	hasChan := db.ChannelPreferred("snap", pkg)
	if hasChan == true {
		rcerr := db.RemoveChan("snap", pkg)
		if rcerr != nil {
			utils.PrintErrorExit("Channel Record Removal Error:", rcerr)
		}
	}
	clerr := db.RemoveClassic("snap", pkg)
	if clerr != nil {
		utils.PrintErrorExit("Confinement Record Error:", clerr)
	}
	derr := db.RemovePkg("", "packages", "snap", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

	utils.DeleteDirIfEmpty(env.DBDir + "packages/repo/channel/snap")

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command. Just use app -m snap -a remove -p " + pkg + "...")

}

func AutoRemove() {

	fmt.Println("This action does not apply to Snap...")

}

func ListSystem() {

	err := syscall.Exec(mgr, []string{mgr, "list"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Error:", err)
	}

}

func ListSystemSearch(pkg string) {

	err := syscall.Exec(mgr, []string{mgr, "list", "|", "grep", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("List System Search Error:", err)
	}

}

func Update() {

	err := syscall.Exec(sudo[0], []string{sudo[0], sudo[1], sudo[2], sudo[3], "/usr/bin/snap refresh"}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Update Error:", err)
	}

}

func Upgrade() {

	upgrade := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], "/usr/bin/snap refresh")
	utils.RunCmd(upgrade, "Upgrade Error:")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. Just use app -m snap upgrade...")

}

func Search(pkg string) {

	err := syscall.Exec("/usr/bin/snap", []string{"/usr/bin/snap", "search", pkg}, os.Environ())
	if err != nil {
		utils.PrintErrorExit("Search Error", err)
	}

}

func InstallAll() {

	// snap
	action := " install "

	fmt.Println("Snap:\n")

	nochan, _ := utils.DirIsEmpty(env.DBDir + "/packages/repo/channel/snap")
	if nochan == true {
		pkgs, sperr := db.ReadPkgs("", "packages", "snap")
		if sperr != nil {
			utils.PrintErrorExit("Snap - Read ERROR:", sperr)
		}
		classic := ""
		isClassic, icerr := db.GetClassic("snap", pkgs)
		if icerr != nil {
			utils.PrintErrorExit("Snap - Confinement Preference Read Error:", icerr)
		}
		command := mgr + action
		if isClassic == true {
			classic = " --classic"
		}
		install := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], command+" "+pkgs+classic)
		utils.RunCmd(install, "Installation Error:")
	} else {
		pkgs, sperr := db.ReadPkgSlice("", "packages", "snap")
		if sperr != nil {
			utils.PrintErrorExit("Snap - Read ERROR:", sperr)
		}
		for i := 0; i < len(pkgs); i++ {
			command := mgr + action
			hasChan := db.ChannelPreferred("snap", pkgs[i])
			if hasChan == true {
				ch, err := db.GetChan("snap", pkgs[i])
				if err != nil {
					utils.PrintErrorExit("Get Channel Error:", err)
				}
				command = command + " --channel=" + ch
			}

			classic := ""
			isClassic, err := db.GetClassic("snap", pkgs[i])
			if err != nil {
				utils.PrintErrorExit("Confinement Read Error:", err)
			}

			if isClassic == true {
				classic = " --classic"
			}

			install := exec.Command(sudo[0], sudo[1], sudo[2], sudo[3], command+" "+pkgs[i]+classic)
			utils.RunCmd(install, "Installation Error:")
		}
	}

}
