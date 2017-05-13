#!/bin/sh

set -ex

rm `which channeler`
go install


channeler - demo/Tomate:ChanTomate | grep "Name(it" || exit 1;
channeler - demo/Tomate:ChanTomate | grep "package main" || exit 1;
channeler -p nop - demo/Tomate:ChanTomate | grep "package nop" || exit 1;

channeler - demo/Tomate:ChanTomate | grep "embed Tomate" || exit 1;
channeler - demo/*Tomate:ChanTomate | grep -F "embed *Tomate" || exit 1;

channeler - demo/Tomate:*ChanTomate | grep "embed Tomate" || exit 1;
channeler - demo/*Tomate:*ChanTomate | grep -F "embed *Tomate" || exit 1;

rm -fr gen_test
channeler demo/Tomate:gen_test/ChanTomate || exit 1;
ls -al gen_test | grep "chantomate.go" || exit 1;
cat gen_test/chantomate.go | grep "Name(it" || exit 1;
cat gen_test/chantomate.go | grep "package gen_test" || exit 1;
rm -fr gen_test

rm -fr demo/mytomate.go
go generate demo/main.go
ls -al demo | grep "mytomate.go" || exit 1;
cat demo/mytomate.go | grep "package main" || exit 1;
cat demo/mytomate.go | grep "NewMyTomate(" || exit 1;
cat demo/mytomatepointer.go | grep "NewMyTomatePointer(n" || exit 1;
go run demo/*.go | grep "Hello world!" || exit 1;
# rm -fr demo/tomates.go # keep it for demo

rm -fr demo/mytomate.go
go generate github.com/mh-cbon/channeler/demo
ls -al demo | grep "mytomate.go" || exit 1;
cat demo/mytomate.go | grep "package main" || exit 1;
go run demo/*.go | grep "Hello world!" || exit 1;
# rm -fr demo/gen # keep it for demo

# go test


echo ""
echo "ALL GOOD!"
