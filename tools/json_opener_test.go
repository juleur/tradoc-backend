package tools

import "testing"

func TestJSONOpener(t *testing.T) {
	json := OpenDialectsJSONFile()
	t.Log(json)
}
