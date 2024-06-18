#!/usr/bin/env bash

go mod tidy
go build
cp app $HOME/.local/bin/
