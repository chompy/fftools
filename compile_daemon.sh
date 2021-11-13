#!/bin/sh
cd src/daemon
CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -o ../../bin/fftk_daemon.exe