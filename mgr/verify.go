package mgr

import (
	"github.com/hkdb/app/env"
)

func isEnabled(pm string) bool {

	switch pm {
	case "pkg", "app":
		return true
	case "apt":
		if env.Apt == false {
			return false
		}
		return true
	case "dnf":
		if env.Dnf == false {
			return false
		}
		return true
	case "pacman":
		if env.Pacman == false {
			return false
		}
		return true
	case "zypper":
		if env.Zypper == false {
			return false
		}
		return true
	case "nix-env":
		if env.NixEnv == false {
			return false
		}
		return true
	case "yay":
		if env.Yay == false {
			return false
		}
		return true
	case "paru":
		if env.Paru == false {
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
