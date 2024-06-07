package appimage

import (
	"github.com/hkdb/app/db"
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"strings"
)

func installPkgs(input string, restore bool) string {

	pkgs := strings.Split(input, " ")

	if len(pkgs) < 1 {
		utils.PrintErrorMsgExit("Input Error:", "No package was specified...\n")
	}
	if len(pkgs) < 2 {
		pname := installPkg(input, restore)
		return pname
	}

	pkgNames := ""
	if len(pkgs) >= 2 {
		for i := 0; i < len(pkgs); i++ {
			name := installPkg(pkgs[i], restore)
			pkgNames = pkgNames + " " + name
		}
	}

	return pkgNames

}

func removePkgs(input string) {

	pkgs := strings.Split(input, " ")

	if len(pkgs) < 1 {
		utils.PrintErrorMsgExit("Input Error:", "No package was specified...\n")
	}
	if len(pkgs) < 2 {
		removePkg(input)
		return
	}
	if len(pkgs) >= 2 {
		for i := 0; i < len(pkgs); i++ {
			removePkg(pkgs[i])
		}
	}

}

func installPkg(pkg string, restore bool) string {

	tmp := "/tmp/app/appimage"
	if _, terr := os.Stat(tmp); os.IsNotExist(terr) {
		merr := os.MkdirAll(tmp, os.ModePerm)
		if merr != nil {
			utils.PrintErrorExit("AppImage File Error:", merr)
		}
	}

	// Generate a random string of 10 characters as the work folder
	workdir, err := utils.RandString(10)
	workdirFull := tmp + "/" + workdir
	werr := os.MkdirAll(workdirFull, os.ModePerm)
	if werr != nil {
		utils.PrintErrorExit("AppImage File Error:", werr)
	}

	// Get current work path of where command was executed
	path, err := os.Getwd()
	if err != nil {
		utils.PrintErrorExit("AppImage Path Read Error:", err)
	}

	tappimg := path + "/" + pkg
	// Copy AppImage file to tmp dir
	utils.Copy(tappimg, workdirFull+"/"+pkg)

	// Set perms for AppImage file
	os.Chmod(workdirFull+"/"+pkg, 0755)

	// Extract AppImage to get .desktop and icon
	extract := exec.Command(workdirFull+"/"+pkg, "--appimage-extract")
	extract.Dir = workdirFull
	fmt.Println("Extracting AppImage... please wait...")
	utils.RunCmdQuiet(extract, "Extract Error:")

	confdir := env.DBDir + "/packages/local/appimage"

	// Check if conf dir exists and make it if it doesn't
	if _, cerr := os.Stat(confdir); os.IsNotExist(cerr) {
		cmerr := os.MkdirAll(confdir+"/"+"", os.ModePerm)
		if cmerr != nil {
			utils.PrintErrorExit("AppImage File Error:", cmerr)
		}
	}

	tmpAppDir := workdirFull + "/" + "squashfs-root"

	// Get a list of files inside AppImage
	fileList, ferr := os.ReadDir(tmpAppDir)
	if ferr != nil {
		utils.PrintErrorExit("Read AppImage Dir Error:", ferr)
	}

	// Get .desktop file name
	var dotDesktop = ""
	for i := 0; i < len(fileList); i++ {
		ext := utils.GetFileExtension(fileList[i].Name())
		if ext == "desktop" {
			dotDesktop = fileList[i].Name()
			break
		}
	}
	if dotDesktop == "" {
		utils.PrintErrorMsgExit("Error:", "No .desktop file found in AppImage...")
	}

	tmpDesktop := tmpAppDir + "/tmp.desktop"
	utils.Copy(tmpAppDir+"/"+dotDesktop, tmpDesktop)
	utils.RemoveFileLine(tmpDesktop, 1)

	// Load .desktop values
	dd, gerr := utils.ReadDotDesktop(tmpDesktop)
	if gerr != nil {
		utils.PrintErrorExit("Load .desktop Error:", gerr)
	}
	name := dd.Name
	if name == "" {
		utils.PrintErrorMsgExit("Error:", "AppImage application name could not be read...")
	}

	name = strings.ReplaceAll(name, " ", "_")
	iconName := dd.Icon
	// Make dir for app
	appDir := confdir + "/" + name
	if nerr := os.MkdirAll(appDir, os.ModePerm); nerr != nil {
		utils.PrintErrorExit("Make Dir Error:", nerr)
	}

	// Get icon file name
	var icon = ""
	for i := 0; i < len(fileList); i++ {
		if fileList[i].Name() == iconName+".png" {
			icon = fileList[i].Name()
			break
		}
	}
	if icon == "" {
		utils.PrintErrorMsgExit("Error:", "No icon found in AppImage...")
	}

	// Copy AppImage to conf dir
	utils.Copy(tappimg, appDir+"/"+pkg)
	// Set perms for AppImage file
	os.Chmod(appDir+"/"+pkg, 0755)
	// Copy AppImage icon file to conf dir
	utils.Copy(tmpAppDir+"/"+icon, appDir+"/"+icon)
	// Edit AppImage .desktop and write to config dir for storing and to ~/.local/share/applications/
	passPkg := appDir + "/" + pkg
	passIcon := appDir + "/" + icon
	passDotDesktop := tmpAppDir + "/" + dotDesktop
	passDestDotDesktop := appDir + "/" + dotDesktop
	passShareDotDesktop := env.HomeDir + "/.local/share/applications/" + dotDesktop
	utils.WriteDotDesktop(passPkg, passIcon, passDotDesktop, passDestDotDesktop, passShareDotDesktop)

	// Record package and file name association
	rferr := db.RecordPkg("packages/local", "appimage", name, dotDesktop)
	if rferr != nil {
		utils.PrintErrorExit("Record Local Package Association Error: ", rferr)
	}

	// Cleanup
	utils.RemoveDirRecursive(workdirFull)

	return name

}

func removePkg(pkg string) {

	sysDir := env.HomeDir + "/.local/share/applications"
	confDir := env.DBDir + "/packages/local/appimage"
	appDir := confDir + "/" + pkg

	rmFileRecord, err := db.ReadPkgs("packages/local", "appimage", pkg)
	if err != nil {
		utils.PrintErrorExit("Read Package Error:", err)
	}
	rmFile := strings.TrimSpace(rmFileRecord)

	// Remove .desktop from ~/.local/share/applications
	utils.RemoveFile(sysDir + "/" + rmFile)

	// Remove all App Dir
	utils.RemoveDirRecursive(appDir)

	// Remove package.json
	utils.RemoveFile(confDir + "/" + pkg + ".json")

	// Remove conf dir if empty
	confDirIsEmpty, err := utils.DirIsEmpty(confDir)
	if err != nil {
		utils.PrintErrorExit("Read Dir Error:", err)
	}
	if confDirIsEmpty == true {
		utils.RemoveDirRecursive(confDir)
	}

}
