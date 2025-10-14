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
var arch_base = []string{"arch", "garuda", "manjaro", "endeavour", "cachyos"}
var nixos_base = []string{"nixos"}

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
		fmt.Print(utils.ColorRed, "🚨 Unsupported Operating System... Exiting...\n\n", utils.ColorReset)
		os.Exit(1)
	}

	// Getting settings from settings.conf if it exists
	if _, conferr := os.Stat(dir + "/settings.conf"); conferr == nil {
		err := godotenv.Load(dir + "/settings.conf")
		if err != nil {
			fmt.Println(utils.ColorRed, "🚨 Error loading settings.conf", utils.ColorReset)
		}
		migration := migrationNeeded(dir+"/settings.conf")
		if migration {
			fmt.Println("🏗️ Settings file requires a migration to the latest format...\n")
			err := migrateConf(dir+"/settings.conf")
			if err != nil {
				utils.PrintErrorExit("Migration Error:", err)
			}
			fmt.Println("🚀 Settings file migration completed...\n")
			fmt.Println("⚓ Reloading settings file...\n")
			lerr := godotenv.Load(dir + "/settings.conf")
			if lerr != nil {
				fmt.Println(utils.ColorRed, "🚨 Error loading settings.conf", utils.ColorReset)
			}
		}
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println(utils.ColorYellow, "\n🏗️ First time running... Creating config dir...\n\n", utils.ColorReset)
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				fmt.Println("🚨 Error:", err)
				fmt.Println("Exiting...\n")
				os.Exit(1)
			}
		}

		if werr := utils.WriteToFile("APP_FLATPAK = n\nAPP_SNAP = n\nAPP_APPIMAGE = n", dir+"/settings.conf"); werr != nil {
			utils.PrintErrorExit("Write settings.conf Error:", werr)
		}

	}

	switch osType {
	case "linux":
		yay := os.Getenv("APP_YAY")
		if yay != "y" {
			env.Yay = false
		}
		if yay == "" || yay == " " {
			utils.AppendToFile("APP_YAY = n", env.DBDir+"/settings.conf")
		}
		paru := os.Getenv("APP_PARU")
		if paru != "y" {
			env.Paru = false
		}
		if paru == "" || paru == " " {
			utils.AppendToFile("APP_PARU = n", env.DBDir+"/settings.conf")
		}
		flatpak := os.Getenv("APP_FLATPAK")
		if flatpak != "y" {
			env.Flatpak = false
		}
		if flatpak == "" || flatpak == " " {
			utils.AppendToFile("APP_FLATPAK = n", env.DBDir+"/settings.conf")
		}
		snap := os.Getenv("APP_SNAP")
		if snap != "y" {
			env.Snap = false
		}
		if snap == "" || snap == " " {
			utils.AppendToFile("APP_SNAP = n", env.DBDir+"/settings.conf")
		}
		env.BrewCmd = env.BrewLinuxCmd
		if brew := os.Getenv("APP_BREW"); brew != "y" {
			env.Brew = false
		}
		if appimage := os.Getenv("APP_APPIMAGE"); appimage != "y" {
			env.AppImage = false
		}

		id, err := exec.Command(env.Env, env.Bash, "-c", "cat /etc/os-release | grep \"^ID=\" | head -1 | cut -d '=' -f 2").Output()
		if err != nil {
			fmt.Print(utils.ColorRed, "Unable to determine distribution... Exiting...\n\n", utils.ColorReset)
			os.Exit(1)
		}
		distro := string(id)

		//fmt.Println("Distro:", distro)
		env.Distro = distro

		i, err := exec.Command(env.Env, env.Bash, "-c", "cat /etc/os-release").Output()
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

		// Check if it's a Nixos based
		for i := 0; i < len(nixos_base); i++ {
			if strings.Contains(infos, nixos_base[i]) {
				dBase = "nixos"
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

		npm := utils.GetNativePkgMgr()
		apt := os.Getenv("APP_APT")
		onoroff := "n"
		if apt != "y" {
			env.Apt = false
		}
		if apt == "" || apt == " " {
			if npm == "apt" {
				onoroff = "y"
			}
			utils.AppendToFile("APP_APT = " + onoroff, env.DBDir+"/settings.conf")
			onoroff = "n"
		}
		dnf := os.Getenv("APP_DNF")
		if dnf != "y" {
			env.Dnf = false
		}
		if dnf == "" || dnf == " " {
			if npm == "dnf" {
				onoroff = "y"
			}
			utils.AppendToFile("APP_DNF = " + onoroff, env.DBDir+"/settings.conf")
			onoroff = "n"
		}
		pacman := os.Getenv("APP_PACMAN")
		if pacman != "y" {
			env.Pacman = false
		}
		if pacman == "" || pacman == " " {
			if npm == "pacman" {
				onoroff = "y"
			}
			utils.AppendToFile("APP_PACMAN = " + onoroff, env.DBDir+"/settings.conf")
			onoroff = "n"
		}
		zypper := os.Getenv("APP_ZYPPER")
		if zypper != "y" {
			env.Pacman = false
		}
		if zypper == "" || zypper == " " {
			if npm == "zypper" {
				onoroff = "y"
			}
			utils.AppendToFile("APP_ZYPPER = " + onoroff, env.DBDir+"/settings.conf")
			onoroff = "n"
		}
		nixenv := os.Getenv("APP_NIXENV")
		if nixenv != "y" {
			env.NixEnv = false
		}
		if nixenv == "" || nixenv == " " {
			if npm == "nix-env" {
				onoroff = "y"
			}
			utils.AppendToFile("APP_NIXENV = " + onoroff, env.DBDir+"/settings.conf")
			onoroff = "n"
		}


		// Temporarily disabling package manager that don't exist
		if env.Base == "arch" && env.Yay == true {
			yay, _ := utils.CheckIfExists(env.YayCmd)
			if yay == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Yay because it's not installed on your system. Suppress this message by disabling Yay on app by running \"app -m yay disable\"...\n", utils.ColorReset)
			}
		}

		if env.Base == "arch" && env.Paru == true {
			paru, _ := utils.CheckIfExists(env.ParuCmd)
			if paru == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Paru because it's not installed on your system. Suppress this message by disabling Paru on app by running \"app -m paru disable\"...\n", utils.ColorReset)
			}
		}

		if env.Flatpak == true { 
			flatpak, _ := utils.CheckIfExists(env.FlatpakCmd)
			if env.Base == "nixos" {
				flatpak = nixosCheck("flatpak", 8)
			}
			if flatpak == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Flatpak because it's not installed on your system. Suppress this message by disabling Flatpak on app by running \"app -m flatpak disable\"...\n", utils.ColorReset)
			}
		}

		if env.Snap == true {
			snap, _ := utils.CheckIfExists(env.SnapCmd)
			if env.Base == "nixos" {
				snap = nixosCheck("snap", 5)
			}
			if snap == false {
				fmt.Println(utils.ColorYellow, "Temporarily disabling Snap because it's not installed on your system. Suppress this message by disabling Snap on app by running \"app -m snap disable\"...\n", utils.ColorReset)
			}
		}
	case "darwin":
		settings := os.Getenv("APP_BREW")
		sCount := len(settings)
		if sCount == 0 {
			utils.AppendToFile("APP_BREW = y", env.DBDir+"/settings.conf")
			err := godotenv.Load(dir + "/settings.conf")
			if err != nil {
				fmt.Println(utils.ColorRed, "Error loading settings.conf", utils.ColorReset)
			}
		}
		env.Brew = true
		bsdPath()
	case "freebsd":
		// bashPath()
		bsdPath()
	case "windows":
		utils.PrintErrorMsgExit("Error:", "Windows is not supported yet...")
	}

	brew := os.Getenv("APP_BREW")
	if brew != "y" {
		env.Brew = false
	}
	if brew == "" || brew == " " {
		utils.AppendToFile("APP_BREW = n", env.DBDir+"/settings.conf")
	}
	golang := os.Getenv("APP_GOLANG")
	if golang != "y" {
		env.Go = false
	}
	if golang == "" || golang == " " {
		utils.AppendToFile("APP_GOLANG = n", env.DBDir+"/settings.conf")
	}
	pip := os.Getenv("APP_PIP")
	if pip != "y" {
		env.Pip = false
	}
	if pip == "" || pip == " " {
		utils.AppendToFile("APP_PIP = n", env.DBDir+"/settings.conf")
	}
	cargo := os.Getenv("APP_CARGO")
	if cargo != "y" {
		env.Cargo = false
	}
	if cargo == "" || cargo == " " {
		utils.AppendToFile("APP_CARGO = n", env.DBDir+"/settings.conf")
	}

	if env.Brew == true {
		brew, _ := utils.CheckIfExists(env.BrewCmd)
		if brew == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Brew because it's not installed on your system. Suppress this message by disabling Brew on app by running \"app-m brew disable\"...\n", utils.ColorReset)
		}
	}

	if env.Go == true {
		golang, _ := utils.CheckIfExists(env.GoCmd)
		if env.Base == "nixos" {
			golang = nixosCheck("go", 3)
		}
		if golang == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Go because it's not installed on your system. Suppress this message by disabling Go on app by running \"app -m go disable\"...\n", utils.ColorReset)
		}
	}

	if env.Pip == true {
		pip, _ := utils.CheckIfExists(env.PipCmd)
		if env.Base == "nixos" {
			pip = nixosCheck("pip", 4)
		}
		if pip == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Pip because it's not installed on your system. Suppress this message by disabling Pip on app by running \"app -m pip disable\"...\n", utils.ColorReset)
		}
	}

	if env.Cargo == true {
		cargo, _ := utils.CheckIfExists(env.CargoCmd)
		if env.Base == "nixos" {
			cargo = nixosCheck("cargo", 6)
		}
		cargoLocal, _ := utils.CheckIfExists(env.HomeDir + env.CargoLocalCmd)
		if cargo == false && cargoLocal == false {
			fmt.Println(utils.ColorYellow, "Temporarily disabling Cargo because it's not installed on your system. Suppress this message by disabling Cargo on app by running \"app -m cargo disable\"...\n", utils.ColorReset)
		}
	}

}


func nixosCheck(c string, l int) bool {

	cmd := exec.Command("whereis", c)
	out, err := utils.RunCmdReturn(cmd)
	if err != nil {
		utils.PrintErrorExit("Error:", err)
	}
	// Number of characters in command + :
	if len(out) <= l {
		return false
	}
	return true

}

/*

func bashPath() {

	env.Bash = "/usr/local/bin/bash"

}

*/

func bsdPath() {

	env.GoCmd = "/usr/local/bin/go"
	env.PipCmd = "/usr/local/bin/pip"
	env.CargoCmd = "/usr/local/bin/pip"
	
	if runtime.GOARCH == "arm64" && runtime.GOOS =="darwin" {
		env.BrewCmd = env.BrewSiliconCmd
	}

}

func migrationNeeded(s string) bool {
	confver := os.Getenv("APP_CONFIG_VER")
	if confver != env.ConfVer {
		return true
	}
	return false
}

func migrateConf(s string) error {
	content, err := os.ReadFile(s)
	if err != nil {
			return err
	}

  // Add APP_ prefix to each line
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
			if strings.TrimSpace(line) != "" && strings.Contains(line, "=") {
					lines[i] = "APP_" + line
			}
	}

	// Add version at the top
  newContent := "APP_CONFIG_VER = 2\n" + strings.Join(lines, "\n")

	return os.WriteFile(s, []byte(newContent), 0644)
}
