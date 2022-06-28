#!/bin/sh

cd "${0%/*}"
go build -o dist/most-profitable-pool src/main/main.go