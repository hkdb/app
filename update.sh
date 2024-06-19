#!/usr/bin/env bash

echo -e "\n📦️ Updating...\n"

echo -e "🛰️  Getting modules...\n"
go mod tidy
echo -e "🛠️  Compiling app...\n"
go build
echo -e "💾️ Copying app to $HOME/.local/bin...\n"
cp app $HOME/.local/bin/

echo -e "\n${GREEN}**************"
echo -e " 💯️ COMPLETED"
echo -e "**************${NC}\n"
