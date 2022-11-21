#!/bin/sh

builddate() {
	date +%FT%T%z
}

version() {
	git describe --tags --always --dirty --match=v* 2> /dev/null || \
		cat .version 2> /dev/null || echo v0
}

builddate > builddate.txt
version > version.txt
