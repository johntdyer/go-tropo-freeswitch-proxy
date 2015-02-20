#!/bin/bash
CGO_ENABLED=0
VERSION=`cat version.go | grep "const Version string" | cut -d'"' -f2`

while getopts 'uh' OPT; do
  case $OPT in
    h)  hlp="yes";;
    u)  upload="yes";;
  esac
done

# usage
HELP="
    usage: $0 [ -u -h ]

        -u --> Uploads versioned linux build to s3 artifacts.voxeolabs.net/tropo-auth
        -h --> print this help screen
"

if [ "$hlp" = "yes" ]; then
  echo "$HELP"
  exit 0
fi

# Install dependencies
# Runtime dependencies
go get github.com/tools/godep

BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
# Build for OSX
godep go build -ldflags "-w -X main.buildDate ${BUILD_DATE}" -o builds/tropo-auth-$VERSION.osx
if [ $? -eq 0 ]; then
  echo "Success Build artifact - builds/tropo-auth-$VERSION.osx"
else
  echo "Build error"
  exit $?
fi

# Build for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -ldflags "-w -X main.buildDate ${BUILD_DATE}" -o builds/tropo-auth-$VERSION.linux
if [ $? -eq 0 ]; then
  echo "Success Build artifact - builds/tropo-auth-$VERSION.linux"
else
  echo "Build error"
  exit $?
fi

if [ "$upload" = "yes" ]; then
  echo "Uploading builds/tropo-auth-$VERSION.linux to s3::artifacts.voxeolabs.net/tropo-auth/tropo-auth-$VERSION.linux"
  /usr/local/bin/aws s3 cp ./builds/tropo-auth-$VERSION.linux s3://artifacts.voxeolabs.net/tropo-auth/tropo-auth-$VERSION.linux --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers > /dev/null
  if [ $? -eq 0 ]; then
    echo "Succesfully Uploaded artifact - tropo-auth-$VERSION.linux -- http://artifacts.voxeolabs.net.s3.amazonaws.com/tropo-auth/tropo-auth-$VERSION.linux"
  else
    echo "Upload error"
    exit $?
  fi
fi
