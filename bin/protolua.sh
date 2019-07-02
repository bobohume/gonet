cd ../src/gonet/message
protoc -o client.pb  client.proto
protoc -o game.pb game.proto
protoc -o message.pb message.proto
cp client.pb ./pb/
cp game.pb ./pb/
cp message.pb ./pb/
rm -rf client.pb
rm -rf game.pb
rm -rf message.pb