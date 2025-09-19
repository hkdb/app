package env

type Flags struct {
	A       string // Action
	P       string // Package or Other types of value for Action
	M       string // Package Manager
	R       string // Restore
	Y       bool   // Auto-Yes for certain package managers
	G       string // URL to GPG key for dnf/yum/rpm based distros when adding repo
	C       string // Channel for specifying channel when install snaps
	Classic bool   // Classic confinement for snaps
	Tag     string // Tag (version) for installing git url with cargo
	Sort		bool
}

var Version string
var OSType string
var Distro string
var Base string
var HomeDir string
var DBDir string
var Env = "/usr/bin/env"
var Bash = "bash"
var Header = true
var Apt = true
var Dnf = true
var Pacman = true
var Zypper = true
var Yay = true
var YayCmd = "/usr/bin/yay"
var Paru = true
var ParuCmd = "/usr/bin/paru"
var NixEnv = true
var Flatpak = true
var FlatpakCmd = "/usr/bin/flatpak"
var Snap = true
var SnapCmd = "/usr/bin/snap"
var Brew = true
var BrewIntelCmd = "/usr/local/bin/brew"
var BrewLinuxCmd = "/home/linuxbrew/.linuxbrew/bin/brew"
var BrewSiliconCmd = "/opt/homebrew/bin/brew"
var BrewCmd = BrewIntelCmd
var Go = true
var GoCmd = "/usr/bin/go"
var Pip = true
var PipCmd = "/usr/bin/pip"
var Cargo = true
var CargoCmd = "/usr/bin/cargo"
var CargoLocalCmd = "/.cargo/bin/cargo"
var AppImage = true
var AutoYes = false
