#!/bin/sh

GOPATH=$(go env GOPATH)
GOMOD=$(go env GOMOD)

PKG_VERSION=$(grep "$GOMOD" -e 'github.com/remyoudompheng/bigfft' | awk '{print $2}')
PKG_PATH=$GOPATH/pkg/mod/github.com/remyoudompheng/bigfft@${PKG_VERSION}

BASEDIR=$(dirname "$0")

patch () {
  cp $BASEDIR/arith_wasm.s $PKG_PATH/
}

unpatch () {
  rm $PKG_PATH/arith_wasm.s
}

case $1 in
patch)
  patch
  ;;
unpatch)
  unpatch
  ;;
esac
