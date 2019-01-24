#!/bin/bash

[ -z "$GOPATH" ] && export GOPATH=$HOME/go

[ ! -d "$GOPATH/src/gitlab.com/diamondburned/shittercord" ] && {
	go get gitlab.com/diamondburned/shittercord
}

export LD_LIBRARY_PATH="$GOPATH/src/gitlab.com/diamondburned/shittercord"

case $1 in 
run)
	$GOPATH/bin/shittercord $@
	;;
install|update)
	go get -u -v gitlab.com/diamondburned/shittercord
	;;
*)
	echo "Usage: ./run.sh [command]"
	echo $'\t- run:    runs the application'
	echo $'\t- update: updates the application'
	;;
esac

