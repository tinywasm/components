//go:build !wasm

package selectsearch

import "testing"

func TestSheetValidates(t *testing.T) {
	c := &SelectSearch{}
	c.Init(nil)
	if errs := c.sheet().Validate(); len(errs) > 0 {
		t.Errorf("selectsearch sheet must validate, got:\n%v", errs)
	}
}
