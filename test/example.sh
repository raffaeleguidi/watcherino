#!/bin/sh
# the script receives three arguments: folder, file and CREATE|WRITEPARAM_FOLDER=$1
PARAM_FILE=$2
PARAM_EVENT=$3
echo "hi, this is watcherino on folder '${PARAM_FOLDER}' file '${PARAM_FILE}' with event '${PARAM_EVENT}'" >> test.log