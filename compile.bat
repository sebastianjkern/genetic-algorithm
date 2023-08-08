protoc --go_out=./ ./genom.proto
protoc --go_out=./ ./fitness.proto
protoc --python_out=./python --pyi_out=./python fitness.proto

go build .

./genetic-algorithm

