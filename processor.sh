#!/bin/bash

WORKDIR=${1:-/uru}
RESULTS_BASENAME=${2:-result}

RESULTS_FILENAME="${RESULTS_BASENAME}_.json"
RESULTS_DIRECTORY="${RESULTS_BASENAME}s"

# https://stackoverflow.com/questions/911168/how-to-detect-if-my-shell-script-is-running-through-a-pipe
if [ -t 1 ] ; then
  # terminal
  # https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
  RED='\033[0;31m'
  LIGHT_RED='\033[1;31m'
  LIGHT_GREEN='\033[1;32m'
  BROWN='\033[0;33m'
  NC='\033[0m' # no color
else
  # not a terminal
  RED=''
  LIGHT_GREEN=''
  LIGHT_RED=''
  BROWN=''
  NC=''
fi

# https://stedolan.github.io/jq/manual/

which jq 1>/dev/null 2>& 1
if [ $? != 0 ] ; then
  exit 1
fi
if [ ! -d $WORKDIR  ] ; then
  echo "Invalid work directory: \"${WORKDIR}\""
  exit 1
fi
>/dev/null pushd $WORKDIR
if [ ! -d $RESULTS_DIRECTORY  ] ; then
  echo "Invalid results directory: \"${WORKDIR}/${RESULTS_DIRECTORY}\""
  exit 1
fi
>/dev/null pushd $RESULTS_DIRECTORY
# https://github.com/stedolan/jq/issues/785
# https://stackoverflow.com/questions/39139107/how-to-format-a-json-string-as-a-table-using-jq
jq -M '.examples[] | select(.status != "passed")| (.status)' $RESULTS_FILENAME | grep -q 'failed'
if [ $? != 0 ] ; then
  printf "${LIGHT_GREEN}No failed tests.\n${NC}"
else
  printf "${LIGHT_RED}Failed tests:\n${NC}"
  printf "${RED}"
  jq -M ' .examples[] | select(.status == "failed")| "\(.full_description)\n"' $RESULTS_FILENAME | awk '{print $0}' | sed 's/Command.*$//'
  printf "${NC}"
fi
printf "${LIGHT_GREEN}Summary: ${BROWN}"
jq -M '.summary_line' $RESULTS_FILENAME
printf "${NC}"
>/dev/null popd
>/dev/null popd