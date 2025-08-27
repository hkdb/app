#!/usr/bin/env bash

#################
# app installer #
#    COMPILE    #
#################

VER="v0.12"
CYAN='\033[0;36m'
GREEN='\033[1;32m'
NC='\033[0m' 

echo -e "\n📦️ Installing:"

echo -e "${CYAN}
_____  ______ ______  
\__  \ \____ \\____  \ 
 / __ \|  |_> >  |_> >
(____  /   __/|   __/ 
     \/|__|   |__|    
${NC}"

echo -e "🚀️ ${GREEN}The Cross-Platform Package Management Assistant with Super Powers${NC}\n"

USEROS=""
echo -e "🐧️ Detecting OS..."
if [[ "$OSTYPE" == "linux"* ]]; then
  USEROS="linux"
  echo -e "\n🐧️ Linux\n"
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  USEROS="freebsd"
  echo -e "\n🅱️  FreeBSD\n"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  USEROS="darwin"
  echo -e "\n🍎️ MacOS"
else
  echo -e "❌️ Operating System not supported... Exiting...\n"
  exit 1
fi

echo -e "💻️ Detecting CPU arch...\n"

CPUARCH=""
UNAMEM=$(uname -m)
echo -e "🏰️: $UNAMEM\n"

if [[ "$UNAMEM" == "x86_64" ]] || [[ "$UNAMEM" == "amd64" ]]; then
  CPUARCH="amd64"
elif [[ "$UNAMEM" == "arm" ]]; then
  CPUARCH="arm64"
elif [[ "$UNAMEM" == "aarch64" ]]; then
  CPUARCH="aarch64"
else
  echo -e "❌️ CPU Architecture not supported... Exiting...\n"
  exit 1
fi

echo -e "✅️ Dependencies check...\n"
UCHECK="$(whereis unzip)"
UL=${#UCHECK}
if [[ $UL -lt 7 ]]; then
  echo -e "\n❌️ unzip is not installed on this system. Install it and run the install command again...\n"
  exit 1
fi

CCHECK="$(whereis curl 2>&1)"
CL=${#CCHECK}
if [[ $CL -lt 6 ]]; then
  echo -e "\n❌️ curl is not installed on this system. Install it and run the install command again...\n"
  exit 1
fi

echo -e "✅️ Detecting shell...\n"

SHELLTYPE=$(basename ${SHELL})

echo -e "🐚️ shell: $SHELLTYPE"

SHELLRC="none"
SHELLPROFILE="$HOME/.config/app/.app_profile"

if [[ $SHELLTYPE == "sh" ]]; then
  SHELLRC="$HOME/.shrc"
fi

if [[ $SHELLTYPE == "csh" ]]; then
  SHELLRC="$HOME/.cshrc"
fi

if [[ $SHELLTYPE == "ksh" ]]; then
  SHELLRC="$HOME/.kshrc"
fi

if [[ $SHELLTYPE == "tcsh" ]]; then
  SHELLRC="$HOME/.tcshrc"
fi

if [[ $SHELLTYPE == "bash" ]]; then
  SHELLRC="$HOME/.bashrc"
fi

if [[ $SHELLTYPE == "zsh" ]]; then
  SHELLRC="$HOME/.zshrc"
fi

if [[ $SHELLTYPE == "fish" ]]; then
  SHELLRC="$HOME/.config/fish/config.fish"
fi

if [[ $SHELLRC == "none" ]]; then
  echo -e "\n❌️ Unrecognized shell... app only supports sh, csh, ksh, tcsh, bash, zsh, and fish... exiting...\n"
  exit 1
fi

echo -e "🐚️ config: $SHELLRC\n"

echo -e "✅️ Create app config dir if not already created...\n"
if [[ ! -d "$HOME/.config/app" ]]; then
  mkdir -p $HOME/.config/app
  if [[ $? -ne 0 ]] ; then
      echo -e "\n❌️ Failed to create $HOME/.config/app... Exiting...\n"
      exit 1
  fi
fi

echo -e "✅️ Making sure there's a $HOME/.local/bin...\n"
if [[ ! -d "$HOME/.local/bin" ]]; then
  mkdir -p $HOME/.local/bin
  if [[ $? -ne 0 ]] ; then
      echo -e "\n❌️ Failed to create $HOME/.local/bin... Exiting...\n"
      exit 1
  fi
fi

echo -e "✅️ Making sure $HOME/.local/bin is in PATH...\n"
if [[ -f $SHELLPROFILE ]]; then
  PCHECK=$(grep ".local/bin" $SHELLPROFILE)
  if [[ "$PCHECK" == "" ]]; then
    echo -e "\nif [ -d \"$HOME/.local/bin\" ]; then\n\tPATH=\"$HOME/.local/bin:\$PATH\"\nfi" >> $SHELLPROFILE
    echo -e "\n# Added by app (https://github.com/hkdb/app) installation\nsource $SHELLPROFILE" >> $SHELLRC
  fi
else
    if [[ $SHELLTYPE == "fish" ]]; then
      echo -e "if test -d \"$HOME/.local/bin\"\n   set -U fish_user_paths $HOME/.local/bin \$PATH\nend" >> $SHELLPROFILE
    else
      echo -e "\nif [ -d \"$HOME/.local/bin\" ]; then\n\tPATH=\"$HOME/.local/bin:\$PATH\"\nfi" >> $SHELLPROFILE
    fi
    echo -e "\n# Added by app (https://github.com/hkdb/app) installation\nsource $SHELLPROFILE" >> $SHELLRC
fi

DISTRO=""
PKGMGR=""
IFLAG=""


if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "linux" ]]; then
  
  echo -e "🐧️ Detecting distro...\n"

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
 
  echo -e "\n🛸️ Distro: $DISTRO"
  echo -e "📦️ Native: $PKGMGR\n"

  if [[ "$DISTRO" == "debian" ]] || [[ "$DISTRO" == "ubuntu" ]] || [[ "$DISTRO" == "pop" ]]; then
    read -p  "❔️ Add longsleep-ubuntu-golang-backports? (Y/n) " BACKPORTS
    if [[ "$BACKPORTS" != "N" ]] && [[ "$BACKPORTS" != "n" ]]; then
      # sudo echo -e "deb https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main\n# deb-src https://ppa.launchpadcontent.net/longsleep/golang-backports/ubuntu/ jammy main" |sudo tee /etc/apt/sources.list.d/longsleep-ubuntu-golang-backports-jammy.list
      # sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys F6BC817356A3D45E
			sudo add-apt-repository ppa:longsleep/golang-backports
      sudo apt update -y
    fi
  fi

  if [[ $DISTRO == "nixos" ]]; then
    GCHECK="$(whereis go)"
    GL=${#GCHECK}
    if [[ $GL -lt 3 ]]; then
      echo -e "\n❌️ Golang is not installed on this system. Install it and run the install command again...\n"
      exit 1
    fi
  else 
    read -p  "❔️ Install Go? Only say no if you already have it installed... (Y/n) " GOLANG
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
  fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
  echo -e "\nOS: macos\n"
  PKGMGR="brew"
  read -p "❔️ Install Homebrew? (Y/n) " HOMEBREW
  if [[ "$HOMEBREW" != "N" ]] && [[ "$HOMEBREW" != "n" ]]; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
  read -p "❔️ Install Go? Only say no if you already have it installed (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
    echo -e "Installing Go...\n"
    brew install golang
  fi
elif [[ "$OSTYPE" == "freebsd"* ]]; then
  read -p "❔️ Install Go? Only say no if you already have it installed (Y/n) " GOLANG
  if [[ "$GOLANG" != "N" ]] && [[ "$GOLANG" != "n" ]]; then
    echo -e "Installing Go...\n"
    sudo pkg install golang
  fi
fi

echo -e "\n🛰️ Getting modules...\n"
go mod tidy
echo -e "🛠️ Compiling app...\n"
go build
echo -e "💾️ Copying app to $HOME/.local/bin...\n"
cp app $HOME/.local/bin/



echo -e "\n${GREEN}**************"
echo -e " 💯️ COMPLETED"
echo -e "**************${NC}\n"

echo -e "⚠️  You may need to close and reopen your existing terminal windows for app to work as expected...\n"



