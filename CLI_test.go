package oregontrail_test

import (
	"bytes"
	"strings"
	"testing"

	trail "github.com/jwc20/oregontrail"
)

func TestCLIInitSVT(t *testing.T) {
	out := &bytes.Buffer{}
	cli := trail.NewCLI(strings.NewReader(""), out)
	cli.InitSVT()

	if !cli.State.Trip.FortAvailable {
		t.Error("FortAvailable should be true")
	}

}
