#!/usr/bin/env bash

path=$(pwd)
echo ${path}
#export $GOPATH=${path}

cd ${path}/src

pwd

go build main.go

./main