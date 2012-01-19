#!/bin/sh

set -e

java -cp ../lib/deps/jackson-core-lgpl-1.9.3.jar:../lib/deps/jackson-mapper-lgpl-1.9.3.jar:example1 foolib.Example1Client
