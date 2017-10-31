#!/bin/bash

mkdir -p ./bin
go build -o ./bin/perman-okta
chmod +x ./bin/perman-okta
./bin/perman-okta
