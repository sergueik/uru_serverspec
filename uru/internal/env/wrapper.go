// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

var BashWrapper = `uru()
{
  export URU_INVOKER='bash'

  # uru_rt must already be on PATH
  uru_rt "$@"

  if [[ -d "$URU_HOME" ]]; then
    if [[ -f "$URU_HOME/uru_lackee" ]]; then
      . "$URU_HOME/uru_lackee"
    fi
  else
    if [[ -f "$HOME/.uru/uru_lackee" ]]; then
      . "$HOME/.uru/uru_lackee"
    fi
  fi
}
`

var FishWrapper = `function uru -d "Manage your ruby versions"
  set -x URU_INVOKER fish

  # uru_rt must already be on PATH
  uru_rt $argv

  if test -d "$URU_HOME" -a -f "$URU_HOME/uru_lackee.fish"
    source "$URU_HOME/uru_lackee.fish"
  else if test -f "$HOME/.uru/uru_lackee.fish"
    source "$HOME/.uru/uru_lackee.fish"
  end
end
`
