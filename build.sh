#!/usr/bin/env bash

mkdir -p build
#clean up
rm build/*
##build the portal
go build -o build/ScriptPortal scriptPortal.go
##build the plugins
if [ -f plugins/*.go ]
then
	cd plugins
	go get -d ./...
	cd ..
	for f in $(ls plugins/*.go)
	do
		go build -buildmode=plugin -o "build/$(basename $f .go).so" $f
	done
	"build finished"
else
	echo "no plugins found, build finished"
fi

