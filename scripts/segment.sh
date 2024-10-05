#!/bin/bash

# path to the video file
VIDPATH=$1
# path to the segment list file (manifest path)
MANPATH=$2
# path to the chunk file (chunk file template)
CHUNKPATH=$3
# segmentation interval length
SEGTIME=${SEGMENT_TIME:=2}

ffmpeg -i $VIDPATH -codec copy -f ssegment -hls_time $SEGTIME -segment_list $MANPATH -segment_list_type m3u8 $CHUNKPATH
