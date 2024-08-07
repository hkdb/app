package mgr

import (
	"github.com/hkdb/app/env"
)

func isEnabled(pm string) bool {

	switch pm {
	case "apt", "dnf", "pacman", "pkg", "zypper", "nix-env", "app":
		return true
	case "yay":
		if env.Yay == false {
			return false
		}
		return true
	case "flatpak":
		if env.Flatpak == false {
			return false
		}
		return true
	case "snap":
		if env.Snap == false {
			return false
		}
		return true
	case "appimage":
		if env.AppImage == false {
			return false
		}
		return true
	case "brew":
		if env.Brew == false {
			return false
		}
		return true
	case "go":
		if env.Go == false {
			return false
		}
		return true
	case "pip":
		if env.Pip == false {
			return false
		}
		return true
	case "cargo":
		if env.Cargo == false {
			return false
		}
		return true
	}

	return false

}
