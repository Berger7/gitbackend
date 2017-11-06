#!/usr/bin/env bash

TEST_GIT_PATH="./git-test.git"
if [ ! -x "$TEST_GIT_PATH" ]; then
    git clone https://github.com/Berger7/git-test.git --bare
fi