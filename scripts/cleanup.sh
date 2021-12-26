rm -f ../genetic-algorithm
rm -f ../genom.pb.go
rm -f ../fitness.pb.go
rm -f ./fitness_pb2.py

cd ../out/ || exit

find *.png -maxdepth 1 -type f | while read -r textfiles; do
  rm -v "$textfiles"
done

cd ../scripts/ || exit

rmdir ../out/
mkdir ../out
rm -rf ../data/logs.txt
rm -rf ../data/fitness.bin
