rsrc -manifest rsc/ogglooper.exe.manifest -ico rsc/maki.ico -o rsrc.syso
go build -ldflags -H=windowsgui -o GMRepeater.exe
