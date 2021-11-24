protoc --go_out=. genom.proto
protoc --go_out=. fitness.proto
protoc --python_out=. fitness.proto
go build
./genetic-algorithm

python3 plot_fitness.py