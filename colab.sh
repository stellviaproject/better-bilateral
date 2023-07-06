#!/usr/bin/bash
GOOS=linux GOARCH=amd64 go build -o bilateral
echo builded succesfully
./bilateral -input ./input.tif -output ./output.tif -population 10000 -generation 100 -log ./log.txt -parallel 1000