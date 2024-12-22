#!/usr/bin/env bash

#################
# app installer #
#################

VER="v0.19"
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
  USEROS="macos"
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
elif [[ "$UNAMEM" == "arm64" ]]; then
  CPUARCH="arm64"
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
CCHECK="$(whereis curl)"
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
    echo -e "\nif [ -d \"$HOME/.local/bin\" ]; then\n\tPATH=\"$HOME/.local/bin:\$PATH\"\nfi" >> $SHELLPROFILE
    echo -e "\n# Added by app (https://github.com/hkdb/app) installation\nsource $SHELLPROFILE" >> $SHELLRC
fi

echo -e "⏳️ Downloading app binary...\n"
curl -L -o $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip https://github.com/hkdb/app/releases/download/$VER/app-$USEROS-$CPUARCH.zip
if [[ $? -ne 0 ]] ; then
    echo -e "\n❌️ Failed to download app binary... Exiting...\n"
    exit 1
fi

echo -e "\n💫️ Installing binary...\n"
unzip -d $HOME/.local/bin/ $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip
mv $HOME/.local/bin/app-$USEROS-$CPUARCH $HOME/.local/bin/app

echo -e "\n🧹️ Clean-up...\n"
#rm $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip

echo -e "\n${GREEN}**************"
echo -e " 💯️ COMPLETED"
echo -e "**************${NC}\n"

echo -e "⚠️  You may need to close and reopen your existing terminal windows for app to work as expected...\n"

