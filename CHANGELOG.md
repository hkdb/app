# CHANGELOG

#### v0.33

- Fixed suse InstallAll
- Fixed suse RestoreAll
- Fixed output typo
- Fixed brew install multiple packages


#### v0.32

- Added `pipx` support
- Restructured flag checks for future refactor


#### v0.31

v0.29 breaks enable/disable and also pretty much cripples app from working for first time installers. This release fixes that:

- Fixed enable/disable of package managers
- Fixed first time install logic


#### v0.30

- Fixed brew upgrade


#### v0.29

- Changed settings.conf format to resolve paru and yay variable naming conflict that causes installations to fail
- Auto-migration of settings file to v2 format.


#### v0.28

- Added -S (--stdin) to all sudo commands so that users can directly run app inside a container from the host(ie. `incus exec container -- app -r all`)


#### v0.27

- Fixed restore repositories
- Fixed restore go packages
- Fixed restore cargo packages


#### v0.26

- Arch, AUR, & Enable/Disable Native Package Managers [#24](https://github.com/hkdb/app/issues/24)


#### v0.25

- Changed single installation command to be more universal
- Fixed fish shell support
- Fixed aarch64 support [#22](https://github.com/hkdb/app/issues/22)
- Added reminder for enable/disable synatax errors [#23](https://github.com/hkdb/app/issues/23)


#### v0.24

- Fixed restoration of .sh added repositories


#### v0.23

- Really fixed mac self update...


#### v0.22

- Fixed mac self update


#### v0.21

- Fixed regression of brew not working in Linux derived from v0.19 changes


#### v0.20

- Fixed dnf update error handling [#13](https://github.com/hkdb/app/issues/13)


#### v0.19

- Fixed install/update scripts for M1/2/3 macs [#15](https://github.com/hkdb/app/issues/15)
- Fixed environment config for macs
- Added fullpath for brew including specific path for Apple Silicon


#### v0.18

- Fixed FreeBSD & Mac path regression


#### v0.17

- Added nixos support so that app can be used to handle Flatpak, Snap, Go, Pip, and Cargo if desired.
- Note: `whereis` must be installed on nixos for this to work.


#### v0.16

- Really fixed cargo command check this time... v0.15 didn't do the trick
- Fixed distro support regression identified by @AnonJervis specifically making LinuxMint work again


#### v0.15

- Fixed cargo command check to include $HOME/.cargo/bin/cargo so that it handles non-package-manager installed cargo properly


#### v0.14

- Fixed install.sh, update.sh, getapp.sh, updateapp.sh for Suse
- Added show OS during install and `app -m app update`


#### v0.13

- Fixed FreeBSD native package manager check logic
- Fixed updateapp.sh for Mac and FreeBSD


#### v0.12

- Fixed .app_profile logic in getapp.sh
- Aligned UX between binary (getapp.sh/updateapp.sh) and compile (install.sh and update.sh) installation
- Removed redundant header from `app -m app update`
- Added latest version check to `app -m app update` so that it doesn't update when not needed 
- Added copy button and latest release version to site


#### v0.11

- Changed bash scripts headers to allow for ./ execution on FreeBSD
- Added Update script for easy installation with less dependencies
- Implemented self-update `app -m app update`

- Launched https://hkdb.github.io/app


#### v0.10

- Added autoremove to Flatpak - Thanks to @RichardFevrier
- Fixed AppImage default enable/disable value
- Added check to ensure native package manager specified is running on the right distro
- Added SUSE derivatives support


#### v0.09

- Fixed brew and flatpak restore - Thanks to @RichardFevrier
- Fixed default enabled values to avoid errors on first launch - Thanks to @RichardFevrier
- Output history in alphabetical order by default and allow users to use the -raw flag to output history in the order packages were installed
- Fixed -y behavior for flatpak

- Development Feature: Added pre-commit linter - Thanks to @RicherFevrier
- Development Feature: Auto-compile binaries on release in preparation for better distribution/installation solution - Thanks to @RicherFevrier


#### v0.08

- Fixed FreeBSD handling of go, pip, and cargo
- Removed lsb_release dependancy - Special thanks to @RichardFevrier for this valuable pull request!
- Added update to pacman and yay
- Minor clean up of pacman and yay code


#### v0.07

- Added FreeBSD support
- Added LinuxMint support
- Fixed Fedora installation
- Added update action for Redhat derivatives
- Cleaned-up settings handling


#### v0.06

- Fixed flatpak and apt list to support keyword searches better.


#### v0.05 - Feb 10th, 2024

- Fixed .desktop file handler to avoid unintended read errors for AppImage installation


#### v0.04 - Feb 1st, 2024

- Allow for debian packages with a . in the name (ie. containerd.io or docker.io)
- Added go/pip/cargo to dist-upgrade all for Debian based distros
- Fixed UrlCheck to allow for installing cargo packages with short package names
- Updated `app help` to reflect go, pip, and cargo


#### v0.03 - Jan 17th, 2024

- Minor fixes
- Added go, pip, and cargo support


#### v0.02 - Jan 10th, 2024

- Added Homebrew & macOS support
- Fixed `app settings`

If you are upgrading rather than installing for the first time on Linux, after upgrading to this version, execute the following command:
```
echo "BREW = n" >> ~/.config/app/settings.conf
```

#### v0.01 - Jan 6th, 2024

- First Alpha Release

