## USAGE

```
USAGE:
	app [OPTIONS] <ACTION> <PACKAGE>
	EXAMPLE:
		app install neovim
		app -m flatpak install Geary

ACTIONS:
	- install ~ Install package. Takes package name as argument
	- remove ~ Uninstall package. Takes package name as argument
	- update ~ Refreshes repos
	- upgrade ~ Upgrade packages. Takes "all" as a value to upgrade with all package managers
	- dist-upgrade ~ A more advanced upgrade that can add or remove packages during upgrade (APT Only)
	- autoremove ~ Remove dependency packages that are no longer required
	- purge ~ Same as remove but removes configs too (APT only)
	- search ~ Search for packages in repos
	- list ~ List packages installed on system. Greps for package if argument is provided.
	- history ~ List pacakges installed by app. Takes package name as argument to search.
	- enable ~ Enable Package Manager (Flatpak, Snap, AppImage)
	- disable ~ Disable Package Manager (Flatpak, Snap, AppImage)
	- add-repo ~ Add package manager repo. Takes a .sh, ppa, or url as argument.
	- rm-repo ~ Remove package manager repo. Takes repo identifier as argument
	- ls-repo ~ List package manager repos
	- settings ~ List settings including the status of packages managers (enabled/disabled)

PACKAGE:
	Package name(s). For multiple packages, wrap the argument with quotes.
	EXAMPLE:
		app install 'neovim whois nmap'

OPTIONS:
  -c string
    	Channel - Used in combination with installing snap packages (SNAP Only)
    	   usage: app -m snap -c <channel> install vlc
    	   options:
    		- beta
    		- candidate
    		- edge
    		- stable
    	
  -classic
    	Classic Confinement for Snaps (SNAP Only)
    	   usage: app -m snap -classic install flow
    	
  -gpg string
    	PGP Key URL - Used in combination with a url arg for add-repo (DNF Only)
    	   usage: app -gpg <url> add-repo <url>
    	
  -m string
    	Package Manager
    	   usage: app -m <package manager> install neovim
    	   default: auto-detect of native pkg manager <apt/dnf/pacman>
    	   example: app -a install -p neovim
    	   options:
    		- apt (default if debian based)
    		- dnf (default if redhat based)
    		- pacman (default if arch based)
    		- yay
    		- flatpak
    		- snap
            - brew (default if macOS)
            - go
            - pip
            - cargo
    		- appimage
    	
  -r string
    	Restore / Install all on new system
    	   usage: app -r <type>
    	   option:
    		- apt
    		- dnf
    		- pacman
    		- yay
    		- flatpak
    		- snap
            - brew
            - go
            - pip
            - cargo
    		- appimage
    		- all
  
  -tag string
    	Tag (version) for cargo
    	   usage: app -m cargo -tag <version> install <git url>
    	   example: app -m cargo -tag 0.2.0 install https://github.com/donovanglover/hyprnome 	
  
  -y	Auto Yes - Skips the package manager confirmation (APT & DNF Only)
    	   usage: app -y install neovim
    	
```

