cd ../src/gonet/message
protoc.exe -o client.pb client.proto
protoc.exe -o game.pb game.proto
protoc.exe -o message.pb message.proto
copy /y client.pb .\pb\
copy /y game.pb .\pb\
copy /y message.pb .\pb\
del client.pb
del game.pb
del message.pb