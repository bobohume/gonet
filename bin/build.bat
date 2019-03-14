cd ./../src/server
go build
cp server ./../../bin
::go install

cd ./../client
go build
cp client ./../../bin
::go install
