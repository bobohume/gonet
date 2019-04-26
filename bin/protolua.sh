cd ../src/gonet/message
proto -o client.pb  client.proto
proto -o game.pb game.proto
proto -o message.pb message.proto
cp client.pb ./pb/
cp game.pb ./pb/
cp message.pb ./pb/
rm -rf client.pb
rm -rf game.pb
rm -rf message.pb