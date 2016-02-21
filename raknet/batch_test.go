package raknet

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

func TestFormatParse(t *testing.T) {
	var tests = []struct {
		fmt    string
		expect []util.FmtElement
	}{
		{"B", []util.FmtElement{
			{[]byte("B")[0], 1},
		}},
		{"S", []util.FmtElement{
			{[]byte("S")[0], 1},
		}},
		{"D32", []util.FmtElement{
			{[]byte("D")[0], 32},
		}},
		{"BBS", []util.FmtElement{
			{[]byte("B")[0], 1},
			{[]byte("B")[0], 1},
			{[]byte("S")[0], 1},
		}},
		{"B13", []util.FmtElement{
			{[]byte("B")[0], 13},
		}},
		{"S4D1B", []util.FmtElement{
			{[]byte("S")[0], 4},
			{[]byte("D")[0], 1},
			{[]byte("B")[0], 1},
		}},
		{"L7D90B32", []util.FmtElement{
			{[]byte("L")[0], 7},
			{[]byte("D")[0], 90},
			{[]byte("B")[0], 32},
		}},
		{"SSSS", []util.FmtElement{
			{[]byte("S")[0], 1},
			{[]byte("S")[0], 1},
			{[]byte("S")[0], 1},
			{[]byte("S")[0], 1},
		}},
		{"D1000", []util.FmtElement{
			{[]byte("D")[0], 1000},
		}},
		{"TTS", []util.FmtElement{
			{[]byte("T")[0], 1},
			{[]byte("T")[0], 1},
			{[]byte("S")[0], 1},
		}},
	}
	for _, test := range tests {
		result := util.ParseBatchFmt(test.fmt)
		if len(test.expect) != len(result) {
			t.Errorf("Parse result length mismatch: got %v, expected %v (%d != %d) (%s)", result, test.expect, len(result), len(test.expect), test.fmt)
			return
		}
		for i := range test.expect {
			if test.expect[i] != result[i] {
				t.Error("Batch format parse test failed: expected", test.expect, "got", util.ParseBatchFmt(test.fmt), "("+test.fmt+")")
				return
			}
		}
	}
}

func TestBatchRead(t *testing.T) {
	var tests = []struct {
		init   []byte
		fmt    string
		expect []interface{}
	}{
		{[]byte{7, 4, 1, 2}, "BBBB", []interface{}{7, 4, 1, 2}},
		{[]byte{0, 3, 6, 3}, "BBS", []interface{}{0, 3, 1539}},
		{[]byte{0, 0, 0, 0}, "SS", []interface{}{0, 0}},
		{[]byte{0, 0, 0, 0}, "I", []interface{}{0}},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 1}, "L", []interface{}{1}},
		{[]byte{0, 6, 0, 2, 0, 0, 3, 0}, "SSI", []interface{}{6, 2, 768}},
		{[]byte{0, 32, 1, 4}, "SS", []interface{}{32, 260}},
	}
	for _, test := range tests {
		b := bytes.NewBuffer(test.init)
		r, err := b.BatchRead(test.fmt)
		if err != nil {
			t.Error(
				"Got error while BatchRead:",
				err,
				"\nInit:\n"+hex.Dump(test.init)+
					"Format:", test.fmt+
					"\nExpected:", test.expect,
			)
			return
		}
		if len(r) != len(test.expect) {
			t.Error(
				"BatchRead return value length mismatch: Got",
				strconv.Itoa(len(r))+",",
				"expected", len(test.expect),
				"\nInit:\n"+hex.Dump(test.init)+
					"Format:", test.fmt+
					"\nExpected:", test.expect,
			)
			return
		}
		if fmt.Sprint(r) != fmt.Sprint(test.expect) {
			t.Error(
				"BatchRead return value mismatch: Got",
				r,
				"expected", test.expect,
				"\nInit:\n"+hex.Dump(test.init)+
					"Format:", test.fmt+
					"\nExpected:", test.expect,
			)
		}
	}
}
