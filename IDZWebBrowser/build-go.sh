#!/bin/sh

#  build-go.sh
#  GoTest
#
#  Created by Saber Gilani on 02/07/15.
#  Copyright (c) 2015 DevCups SL. All rights reserved.

#####################################
# Change the paths to your porjects #
#####################################

set -e

GO=/Users/sgilani/go1.4/ios/bin/go

GG_OBJ=/Users/sgilani/Documents/WebBrowser/IDZWebBrowser/go-obj
GG_CGO_OBJ=/Users/sgilani/Documents/WebBrowser/IDZWebBrowser/cgo-obj
ARCHIVE=/Users/sgilani/Documents/WebBrowser/IDZWebBrowser/go-output.a

export GOPATH=$GOPATH/Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing/syncthing

mkdir -p $GG_OBJ
mkdir -p $GG_CGO_OBJ

cd /Users/sgilani/Documents/WebBrowser/IDZWebBrowser/syncthing

echo "Done 1"

#CGO_ENABLED=1 GOARCH=arm GOARM=7 $GO tool cgo -objdir $GG_CGO_OBJ main.go
CGO_ENABLED=1 GOARCH=arm64 GOARM=7 $GO tool cgo -objdir ../cgo-obj cmd/syncthing/main.go
cp $GG_CGO_OBJ/_cgo_export.h /Users/sgilani/Documents/WebBrowser/IDZWebBrowser

echo "Done 2 GO TOOL"

$GO run build.go clean
$GO run build.go assets

echo "Done 3 GO BUILD ASSETS"

#CGO_ENABLED=1 GOARCH=arm GOARM=7 $GO build -ldflags '-tmpdir ../go-obj -linkmode external'
$GO run build.go -goos darwin -goarch arm64 build

echo "Done 4 GO BUILD"

ar rcs $ARCHIVE $GG_OBJ/*.o
echo "Generated $ARCHIVE"