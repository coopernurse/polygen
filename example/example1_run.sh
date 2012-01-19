#!/bin/sh

set -e

echo "Generating code for; example1_idl.go"
polygen -c -dir=example1 example1_idl.go
rm -rf example1/foolib
mkdir example1/foolib
cp example1/java/*.java example1/foolib
cp Example1*.java example1/foolib

echo "Compiling foolib/*.java files"
javac -cp ../lib/deps/jackson-core-lgpl-1.9.3.jar:../lib/deps/jackson-mapper-lgpl-1.9.3.jar example1/foolib/*.java

echo "Starting Example1Server"
java -cp ../lib/deps/jackson-core-lgpl-1.9.3.jar:../lib/deps/jackson-mapper-lgpl-1.9.3.jar:example1 foolib.Example1Server &
server_pid=$!
echo "Server started with pid: $server_pid"

echo "Running Java client"
java -cp ../lib/deps/jackson-core-lgpl-1.9.3.jar:../lib/deps/jackson-mapper-lgpl-1.9.3.jar:example1 foolib.Example1Client

echo "Stopping server"
kill $server_pid
