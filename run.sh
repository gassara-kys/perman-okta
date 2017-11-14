#!/bin/bash

# load env
if [ -e ./myenv.sh ]; then
  source ./myenv.sh
fi

# build
mkdir -p ./bin
go build -o ./bin/perman-okta

# exec
chmod +x ./bin/perman-okta
./bin/perman-okta
