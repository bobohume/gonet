#!/bin/sh

if [ ! -n "$1" ] ;then

    echo "error:Lack of parameters"

else
    ssh -p 22 -t root@192.168.101.49  "/data/push.sh $1"
fi