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

if [[ PROFILE == "" ]]; then
  echo -e 'if [ -d "$HOME/.local/bin" ]; then\n\tPATH="$HOME/.local/bin:$PATH"\nfi' >> $HOME/.profile
fi

source $HOME/.profile

DISTRO=""
PKGMGR=""
IFLAG=""

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  echo -e "\nOS: Linux\n"
  DISTRO=$(lsb_release -i -s)

  if [[ "$DISTRO" == "Debian" ]] || [[ "$DISTRO" == "Ubuntu" ]] || [[ "$DISTRO" == "Pop" ]] || [[ "$DISTRO" == "Linuxmint" ]]; then
    PKGMGR="apt"
    IFLAG="install"
  elif [[ "$DISTRO" == "Fedora" ]] || [[ "$DISTRO" == "Rocky" ]] || [[ "$DISTRO" == "AlmaLinux" ]] || [[ "$DISTRO" == "CentOS" ]] || [[ "$DISTRO" == "RedHatEnterpriseServer" ]] || [[  "$DISTRO" == "Oracle" ]] || [[ "$DISTRO" == "ClearOS" ]] || [[ "$DISTRO" == "AmazonAMI" ]]; then
    PKGMGR="dnf"
    IFLAG="install"
  elif [[ "$DISTRO" == "Arch" ]] || [[ "$DISTRO" == "Garuda" ]] || [[ "$DISTRO" == "Manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
    PKGMGR="pacman"
    IFLAG="-S"
  fi

  read -p "Would you like to install Flatpak? (Y/n) " FLATPAK
  if [[ $FLATPAK != "N" ]] && [[ $FLATPAK != "n" ]]; then
   sudo $PKGMGR $IFLAG flatpak
   if [[ "$DISTRO" == "Arch" ]] || [[ "$DISTRO" == "Garuda" ]] || [[ "$DISTRO" == "Manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
    U=$USER
    sudo adduser $U _flatpak
    sudo flatpak repair
   fi
   flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
  fi

  read -p "Would you like to install Snap? (Y/n) " SNAP
  if [[ $SNAP != "N" ]] && [[ $SNAP != "n" ]]; then
    if [[ "$DISTRO" == "Arch" ]] || [[ "$DISTRO" == "Garuda" ]] || [[ "$DISTRO" == "Manjaro" ]] || [[ "$DISTRO" == "Endeavour" ]]; then
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
  if [[ "$DISTRO" == "Debian" ]] || [[ "$DISTRO" == "Ubuntu" ]] || [[ "$DISTRO" == "Pop" ]]; then
    read -p  "Add longsleep-ubuntu-golang-backports? (Y/n) " BACKPORTS
    if [[ "$BACKPORTS" != "N" ]] && [[ "$BACKPORTS" != "n" ]]; then
      # sudo echo -e "deb https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main\n# deb-src https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main" |sudo tee /etc/apt/sources.list.d/longsleep-ubuntu-golang-backports-jammy.list
      # sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys F6BC817356A3D45E
			sudo add-apt-repository ppa:longsleep/golang-backports
      sudo apt update -y
    fi
  fi
   
  read -p  "Install Go 1.21? Only say no if you already have it installed... (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
		if [[ "$DISTRO" == "Fedora" ]] || [[ "$DISTRO" == "Rocky" ]] || [[ "$DISTRO" == "AlmaLinux" ]] || [[ "$DISTRO" == "CentOS" ]] || [[ "$DISTRO" == "RedHatEnterpriseServer" ]] || [[  "$DISTRO" == "Oracle" ]] || [[ "$DISTRO" == "ClearOS" ]] || [[ "$DISTRO" == "AmazonAMI" ]]; then
			wget -O /tmp/go1.21.5.linux-amd64.tar.gz https://go.dev/dl/go1.21.5.linux-amd64.tar.gz  
			sudo tar -C /usr/local -xzf /tmp/go1.21.5.linux-amd64.tar.gz
			export PATH=$PATH:/usr/local/go/bin
			echo -e "if [[ -d /usr/local/go/bin ]]; then\n\tPATH=$PATH:/usr/local/go/bin\nfi" >> $HOME/.bashrc
		else
    	sudo $PKGMGR $IFLAG golang
		fi
  fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
  echo -e "\nOS: macos\n"
  PKGMGR="brew"
  read -p "Install Homebrew? (Y/n)" HOMEBREW
  if [[ "$HOMEBREW" != "N" ]] && [[ "$HOMEBREW" != "n" ]]; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
  read -p "Install Go? Only say no if you already have it installed (Y/n)" GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
    echo -e "Installing Go...\n"
    brew install golang
  fi
fi

go mod tidy
echo -e "Compiling app...\n"
go build
echo -e "Copying app to $HOME/.local/bin...\n"
cp app $HOME/.local/bin/
echo -e "COMPLETE\n"

