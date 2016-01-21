package raknet

import (
	"testing"
	"time"
)

func TestCreateRouter(t *testing.T) {
	_, err := CreateRouter(19132)
	time.Sleep(time.Millisecond * 250)
	if err != nil {
		t.Error("Test failed: error occured while creating router:", err.Error())
	}
}
