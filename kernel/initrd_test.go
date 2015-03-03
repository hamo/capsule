package kernel

import (
	"testing"
)

func TestBuildInitrd(t *testing.T) {
	if name, err := BuildInitrd("/lib/modules/3.14.33-1-lts"); err != nil {
		t.Errorf("%s", err)
	} else {
		t.Errorf("%s", name)

	}

}
