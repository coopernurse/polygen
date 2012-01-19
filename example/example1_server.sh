#!/bin/sh

set -e

echo "Starting Example1Server"
java -cp ../lib/deps/jackson-core-lgpl-1.9.3.jar:../lib/deps/jackson-mapper-lgpl-1.9.3.jar:example1 foolib.Example1Server
