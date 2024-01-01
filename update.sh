#!/bin/bash

go mod tidy
go build
cp app $HOME/.local/bin/
