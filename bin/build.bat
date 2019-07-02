cd ./../src/gonet/server
go build
copy /y server.exe .\..\..\..\bin
del server.exe
::go install

cd ./../client
go build
copy /y client.exe .\..\..\..\bin
::copy /y client.exe ./../../../bin
del client.exe
::go install
