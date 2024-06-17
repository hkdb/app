#!/bin/bash

echo -e "\nInstalling app...\n\n"

echo -e "Make sure there's a $HOME/.local...\n"

if [[ ! -d $HOME/.local ]]; then
  mkdir -p $HOME/.local/bin
fi

echo -e "Make sure there's a $HOME/.local/bin...\n"

if [[ ! -d $HOME/.local/bin ]]; then
  mkdir -p $HOME/.local/bin
fi

echo -e "Check to make sure that $HOME/.local/bin is part of PATH...\n"

PROFILE=$(cat ~/.profile |grep .local/bin)
if [[ "$OSTYPE" == "freebsd"* ]]; then
  PROFILE=$(cat ~/.zshrc |grep .local/bin)
fi

if [[ "$OSTYPE" == "freebsd"* ]]; then
  if [[ ! -f "$HOME/.zsh_profile" ]]; then
    echo -e "\nCreating .zshrc_profile to have ~/.local/bin be part of PATH and inserting lines to .zshrc to source it...\n"
    echo -e "\nif [ -d \"$HOME/.local/bin\" ]; then\n\tPATH=\"$HOME/.local/bin:\$PATH\"\nfi" >> $HOME/.zsh_profile
    echo -e "\nsource ~/.zsh_profile" >> $HOME/.zshrc
  fi
else 
  if [[ "$PROFILE" == "" ]]; then
    echo -e "\nAdding lines in .profile to have ~/.local/bin be part of PATH...\n"
    echo -e '\nif [ -d "$HOME/.local/bin" ]; then\n\tPATH="$HOME/.local/bin:$PATH"\nfi' >> $HOME/.profile
  fi
fi

if [[ "$OSTYPE" == "freebsd"* ]]; then
  echo -e "\nSourcing .zsh_profile to ensure that ~/.local/bin is in PATH...\n"
  source $HOME/.zsh_profile
else
  echo -e "\nSourcing .profile to ensure that ~/.local/bin is in PATH...\n"
  source $HOME/.profile
fi
echo -e "\n"

DISTRO=""
PKGMGR=""
IFLAG=""

if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "linux" ]]; then
  
  OSR=$(cat /etc/os-release |grep ^NAME)
  OSRNV=${OSR:5}
  DIST=${OSRNV:1:${#OSRNV}-2}

  echo -e "\nOS: Linux\n"
  DISTRO=$(cat /etc/*-release | grep "^ID=" | head -1 | cut -d '=' -f 2)

  if [[ "$DISTRO" == "\"opensuse\"" ]] || [[ "$DISTRO" == "\"opensuse-leap\"" ]] || [[ "$DISTRO" == "\"suse\"" ]]; then
    DISTRO="opensuse-leap"
  fi

  if [[ "$DISTRO" == "debian" ]] || [[ "$DISTRO" == "ubuntu" ]] || [[ "$DISTRO" == "pop" ]] || [[ "$DISTRO" == "linuxmint" ]]; then
    PKGMGR="apt"
    IFLAG="install"
  elif [[ "$DISTRO" == "fedora" ]] || [[ "$DISTRO" == "rocky" ]] || [[ "$DISTRO" == "almalinux" ]] || [[ "$DISTRO" == "centos" ]] || [[ "$DISTRO" == "RedHatEnterpriseServer" ]] || [[  "$DISTRO" == "ol" ]] || [[ "$DISTRO" == "clear-linux-os" ]] || [[ "$DISTRO" == "AmazonAMI" ]]; then
    PKGMGR="dnf"
    IFLAG="install"
  elif [[ "$DISTRO" == "arch" ]] || [[ "$DISTRO" == "garuda" ]] || [[ "$DISTRO" == "manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
    PKGMGR="pacman"
    IFLAG="-S"
  elif [[ "$DISTRO" == "opensuse" ]] || [[ "$DISTRO" == "opensuse-leap" ]] || [[ "$DISTRO" == "suse" ]]; then
    PKGMGR="zypper"
    IFLAG="install"
  fi

  read -p "Would you like to install Flatpak? (Y/n) " FLATPAK
  if [[ $FLATPAK != "N" ]] && [[ $FLATPAK != "n" ]]; then
   sudo $PKGMGR $IFLAG flatpak
   if [[ "$DISTRO" == "arch" ]] || [[ "$DISTRO" == "garuda" ]] || [[ "$DISTRO" == "manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
    U=$USER
    sudo adduser $U _flatpak
    sudo flatpak repair
   fi
   flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
  fi

  read -p "Would you like to install Snap? (Y/n) " SNAP
  if [[ $SNAP != "N" ]] && [[ $SNAP != "n" ]]; then
    if [[ "$DISTRO" == "arch" ]] || [[ "$DISTRO" == "garuda" ]] || [[ "$DISTRO" == "manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
      OPATH=$(pwd)
      cd /tmp/
      git clone https://aur.archlinux.org/snapd.git
      cd snapd
      makepkg -si
      sudo systemctl enable --now snapd.socket
      sudo ln -s /var/lib/snapd/snap /snap
      echo "export PATH=\$PATH:\/snap/bin/" | sudo tee -a /etc/profile
      source /etc/profile
      cd $OPATH
    else
      sudo $PKGMGR $IFLAG snapcraft
    fi
  fi

  echo -e "\n"
  if [[ "$DISTRO" == "linuxmint" ]] || [[ "$DISTRO" == "debian" ]]; then
    read -p  "Add software-properties-common? (Y/n) " SPC
    if [[ "$SPC" != "N" ]] && [[ "$SPC" != "n" ]]; then
			sudo apt install software-properties-common
    fi
  fi
   
  echo -e "\n"
  if [[ "$DISTRO" == "debian" ]] || [[ "$DISTRO" == "ubuntu" ]] || [[ "$DISTRO" == "pop" ]]; then
    read -p  "Add longsleep-ubuntu-golang-backports? (Y/n) " BACKPORTS
    if [[ "$BACKPORTS" != "N" ]] && [[ "$BACKPORTS" != "n" ]]; then
      # sudo echo -e "deb https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main\n# deb-src https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main" |sudo tee /etc/apt/sources.list.d/longsleep-ubuntu-golang-backports-jammy.list
      # sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys F6BC817356A3D45E
			sudo add-apt-repository ppa:longsleep/golang-backports
      sudo apt update -y
    fi
  fi
   
  read -p  "Install Go? Only say no if you already have it installed... (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
		if [[ "$DISTRO" == "fedora" ]] || [[ "$DISTRO" == "rocky" ]] || [[ "$DISTRO" == "almalinux" ]] || [[ "$DISTRO" == "centos" ]] || [[ "$DISTRO" == "RedHatEnterpriseServer" ]] || [[  "$DISTRO" == "ol" ]] || [[ "$DISTRO" == "clear-linux-os" ]] || [[ "$DISTRO" == "AmazonAMI" ]]; then
			wget -O /tmp/go1.22.4.linux-amd64.tar.gz https://go.dev/dl/go1.22.4.linux-amd64.tar.gz  
			sudo tar -C /usr/local -xzf /tmp/go1.22.4.linux-amd64.tar.gz
			export PATH=$PATH:/usr/local/go/bin
			echo -e "if [[ -d /usr/local/go/bin ]]; then\n\tPATH=$PATH:/usr/local/go/bin\nfi" >> $HOME/.bashrc
    elif [[ "$DISTRO" == "opensuse" ]] || [[ "$DISTRO" == "opensuse-leap" ]] || [[ "$DISTRO" == "suse" ]]; then
      sudo $PKGMGR $IFLAG go1.22
		else
    	sudo $PKGMGR $IFLAG golang
		fi
  fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
  echo -e "\nOS: macos\n"
  PKGMGR="brew"
  read -p "Install Homebrew? (Y/n) " HOMEBREW
  if [[ "$HOMEBREW" != "N" ]] && [[ "$HOMEBREW" != "n" ]]; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
  read -p "Install Go? Only say no if you already have it installed (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
    echo -e "Installing Go...\n"
    brew install golang
  fi
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  read -p "Install Go? Only say no if you already have it installed (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
    echo -e "Installing Go...\n"
    sudo pkg install golang
  fi
fi

go mod tidy
echo -e "Compiling app...\n"
go build
echo -e "Copying app to $HOME/.local/bin...\n"
cp app $HOME/.local/bin/
echo -e "\n********"
echo -e "COMPLETE"
echo -e "********\n"

echo -e "You will need to logout and log back in to ensure that app is in PATH...\n"
