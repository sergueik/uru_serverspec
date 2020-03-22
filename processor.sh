#!/bin/bash

# RESULTS_DIRECTORY='reports'
# DEFAULT_RESULTS_BASENAME='report_'

RESULTS_DIRECTORY='results'
DEFAULT_RESULTS_BASENAME='result_'
RESULTS_FILENAME=${1:-$DEFAULT_RESULTS_BASENAME.json}
RESULTS_BASENAME=${RESULTS_FILENAME%.*}

WORKDIR=${2:-/uru}

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
  echo 'jq not found'
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
RESULTS_FILTERED_FILENAME="${RESULTS_BASENAME}_filtered.json"
# assuming there are two distinct kinds of spec files and we are only interested in the errors of one
# the below jq command filter that, but currently loses the summary details
# comment if throwing an error
SPEC_FILE='role_spec.rb'
jq '[.examples[]|select(.file_path|endswith("role_spec.rb"))]' $RESULTS_FILENAME > $RESULTS_FILTERED_FILENAME
# NOTE: the jq command seems to require single quotes
# jq: error: syntax error, unexpected INVALID_CHARACTER
EXPR='[.examples[]|select(.file_path|endswith("'$SPEC_FILE'"))]'
jq $EXPR $RESULTS_FILENAME > $RESULTS_FILTERED_FILENAME
printf "${NC}"
>/dev/null popd
>/dev/null popd
