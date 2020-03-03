#!/bin/bash

# set -x
# URU_HOME='<%= @uru_home -%>'
if [ -z $URU_HOME ]
then
  URU_HOME='/uru'
  URU_HOME='/home/uru_user/uru_homebased'
fi
export URU_HOME

export URU_INVOKER=bash
export LD_LIBRARY_PATH=$URU_HOME/ruby/lib

GEM_VERSION='2.1.0'
RAKE_VERSION='10.1.0'
RUBY_VERSION='2.1.0'
RUBY_VERSION_LONG='2.1.9p490'
RUBY_TAG_LABEL='219p490'

GEM_VERSION='2.3.0'
RAKE_VERSION='10.4.2'
RUBY_VERSION='2.3.6'
RUBY_VERSION_LONG='2.3.6p384'
RUBY_TAG_LABEL='236p384'

URU_RUNNER=$URU_HOME/uru_rt

pushd $URU_HOME

# TODO: execute
# $URU_RUNNER admin refresh
# when the ~/.uru/rubies.json, in particular the GemHome, is different

if [ -z $HOME ] ; then
  if [[ "$EUID" -ne 0 ]]
  then
    HOME="/home/${USER}"
  else
    HOME='/root'
  fi

fi

if [[ ! -d $HOME/.uru ]]; then mkdir "$HOME/.uru"; fi
RUBIES="$HOME/.uru/rubies.json"
rm -f $RUBIES
# for Ruby 2.1.0
cat <<EOF>$RUBIES
{
  "Version": "1.0.0",
  "Rubies": {
  "2357568376": {
    "ID": "$RUBY_VERSION_LONG",
    "TagLabel": "$RUBY_TAG_LABEL",
    "Exe": "ruby",
    "Home": "$URU_HOME/ruby/bin",
    "GemHome": "$URU_HOME/.gem/ruby/$GEM_VERSION",
    "Description": "ruby $RUBY_VERSION_LONG (2016-03-30 revision 54437) [x86_64-linux]"
    }
 }
}
EOF
# for Ruby 2.3.6
cat <<EOF>$HOME/.uru/rubies.json
{
  "Version": "1.0.0",
  "Rubies": {
  "2357568376": {
    "ID": "$RUBY_VERSION_LONG",
    "TagLabel": "$RUBY_TAG_LABEL",
    "Exe": "ruby",
    "Home": "$URU_HOME/ruby/bin",
    "GemHome": "$URU_HOME/.gem/ruby/$GEM_VERSION",
    "Description": "ruby $RUBY_VERSION_LONG (2017-12-14 revision 61254) [x86_64-linux]"
    }
 }
}
EOF

echo Y |$URU_RUNNER admin rm $RUBY_TAG_LABEL > /dev/null
$URU_RUNNER admin add ruby/bin

if [ ! -z $DEBUG ]
then
  $URU_RUNNER ls --verbose
fi
if [ ! -z $DEBUG ]
then
  TAG=`$URU_RUNNER ls 2>& 1|awk '{print $1}'`
  $URU_RUNNER $TAG
fi

# Copy .gems to default location

# TODO: fix it properly
cp -R .gem $HOME

if [ ! -z $DEBUG ]
then
  # Verify the gems
  $URU_RUNNER gem list --local
fi

# Check that the required gems are present
$URU_RUNNER gem list| grep -qi serverspec
if [ $? != 0 ]; then
  echo 'ERROR: serverspec gem is not found'
  exit 1
fi

# Run the serverspec
$URU_RUNNER ruby ruby/lib/ruby/gems/$GEM_VERSION/gems/rake-$RAKE_VERSION/bin/rake spec
