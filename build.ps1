if ( Test-Path release ) {
    rm -Recurse -Force release | Out-Null
}

$Env:CGO_ENABLED='0'
$Env:GOROOT_FINAL='/usr'

$Env:GOOS='linux'
$Env:GOARCH='amd64'
go build -a -trimpath -asmflags '-s -w' -ldflags '-s -w' -o release\stream
if ( -Not $? ) { exit $lastExitCode }

cp -Force default.json release
exit 0
