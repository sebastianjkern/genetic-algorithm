# Compile proto buffers
protoc --go_out=.. genom.proto
protoc --go_out=.. fitness.proto
protoc --python_out=. fitness.proto

# Compile go files
cd ..
go build

# Run go project
./genetic-algorithm

# Run plotting script
cd scripts || exit
python3 plot_fitness.py
