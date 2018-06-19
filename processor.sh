#!/bin/bash

RESULTS_FILENAME='result_.json'
RESULTS_DIRECTORY='results'

# https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
RED='\033[0;31m'
LIGHT_GREEN='\033[1;32m'
BROWN='\033[0;33m'
NC='\033[0m' # No Color

# https://stedolan.github.io/jq/manual/

which jq 1>/dev/null 2>& 1
if [ $? == 0 ] ; then
  >/dev/null pushd $RESULTS_DIRECTORY
  printf "${RED}"
  jq -M ' .examples[] | select(.status != "passed")| .description,.status' $RESULTS_FILENAME
  printf "${NC}"
  printf "${LIGHT_GREEN}Summary: ${BROWN}"
  jq -M '.summary_line' $RESULTS_FILENAME
  printf "${NC}"
  >/dev/null popd
fi