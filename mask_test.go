package config

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type TestJSONStruct struct {
	Key        string            `json:"key" masked:"true"`
	Value      string            `json:"value"`
	Value2     string            `json:"value2"`
	Marshaler  testJSONMarshaler `json:"marshaler"`
	Marshaler2 testJSONMarshaler `json:"marshaler2" masked:"true"`
	Struct     testJSONStruct    `json:"struct"`
	Struct2    testJSONStruct    `json:"struct2" masked:"true"`
	LastField  string            `json:"last"`
}

func Test_MarshalJSONMasked(t *testing.T) {
	value := TestJSONStruct{
		Key:    "Don't show this",
		Value:  "Show this",
		Value2: "Show this too",
		Marshaler: testJSONMarshaler{
			field1: 1,
			field2: 2,
		},
		Marshaler2: testJSONMarshaler{
			field1: 3,
			field2: 4,
		},
		Struct: testJSONStruct{
			Field1: 5,
			Field2: -1,
			field3: 7,
		},
		Struct2: testJSONStruct{
			Field1: 8,
			Field2: -2,
			field3: 10,
		},
		LastField: "final value",
	}

	b, err := MarshalJSONMasked(value)
	if err != nil {
		t.Fatalf("Failed to marshal value : %s", err)
	}

	t.Logf("JSON : %s\n", string(b))

	jsRaw, err := MarshalJSONMaskedRaw(value)
	if err != nil {
		t.Fatalf("Failed to marshal value raw : %s", err)
	}

	bRaw, err := jsRaw.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal raw json : %s", err)
	}
	t.Logf("Raw JSON : %s", string(bRaw))

	if !bytes.Equal(b, bRaw) {
		t.Errorf("Bytes not equal to Raw Bytes")
	}

	s := string(b)

	if strings.Contains(s, "Don't") {
		t.Errorf("Should not contain \"Don't\"")
	}

	if !strings.Contains(s, "\"marshaler2\":\"***\"") {
		t.Errorf("Should contain \"marshaler2\":\"***\"")
	}

	if !strings.Contains(s, "\"struct2\":\"***\"") {
		t.Errorf("Should contain \"struct2\":\"***\"")
	}

	if strings.Contains(s, "-1") {
		t.Errorf("Should not contain \"-1\"")
	}

	if strings.Contains(s, "-2") {
		t.Errorf("Should not contain \"-2\"")
	}
}

type testJSONMarshaler struct {
	field1 int `json:"field_1"`
	field2 int `json:"field_2"`
}

func (t testJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t.field1+t.field2)), nil
}

type testJSONStruct struct {
	Field1 int `json:"field_1"`
	Field2 int `json:"field_2_m" masked:"true"`
	field3 int `json:"field_3"`
}
