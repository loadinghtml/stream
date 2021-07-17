Remove-Item -Force go.*

go mod init
go mod tidy
exit $?