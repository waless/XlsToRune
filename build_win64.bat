@echo off
pushd .\src

set GOOS=windows
set GOARCH=amd64
go build -o .\..\build\win64\XlsToRune.exe

popd
