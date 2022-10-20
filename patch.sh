#!/bin/sh

set -e

export GOPATH=$(pwd)/.go

./patches/bigfft/script.sh $1
