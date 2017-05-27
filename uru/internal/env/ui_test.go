// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"testing"
)

func TestYesRegex(t *testing.T) {
	if yResp.MatchString(`N`) || yResp.MatchString(`n`) {
		t.Error("incorrectly matched a `no` type response")
	}
	if !(yResp.MatchString(`yes`) && yResp.MatchString(`YesPlease`)) {
		t.Error("did not match a `yes` type response")
	}
}
