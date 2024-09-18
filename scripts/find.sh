#!/bin/bash

VIDPATH=$1

if [ -f $VIDPATH ]; then
    exit 0
else 
    exit 1
fi