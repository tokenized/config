package config

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type DurationTestStruct struct {
	Duration1  Duration `json:"duration_1"`
	OtherField string   `json:"other_field"`
	Duration2  Duration `json:"duration_2"`
}

type DurationTestStringStruct struct {
	Duration1  string `json:"duration_1"`
	OtherField string `json:"other_field"`
	Duration2  string `json:"duration_2"`
}

func Test_Duration(t *testing.T) {
	var tests = []struct {
		Struct          DurationTestStruct
		duration1       time.Duration
		duration2       time.Duration
		duration1String string
		duration2String string
	}{
		{
			Struct: DurationTestStruct{
				Duration1:  NewDuration(time.Second * 5),
				OtherField: "test1",
				Duration2:  NewDuration(time.Minute * 30),
			},
			duration1:       time.Second * 5,
			duration2:       time.Minute * 30,
			duration1String: "5s",
			duration2String: "30m0s",
		},
		{
			Struct: DurationTestStruct{
				Duration1:  NewDuration(time.Second),
				OtherField: "test2",
				Duration2:  NewDuration(time.Hour * 4),
			},
			duration1:       time.Second,
			duration2:       time.Hour * 4,
			duration1String: "1s",
			duration2String: "4h0m0s",
		},
	}

	for i, test := range tests {
		js, _ := json.MarshalIndent(test.Struct, "", "  ")
		t.Logf("JSON (%d) : %s", i, js)

		jsTest := &DurationTestStruct{}
		if err := json.Unmarshal(js, jsTest); err != nil {
			t.Fatalf("Failed to unmarshal json : %s", err)
		}

		if jsTest.Duration1 != test.Struct.Duration1 {
			t.Errorf("Wrong duration 1 (%d) : got %s, want %s", i, jsTest.Duration1,
				test.Struct.Duration1)
		}

		if jsTest.Duration1.Duration != test.duration1 {
			t.Errorf("Wrong duration 1 (%d) : got %s, want %s", i, jsTest.Duration1,
				test.duration1)
		}

		if jsTest.OtherField != test.Struct.OtherField {
			t.Errorf("Wrong other field (%d) : got %s, want %s", i, jsTest.OtherField,
				test.Struct.OtherField)
		}

		if jsTest.Duration2 != test.Struct.Duration2 {
			t.Errorf("Wrong duration 2 (%d) : got %s, want %s", i, jsTest.Duration2,
				test.Struct.Duration2)
		}

		if jsTest.Duration2.Duration != test.duration2 {
			t.Errorf("Wrong duration 2 (%d) : got %s, want %s", i, jsTest.Duration2,
				test.duration2)
		}

		stringTest := &DurationTestStringStruct{}
		if err := json.Unmarshal(js, stringTest); err != nil {
			t.Fatalf("Failed to unmarshal json to string test : %s", err)
		}

		if stringTest.Duration1 != test.duration1String {
			t.Errorf("Wrong duration 1 string (%d) : got %s, want %s", i, stringTest.Duration1,
				test.duration1String)
		}

		if stringTest.OtherField != test.Struct.OtherField {
			t.Errorf("Wrong other field (%d) : got %s, want %s", i, stringTest.OtherField,
				test.Struct.OtherField)
		}

		if stringTest.Duration2 != test.duration2String {
			t.Errorf("Wrong duration 2 string (%d) : got %s, want %s", i, stringTest.Duration2,
				test.duration2String)
		}

		// Test conversion to string
		duration1String := fmt.Sprintf("%s", test.Struct.Duration1)
		if duration1String != test.duration1String {
			t.Errorf("Wrong duration 1 string (%d) : got %s, want %s", i, duration1String,
				test.duration1String)
		}
		t.Logf("Duration 1 String (%d) : %s", i, duration1String)

		duration2String := fmt.Sprintf("%s", test.Struct.Duration2)
		if duration2String != test.duration2String {
			t.Errorf("Wrong duration 1 string (%d) : got %s, want %s", i, duration2String,
				test.duration2String)
		}
		t.Logf("Duration 2 String (%d) : %s", i, duration2String)

		// Test conversion from string
		var duration1 Duration
		if err := duration1.UnmarshalText([]byte(test.duration1String)); err != nil {
			t.Fatalf("Failed to unmarshal duration 1 to text : %s", err)
		}
		if duration1 != test.Struct.Duration1 {
			t.Errorf("Wrong duration 1 (%d) : got %s, want %s", i, duration1,
				test.Struct.Duration1)
		}

		var duration2 Duration
		if err := duration2.UnmarshalText([]byte(test.duration2String)); err != nil {
			t.Fatalf("Failed to unmarshal duration 2 to text : %s", err)
		}
		if duration2 != test.Struct.Duration2 {
			t.Errorf("Wrong duration 2 (%d) : got %s, want %s", i, duration2,
				test.Struct.Duration2)
		}
	}
}
