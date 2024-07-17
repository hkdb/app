package cli

import (
	"github.com/hkdb/app/env"
	"github.com/hkdb/app/utils"

	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var deb_base = []string{"debian", "ubuntu", "pop", "mx", "kali", "raspbian", "linuxmint"}
var rh_base = []string{"redhat", "fedora", "clearos", "oracle", "rocky", "amazonami" }
var suse_base = []string{"opensuse", "opensuse-leap", "suse"}
var arch_base = []string{"arch", "garuda", "manjaro", "endeavour"}

// Load envfile and get environment variables
func GetEnv() {

	// Get home dir path
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can't find home directory of system...", err)
		os.Exit(1)
	}

	env.HomeDir = homedir
	env.DBDir = homedir + "/.config/app"

	dir := homedir + "/.config/app"
	dBase := ""

	// determine OS
	osType := runtime.GOOS
	switch osType {
	case "linux":
		env.OSType = "Linux"
	case "darwin":
		env.OSType = "Mac"
	case "freebsd":
		env.OSType = "FreeBSD"
		env.Base = "freebsd"
	case "windows":
		env.OSType = "Windows"
	default:
		fmt.Print(utils.ColorRed, "Unsupported Operating System... Exiting...\n\n", utils.ColorReset)
		os.Exit(1)
	}

	// Getting settings from settings.conf if it exists
	if _, conferr := os.Stat(dir + "/settings.conf"); conferr == nil {
		err := godotenv.Load(dir + "/settings.conf")
		if err != nil {
			fmt.Println(utils.ColorRed, "Error loading settings.conf", utils.ColorReset)
		}
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println(utils.ColorYellow, "\nFirst time running... Creating config dir...\n\n", utils.ColorReset)
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Exiting...\n")
				os.Exit(1)
			}
		}

		if werr := utils.WriteToFile("YAY = n\nFLATPAK = n\nSNAP = n\nAPPIMAGE = n", dir+"/settings.conf"); werr != nil {
			utils.PrintErrorExit("Write settings.conf Error:", werr)
		}
	}

	switch osType {
	case "linux":
		if yay := os.Getenv("YAY"); yay != "y" {
			env.Yay = false
		}
		if flatpak := os.Getenv("FLATPAK"); flatpak != "y" {
			env.Flatpak = false
		}
		if snap := os.Getenv("SNAP"); snap != "y" {
			env.Snap = false
		}
		if brew := os.Getenv("BREW"); brew != "y" {
			env.Brew = false
		}
		if appimage := os.Getenv("APPIMAGE"); appimage != "y" {
			env.AppImage = false
		}

		id, err := exec.Command(env.Bash, "-c", "cat /etc/*-release | grep \"^ID=\" | head -1 | cut -d '=' -f 2").Output()
		if err != nil {
			fmt.Print(utils.ColorRed, "Unable to determine distribution... Exiting...\n\n", utils.ColorReset)
			os.Exit(1)
		}
		distro := string(id)

		//fmt.Println("Distro:", distro)
		env.Distro = distro

		i, err := exec.Command(env.Bash, "-c", "cat /etc/*-release").Output()
		if err != nil {
			fmt.Print(utils.ColorRed, "Unable to determine base distribution... Exiting...\n\n", utils.ColorReset)
			os.Exit(1)
		}
		//fmt.Println("Infos:", string(i))
		infos := strings.ToLower(strings.ReplaceAll(string(i), " ", ""))

		// Check if it's a Debian based
		for i := 0; i < len(deb_base); i++ {
			if strings.Contains(infos, deb_base[i]) {
				dBase = "debian"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a RedHat based
		for i := 0; i < len(rh_base); i++ {
			if strings.Contains(infos, rh_base[i]) {
				dBase = "redhat"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a Arch based
		for i := 0; i < len(arch_base); i++ {
			if strings.Contains(infos, arch_base[i]) {
				dBase = "arch"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a Suse based
		for i := 0; i < len(suse_base); i++ {
			if strings.Contains(infos, suse_base[i]) {
				dBase = "suse"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Temporarily disabling package manager that don't exist
		if env.Base == "arch" && env.Yay == true {
			yay, _ := utils.CheckIfExists(env.YayCmd)
			if yay == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Yay because it's not installed on your system. Suppress this message by disabling Yay on app by running \"app -m yay disable\"...\n", utils.ColorReset)
			}
		}

		if env.Flatpak == true {
			flatpak, _ := utils.CheckIfExists(env.FlatpakCmd)
			if flatpak == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Flatpak because it's not installed on your system. Suppress this message by disabling Flatpak on app by running \"app -m flatpak disable\"...\n", utils.ColorReset)
			}
		}

		if env.Snap == true {
			snap, _ := utils.CheckIfExists(env.SnapCmd)
			if snap == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Snap because it's not installed on your system. Suppress this message by disabling Snap on app by running \"app -m snap disable\"...\n", utils.ColorReset)
			}
		}
	case "darwin":

		if werr := utils.WriteToFile("BREW = y", dir+"/settings.conf"); werr != nil {
			utils.PrintErrorExit("Write settings.conf Error:", werr)
		}

		env.Brew = true
		bsdPath()
	case "freebsd":
		bashPath()
		bsdPath()
	case "windows":
		utils.PrintErrorMsgExit("Error:", "Windows is not supported yet...")
	}

	brew := os.Getenv("BREW")
	if brew != "y" {
		env.Brew = false
	}
	if brew == "" || brew == " " {
		utils.AppendToFile("BREW = n", env.DBDir+"/settings.conf")
	}
	golang := os.Getenv("GOLANG")
	if golang != "y" {
		env.Go = false
	}
	if golang == "" || golang == " " {
		utils.AppendToFile("GOLANG = n", env.DBDir+"/settings.conf")
	}
	pip := os.Getenv("PIP")
	if pip != "y" {
		env.Pip = false
	}
	if pip == "" || pip == " " {
		utils.AppendToFile("PIP = n", env.DBDir+"/settings.conf")
	}
	cargo := os.Getenv("CARGO")
	if cargo != "y" {
		env.Cargo = false
	}
	if cargo == "" || cargo == " " {
		utils.AppendToFile("CARGO = n", env.DBDir+"/settings.conf")
	}

	if env.Go == true {
		golang, _ := utils.CheckIfExists(env.GoCmd)
		if golang == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Go because it's not installed on your system. Suppress this message by disabling Go on app by running \"app -m go disable\"...\n", utils.ColorReset)
		}
	}

	if env.Pip == true {
		pip, _ := utils.CheckIfExists(env.PipCmd)
		if pip == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Pip because it's not installed on your system. Suppress this message by disabling Pip on app by running \"app -m pip disable\"...\n", utils.ColorReset)
		}
	}

	if env.Cargo == true {
		cargo, _ := utils.CheckIfExists(env.CargoCmd)
		cargoLocal, _ := utils.CheckIfExists(env.HomeDir + env.CargoLocalCmd)
		if cargo == false && cargoLocal == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Cargo because it's not installed on your system. Suppress this message by disabling Cargo on app by running \"app -m cargo disable\"...\n", utils.ColorReset)
		}
	}

}

func bashPath() {

	env.Bash = "/usr/local/bin/bash"

}

func bsdPath() {

	env.GoCmd = "/usr/local/bin/go"
	env.PipCmd = "/usr/local/bin/pip"
	env.CargoCmd = "/usr/local/bin/pip"

}
