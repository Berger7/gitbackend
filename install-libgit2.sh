#!/usr/bin/env bash

set -ex

wget https://github.com/libgit2/libgit2/archive/v0.26.0.tar.gz -O /tmp/libgit2.tar.gz
cd /tmp
tar -xvf /tmp/libgit2.tar.gz
cd ./libgit2-0.26.0
mkdir build && cd build
cmake ..
cmake .. -DCMAKE_INSTALL_PREFIX=/usr/
sudo cmake --build . --target install