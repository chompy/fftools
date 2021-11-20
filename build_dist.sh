#!/bin/sh
sh ./compile_assets.sh
sh ./compile_daemon.sh
rm -rf dist
mkdir dist
mkdir dist/bin
cp bin/fftools_daemon.exe dist/bin
cp README.md dist/
cp LICENSE dist/
cp FFTools_ACT_Plugin.cs dist/
cp -r scripts dist/
mkdir dist/config
cp config/_app.yaml dist/config/
mkdir dist/data
touch dist/data/.placeholder
