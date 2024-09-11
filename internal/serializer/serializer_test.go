package serializer

import (
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestSerializeToString(t *testing.T) {
	testData := TestStruct{Name: "Test", Value: 123}
	expected := `{"name":"Test","value":123}`

	result, err := SerializeToString(testData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestDeserializeFromString(t *testing.T) {
	jsonData := `{"name":"Test","value":123}`
	var result TestStruct

	err := DeserializeFromString(jsonData, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := TestStruct{Name: "Test", Value: 123}
	if result != expected {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
