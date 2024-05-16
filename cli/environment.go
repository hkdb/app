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

var deb_base = []string{"Ubuntu", "Pop", "Debian", "MX", "Raspbian", "Kali", "Linuxmint"}
var rh_base = []string{"Fedora", "Rocky", "AlmaLinux", "CentOS", "RedHatEnterpriseServer", "Oracle", "ClearOS", "AmazonAMI"}
var arch_base = []string{"Arch", "Garuda", "Manjaro", "Endeavour"}

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
		if yay := os.Getenv("YAY"); yay == "n" {
			env.Yay = false
		}
		if flatpak := os.Getenv("FLATPAK"); flatpak == "n" {
			env.Flatpak = false
		}
		if snap := os.Getenv("SNAP"); snap == "n" {
			env.Snap = false
		}
		if brew := os.Getenv("BREW"); brew == "n" {
			env.Brew = false
		}
		if appimage := os.Getenv("APPIMAGE"); appimage == "n" {
			env.AppImage = false
		}

		d, err := exec.Command("/usr/bin/lsb_release", "-i", "-s").Output()
		if err != nil {
			fmt.Println("Error:", err)
		}
		distro := strings.TrimSuffix(string(d), "\n")

		//fmt.Println("Distro:", distro)
		env.Distro = distro

		// Check if it's a Debian based
		for i := 0; i < len(deb_base); i++ {
			if distro == deb_base[i] {
				dBase = "debian"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a RedHat based
		for i := 0; i < len(rh_base); i++ {
			if distro == rh_base[i] {
				dBase = "redhat"
				//fmt.Println("Base:", dBase + "\n")
				env.Base = dBase
				break
			}
		}

		// Check if it's a Arch based
		for i := 0; i < len(arch_base); i++ {
			if distro == arch_base[i] {
				dBase = "arch"
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
	case "freebsd":
	case "windows":
		utils.PrintErrorMsgExit("Error:", "Windows is not supported yet...")
	}
	
	brew := os.Getenv("BREW")
	if brew == "n" {
		env.Brew = false
	}
	if brew == "" || brew == " " {
		utils.AppendToFile("BREW = n", env.DBDir+"/settings.conf")
	}
	golang := os.Getenv("GOLANG")
	if golang == "n" {
		env.Go = false
	}
	if golang == "" || golang == " " {
		utils.AppendToFile("GOLANG = n", env.DBDir+"/settings.conf")
	}
	pip := os.Getenv("PIP")
	if pip == "n" {
		env.Pip = false
	}
	if pip == "" || pip == " " {
		utils.AppendToFile("PIP = n", env.DBDir+"/settings.conf")
	}
	cargo := os.Getenv("CARGO")
	if cargo == "n" {
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
		if cargo == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Cargo because it's not installed on your system. Suppress this message by disabling Cargo on app by running \"app -m cargo disable\"...\n", utils.ColorReset)
		}
	}

}
