#!/bin/bash

RESULTS_FILENAME='result_.json'
RESULTS_DIRECTORY='results'

# https://stackoverflow.com/questions/911168/how-to-detect-if-my-shell-script-is-running-through-a-pipe
if [ -t 1 ] ; then
  # terminal
  # https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
  RED='\033[0;31m'
  LIGHT_GREEN='\033[1;32m'
  BROWN='\033[0;33m'
  NC='\033[0m' # no color
else
  # not a terminal
  RED=''
  LIGHT_GREEN=''
  BROWN=''
  NC=''
fi

# https://stedolan.github.io/jq/manual/

which jq 1>/dev/null 2>& 1
if [ $? == 0 ] ; then
  >/dev/null pushd $RESULTS_DIRECTORY
  printf "${RED}"
  # https://github.com/stedolan/jq/issues/785
  # https://stackoverflow.com/questions/39139107/how-to-format-a-json-string-as-a-table-using-jq
  jq -M ' .examples[] | select(.status != "passed")| "\(.full_description)\n"' $RESULTS_FILENAME | awk '{print $0}' | sed 's/Command.*$//'
  printf "${NC}"
  printf "${LIGHT_GREEN}Summary: ${BROWN}"
  jq -M '.summary_line' $RESULTS_FILENAME
  printf "${NC}"
  >/dev/null popd
fi