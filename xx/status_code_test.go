// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx_test

import (
	"github.com/orivil/morgine/xx"
	"reflect"
	"testing"
)

func TestStatusCodes_StatusTexts(t *testing.T) {
	scs := xx.NewStatusCodes()
	scs.InitStatus(xx.DefaultStatusNamespace, 2000, "success")
	scs.InitStatus("namespace1", 2000, "success")
	got := scs.StatusTexts()
	need := xx.StatusTexts{
		2000: []string{"success", "namespace1-success"},
	}
	if !reflect.DeepEqual(got, need) {
		t.Errorf("need: %v got: %v", need, got)
	}
}
