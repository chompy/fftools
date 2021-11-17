#!/bin/sh
sh ./compile_assets.sh
sh ./compile_daemon.sh
rm -rf dist
mkdir dist
cp bin/fftools_daemon.exe dist/
cp README.md dist/
cp LICENSE dist/
cp FFTools_ACT_Plugin.cs dist/
cp -r scripts dist/
mkdir dist/config
cp config/_app.yaml dist/config/
mkdir dist/data
touch dist/data/.placeholder
sed -i 's/bin\\\\fftools_daemon.exe/fftools_daemon.exe/g' dist/FFTools_ACT_Plugin.cs