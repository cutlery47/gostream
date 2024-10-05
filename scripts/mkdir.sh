#!/bin/bash

DIR=$1

if [ -d $DIR ]; then 
    exit 0
else 
    mkdir $DIR
fi
