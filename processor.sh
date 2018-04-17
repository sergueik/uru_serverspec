#!/bin/bash

RESULTS_FILENAME='result_.json'
RESULTS_DIRECTORY='results'

# https://stedolan.github.io/jq/manual/

which jq 1 > /dev/null 2 >& 1
if [ $? == 0 ] ; then
  pusd $RESULTS_DIRECTORY
  jq ' .examples[] | select(.status != "passed")| .description,.status' $RESULTS_FILENAME
  jq '.summary_line' $RESULTS_FILENAME
  popd
fi
