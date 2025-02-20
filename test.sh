#!/bin/sh

set -e

for i in {1..2000}; do
    echo $i
    bash run.sh
    clear
done
