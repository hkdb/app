package env

type Flags struct {
	A string // Action
	P string // Package or Other types of value for Action
	M	string // Package Manager
	R string // Restore
	Y bool	 // Auto-Yes for certain package managers
	G string // URL to GPG key for dnf/yum/rpm based distros when adding repo 
	C string // Channel for specifying channel when install snaps
	Classic bool
}

var OSType string
var Distro string
var Base string
var HomeDir string
var DBDir string
var Header = true
var Yay = true
var YayCmd = "/usr/bin/yay"
var Flatpak = true
var FlatpakCmd = "/usr/bin/flatpak"
var Snap = true
var SnapCmd = "/usr/bin/snap"
var Brew = true
var AppImage = true
var AutoYes = false
