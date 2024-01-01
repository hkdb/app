package appimage

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/utils"

	"os"
	"fmt"
	"strings"
)

func Install(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "appimage", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == true {
		utils.PrintErrorMsgExit(pkg + " is already installed...", "")
	}

	name := installPkgs(pkg, false)

	fmt.Println("\n Recording " + pkg + " to app history...\n")
	rerr := db.RecordPkg("", "packages", "appimage", name)
	if rerr != nil {
		utils.PrintErrorExit("Record Error: ", rerr)
	}

}

func Remove(pkg string) {

	// Check if package is already installed
	inst, ierr := db.IsInstalled("", "packages", "appimage", pkg)
	if ierr != nil {
		utils.PrintErrorExit("Install Check Error:", ierr)
	}
	
	if inst == false {
		utils.PrintErrorMsgExit(pkg + " was not installed by app...", "")
	}

	removePkg(pkg)

	fmt.Println("\n Removing " + pkg + " from app history...\n")
	derr := db.RemovePkg("", "packages", "appimage", pkg)
	if derr != nil {
		utils.PrintErrorExit("Delete Error: ", derr)
	}

}

func Purge(pkg string) {

	fmt.Println("This is an apt only command and does not apply to AppImage...")

}

func AutoRemove() {

	fmt.Println("This action does not apply to AppImage...")

}

func ListSystem() {

	fmt.Println("This action does not apply to AppImage...")

}

func ListSystemSearch(pkg string) {

	fmt.Println("This action does not apply to AppImage...") 
	
}

func Update() {

	fmt.Println("This action does not apply to AppImage...") 

}

func Upgrade(pkg string) {
	
	fmt.Println("This action does not apply to AppImage... To upgrade an AppImage, simply remove the old and install the new...")

}

func DistUpgrade() {

	fmt.Println("This is an apt only command. It does not apply to AppImage...")

}

func Search(pkg string) {

	fmt.Println("This action does not apply to AppImage...") 

}

func History() {

	pkgs, err := db.ReadPkgSlice("", "packages", "appimage")
	if err != nil {
		utils.PrintErrorExit("Read History Error:", err)
		os.Exit(1)
	}

	fmt.Println("Packages manually installed by app:")

	if len(pkgs) == 1 && pkgs[0] == "" {
		fmt.Println("\n\nNo packages have been installed with app yet...")
	} else {
		for i := 0; i < len(pkgs); i++ {
			fmt.Println(pkgs[i])
		}
	}

	fmt.Println("\n")

}

func InstallAll() {
	
	// appimage
	fmt.Println("AppImage:\n")
	pkgs, fperr := db.ReadPkgSlice("", "packages", "appimage")
	if fperr != nil {
		utils.PrintErrorExit("AppImage - Read ERROR:", fperr)
		os.Exit(1)
	}
	
	for i := 1; i < len(pkgs); i++ {
		fmt.Println(pkgs[i])
		cpFile, err := db.ReadPkgs("packages/local", "appimage", pkgs[i])
		if err != nil {
			utils.PrintErrorExit("Read Packages Error:", err)
		}

		cpFile = strings.TrimSpace(cpFile)
		cpFileFull := env.DBDir + "/packages/local/appimage/" + pkgs[i] + "/" + cpFile
		destFile := env.HomeDir + "/.local/share/applications/" + cpFile

		if err := utils.Copy(cpFileFull, destFile); err != nil {
			utils.PrintErrorExit("Restore Error:", err)
		}
	}

	fmt.Println("\nAppImage Restore Completed...\n")

}

