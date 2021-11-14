#!/bin/sh
cd src/asset_compile
go build -o ../../bin/asset_compile.exe
cd ../../
./bin/asset_compile.exe