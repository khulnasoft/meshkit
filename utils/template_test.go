package utils

import (
	"bytes"
	"testing"
)

func TestSimplestTemplate(t *testing.T) {
	template := []byte("{{.namespace}}")
	result, err := MergeToTemplate(template, map[string]string{"namespace": "meshplay"})
	if err != nil {
		t.Errorf("err = %v; want 'nil'", err)
	}
	if !bytes.Equal(result, []byte("meshplay")) {
		t.Errorf("result = %s; want 'meshplay'", result)
	}
}

func TestEmptyTemplate(t *testing.T) {
	template := []byte("")
	result, err := MergeToTemplate(template, map[string]string{"namespace": "meshplay"})
	if err != nil {
		t.Errorf("err = %v; want 'nil'", err)
	}
	if !bytes.Equal(result, []byte("")) {
		t.Errorf("result = %s; want '' (empty string)", result)
	}
}

func TestEmptyDataMap(t *testing.T) {
	template := []byte("{{.namespace}}")
	result, err := MergeToTemplate(template, map[string]string{})
	if err != nil {
		t.Errorf("err = %v; want 'nil'", err)
	}
	if !bytes.Equal(result, []byte("")) {
		t.Errorf("result = %s; want '' (empty string)", result)
	}
}

func TestMultilineTemplate(t *testing.T) {
	template := `KhulnaSoft
{{.project}} is
best`
	expected := `KhulnaSoft
Meshplay is
best`
	result, err := MergeToTemplate([]byte(template), map[string]string{"project": "Meshplay"})
	if err != nil {
		t.Errorf("err = %v; want 'nil'", err)
	}
	if !bytes.Equal(result, []byte(expected)) {
		t.Errorf("result = %s; want '%s'", result, expected)
	}
}
