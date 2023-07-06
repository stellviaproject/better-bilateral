#!/usr/bin/bash
GOOS=linux GOARCH=amd64 go build -o bilateral
echo builded succesfully
./bilateral -input ./input.tif -output ./output.tif -population 100 -generation 30 -log ./log.txt -parallel 10