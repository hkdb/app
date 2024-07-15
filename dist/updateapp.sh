#!/usr/bin/env bash

###############
# app updater #
###############

VER="v0.15"
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
elif [[ "$UNAMEM" == "arm" ]]; then
  CPUARCH="arm64"
else
  echo -e "❌️ CPU Architecture not supported... Exiting...\n"
  exit 1
fi

echo -e "✅️ Dependencies check...\n"
if [[ ! -f "/usr/bin/unzip" ]] && [[ ! -f "/usr/local/bin/unzip" ]]; then
  echo -e "\n❌️ unzip is not installed on this system. Install it and run the install command again...\n"
  exit 1
fi
if [[ ! -f "/usr/bin/curl" ]] && [[ ! -f "/usr/local/bin/curl" ]]; then
  echo -e "\n❌️ curl is not installed on this system. Install it and run the install command again...\n"
  exit 1
fi

echo -e "✅️ Making sure there's a $HOME/.local/bin...\n"
if [[ ! -d "$HOME/.local/bin" ]]; then
  echo -e "\n❌️ $HOME/.local/bin does not exist... Exiting...\n"
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

echo -e "\n${GREEN}**************"
echo -e " 💯️ COMPLETED"
echo -e "**************${NC}\n"
