package hip

import (
	"github.com/pkg/errors"
	"os"
	"testing"
)

func TestIntegrationally(t *testing.T) {

	hc, err := NewHipchatRoomPrinter(os.Getenv("HIPCHAT_TOKEN"), os.Getenv("HIPCHAT_ROOM"))
	if err != nil {
		t.Fatal(err)
	}

	hc.Out("out")
	hc.Err("err")
	hc.Done(errors.New("blah"))
}
