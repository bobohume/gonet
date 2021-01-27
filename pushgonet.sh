#!/bin/sh

cd /data/myserver

if [ ! -n "$1" ] ;then

    echo "error:Lack of parameters"

else

    git checkout $1

    echo "checkout $1 ok"

    git pull origin $1

    echo "pull Success"

    cd /data/myserver/bin
    sh stop.sh
    sh build.sh
    sh start.sh
fi