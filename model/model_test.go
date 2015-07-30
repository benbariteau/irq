package model

import (
	"testing"
)

func TestNewModel(t *testing.T) {
	model, err := NewModel("./test.db")
	if err != nil {
		t.Error("Got non-nil error: ", err)
	}
	if model.db == nil {
		t.Error("Got nil db connection")
	}
}
