//go:build !wasm

package listselect_test

import (
	"testing"

	"github.com/tinywasm/components/listselect"
)

func TestZeroValueIsNormalMode(t *testing.T) {
	var m listselect.Mode
	if m.On().Get() {
		t.Error("Mode{} zero value must not be in selection mode")
	}
	if len(m.CheckedIDs([]string{"a", "b"})) != 0 {
		t.Error("Mode{} zero value must have no checked items")
	}
}

func TestToggleMarksAndUnmarks(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)

	if m.IsChecked("1") {
		t.Error("item 1 should not be checked initially")
	}

	m.Toggle("1")
	if !m.IsChecked("1") {
		t.Error("item 1 should be checked after toggle")
	}

	m.Toggle("1")
	if m.IsChecked("1") {
		t.Error("item 1 should be unchecked after second toggle")
	}
}

func TestLeavingModeClearsTheMarks(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)
	m.Toggle("1")
	m.Toggle("2")

	if !m.IsChecked("1") || !m.IsChecked("2") {
		t.Fatal("items should be checked")
	}

	m.SetOn(false)
	m.SetOn(true)

	if m.IsChecked("1") || m.IsChecked("2") {
		t.Error("leaving selection mode must clear marks")
	}
}

func TestCheckedIDsFollowTheGivenOrder(t *testing.T) {
	var m listselect.Mode
	m.SetOn(true)

	// Tap in reverse order: 3rd item, then 1st item
	m.Toggle("3")
	m.Toggle("1")

	renderOrder := []string{"1", "2", "3"}
	checked := m.CheckedIDs(renderOrder)

	if len(checked) != 2 {
		t.Fatalf("expected 2 checked items, got %d", len(checked))
	}
	if checked[0] != "1" || checked[1] != "3" {
		t.Errorf("CheckedIDs must follow render order [1, 3], got %v", checked)
	}
}

func TestOnChangeReportsTheCount(t *testing.T) {
	var m listselect.Mode
	var lastCount int
	m.OnChange = func(n int) {
		lastCount = n
	}

	m.SetOn(true)
	m.Toggle("1")
	if lastCount != 1 {
		t.Errorf("expected count 1, got %d", lastCount)
	}

	m.Toggle("2")
	if lastCount != 2 {
		t.Errorf("expected count 2, got %d", lastCount)
	}

	m.Toggle("1")
	if lastCount != 1 {
		t.Errorf("expected count 1, got %d", lastCount)
	}

	m.SetOn(false)
	if lastCount != 0 {
		t.Errorf("expected count 0 on SetOn(false), got %d", lastCount)
	}
}
