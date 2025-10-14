#!/usr/bin/env bash

###############
# app updater #
###############

VER="v0.29"
CYAN='\033[0;36m'
GREEN='\033[1;32m'
NC='\033[0m' 

echo -e "\n📦️ Updating...\n"

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
elif [[ "$UNAMEM" == "arm64" ]]; then
  CPUARCH="arm64"
elif [[ "$UNAMEM" == "aarch64" ]]; then
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
CCHECK="$(whereis curl 2>&1)"
CL=${#CCHECK}
if [[ $CL -lt 6 ]]; then
  echo -e "\n❌️ curl is not installed on this system. Install it and run the install command again...\n"
  exit 1
fi

echo -e "✅️ Making sure there's a $HOME/.local/bin...\n"
if [[ ! -d "$HOME/.local/bin" ]]; then
  echo -e "\n❌️ $HOME/.local/bin does not exist... Use the install script instead... Exiting...\n"
  exit 1
fi

echo -e "⏳️ Downloading app binary...\n"
curl -L -o $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip https://github.com/hkdb/app/releases/download/$VER/app-$USEROS-$CPUARCH.zip
if [[ $? -ne 0 ]] ; then
    echo -e "\n❌️ Failed to download app binary... Exiting...\n"
    exit 1
fi

echo -e "\n💫️ Installing $VER binary...\n"
unzip -d $HOME/.local/bin/ $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip
mv $HOME/.local/bin/app-$USEROS-$CPUARCH $HOME/.local/bin/app

echo -e "\n🧹️ Clean-up...\n"
#rm $HOME/.local/bin/app-$USEROS-$CPUARCH-$VER.zip

$HOME/.local/bin/app settings > /dev/null

echo -e "\n${GREEN}**************"
echo -e " 💯️ COMPLETED"
echo -e "**************${NC}\n"
