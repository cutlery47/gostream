#!/bin/bash

DIRPATH=$1

if [ -d $DIRPATH ]; then 
    exit 0
else 
    mkdir $DIRPATH
    exit 1
fi
