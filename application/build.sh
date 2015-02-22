export GOPATH=`pwd`

cd src

go clean -i -r
go build -o ../poc.exe