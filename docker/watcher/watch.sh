#!/bin/bash

N=${1:-4}
IPBASE=${2:-172.77.5.}
PORT=${3:-80}

watch -t -n 1 '
for i in $(seq 1 '$N');
do
    curl -s -m 1 http://'$IPBASE'$i:'$PORT'/Stats | \
        tr -d "{}\"" | \
        awk -F "," '"'"'{gsub (/[,]/," "); print;}'"'"'
done;
'
