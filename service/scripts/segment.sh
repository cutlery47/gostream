#!/bin/bash

VIDNAME=$1
VIDDIR=$2

ffmpeg -i $VIDDIR/$VIDNAME.mp4 -bsf:v h264_mp4toannexb -codec copy -hls_list_size 0 $VIDDIR/segmented/$VIDNAME/$VIDNAME.m3u8