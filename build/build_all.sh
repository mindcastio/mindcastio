#!/bin/bash

export MINDCAST_HOME=/opt/data/build
export MINDCAST_SRC=/opt/data/build/src/github.com/mindcastio/mindcastio
export GOPATH=/opt/data/build

# update
cd $MINDCAST_SRC
git pull origin master

# build binaries
cd $MINDCAST_SRC/services/crawler
echo "Building the crawler ..."
go get && go build

cd $MINDCAST_SRC/services/indexer
echo "Building the indexer ..."
go get && go build

cd $MINDCAST_SRC/services/api
echo "Building the api service ..."
go get && go build

echo "Building tools"

cd $MINDCAST_SRC/tools/export_index
go get && go build

cd $MINDCAST_SRC/tools/research
go get && go build

echo "Building the migrations"

cd $MINDCAST_SRC/tools/migrations
go build migration_001.go

echo "Addding symbolic links"

if [ ! -L "/usr/local/bin/mindcast-crawler" ]; then
	sudo ln -s "$MINDCAST_SRC/services/crawler/crawler" /usr/local/bin/mindcast-crawler
fi

if [ ! -L "/usr/local/bin/mindcast-indexer" ]; then
	sudo ln -s "$MINDCAST_SRC/services/indexer/indexer" /usr/local/bin/mindcast-indexer
fi

if [ ! -L "/usr/local/bin/mindcast-api" ]; then
	sudo ln -s "$MINDCAST_SRC/services/api/api" /usr/local/bin/mindcast-api
fi

# @cd controller && go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/lambda-engine/lambda-cluster/version.GitCommit=$(COMMIT)" .
