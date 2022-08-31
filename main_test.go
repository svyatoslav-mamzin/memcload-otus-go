package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stkrizh/otus-go-memcload/appsinstalled"
)

type protoTestcase struct {
	lat  float64
	lon  float64
	apps []uint32
}

func areSlicesEqual(a, b []uint32) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestProto(t *testing.T) {

	var protoCases = []protoTestcase{
		{42.345, 33.5677, []uint32{1, 2, 3}},
		{-42.345, 0.0, []uint32{100, 200, 300}},
		{-0.0, 100, nil},
		{0.0, 0.0, nil},
	}

	for _, tcase := range protoCases {

		test := &appsinstalled.UserApps{
			Lon:  &tcase.lon,
			Lat:  &tcase.lat,
			Apps: tcase.apps,
		}

		data, err := proto.Marshal(test)
		if err != nil {
			t.Error("marshaling error: ", err)
		}

		newTest := &appsinstalled.UserApps{}
		err = proto.Unmarshal(data, newTest)
		if err != nil {
			t.Error("unmarshaling error: ", err)
		}

		if *test.Lon != *newTest.Lon {
			t.Error("Lon-s are not equal: ", test.Lon, newTest.Lon)
		}

		if *test.Lat != *newTest.Lat {
			t.Error("Lat-s are not equal: ", test.Lat, newTest.Lat)
		}

		if !areSlicesEqual(test.Apps, newTest.Apps) {
			t.Error("Apps-s are not equal: ", test.Apps, newTest.Apps)
		}

	}

}

func TestParseRecordErrors(t *testing.T) {
	validCases := []string{
		"gaid\t123456\t24.567\t42.1344\t1,2,3,4,5",
		"aaa\t123456\t2\t4\t1",
		"bbb\ta      \t-9999\t0\t",
	}

	for _, tcase := range validCases {
		_, err := ParseRecord(tcase)
		if err != nil {
			t.Error(err)
		}
	}

	invalidCases := []string{
		"",
		"gaid\t123456\t24.56742.1344\t1,2,3,4,5",
		"aaa\t123456\tasd\t4\t1",
		"bbb\ta      \t-9999\ta\t",
		"aaa\t-9999\ta\t",
	}
	for _, tcase := range invalidCases {
		_, err := ParseRecord(tcase)
		if err == nil {
			t.Error("Error was not returned for:", tcase)
		}
	}
}

func TestParseRecord(t *testing.T) {
	cases := []string{
		"gaid\t123456\t24.567\t42.1344\t1,2,3,4,5",
		"aaa\t123456\t2\t4\t1",
		" \t111\t0\t0\t   ",
	}
	expected := []Record{
		{"gaid", "123456", 24.567, 42.1344, []uint32{1, 2, 3, 4, 5}},
		{"aaa", "123456", 2.0, 4.0, []uint32{1}},
		{" ", "111", 0.0, 0.0, []uint32{}},
	}

	for ix, tcase := range cases {
		record, _ := ParseRecord(tcase)

		if record.Type != expected[ix].Type {
			t.Error("Types are not equal:", tcase, expected[ix].Type)
		}

		if record.ID != expected[ix].ID {
			t.Error("IDs are not equal:", tcase, expected[ix].ID)
		}

		if record.Lat != expected[ix].Lat {
			t.Error("Lats are not equal:", tcase, expected[ix].Lat)
		}

		if record.Lon != expected[ix].Lon {
			t.Error("Lons are not equal:", tcase, expected[ix].Lon)
		}

		if len(record.Apps) != len(expected[ix].Apps) {
			t.Error("Len of Apps are not equal:", record.Apps, expected[ix].Apps)
			return
		}

		for appIx, app := range record.Apps {
			if app != expected[ix].Apps[appIx] {
				t.Error("Apps are not equal:", tcase, expected[ix].Apps)
			}
		}
	}

}
