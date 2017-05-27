// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func init() {
	// silence any logging done in the package files
	log.SetOutput(ioutil.Discard)
}

func TestWinPathList2Nix(t *testing.T) {
	pth := []string{canary[0], `C:\Apps\rubies\ruby-2.1.0\bin`, canary[1], `C:\some\fake\path`}
	result := strings.Join(winPathToNix(&pth), `:`)

	if strings.Contains(result, `C:`) {
		t.Errorf("Generated *nix path contains Windows volume designator")
	}
	if strings.ContainsRune(result, '\\') {
		t.Errorf(`Generated *nix path contains '\' char`)
	}
	if !strings.Contains(result, canary[0]) {
		t.Errorf("Generated *nix path missing `%s` start canary", canary[0])
	}
	if !strings.Contains(result, canary[1]) {
		t.Errorf("Generated *nix path missing `%s` stop canary", canary[1])
	}
}
