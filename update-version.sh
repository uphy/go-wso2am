#!/bin/bash

if [ $# != 1 ]; then
  echo Specify the version.
  exit 1
fi

sed -i "" -e "s/const Version = .*/const Version = \"$1\"/" cli/cli.go