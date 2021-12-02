package s3io

import "testing"

func TestGetPartSizeNoChange(t *testing.T) {
	newPartSize := GetPartSize(16*mb, 5*mb)
	expectedPartSize := 16 * mb
	if newPartSize != expectedPartSize {
		t.Errorf("no change: expected part size %s got %s", expectedPartSize, newPartSize)
	}
}

func TestGetPartSizeOneChange(t *testing.T) {
	newPartSize := GetPartSize(16*mb, 170*gb)
	expectedPartSize := 21 * mb
	if newPartSize != expectedPartSize {
		t.Errorf("one change: expected part size %s got %s", expectedPartSize, newPartSize)
	}
}

func TestGetPartSizeTwoChanges(t *testing.T) {
	newPartSize := GetPartSize(16*mb, 250*gb)
	expectedPartSize := 26 * mb
	if newPartSize != expectedPartSize {
		t.Errorf("two changes: expected part size %s got %s", expectedPartSize, newPartSize)
	}
}

func TestGetPartSizeZeroSize(t *testing.T) {
	newPartSize := GetPartSize(16*mb, 0)
	expectedPartSize := 16 * mb
	if newPartSize != expectedPartSize {
		t.Errorf("zero file: expected part size %s got %s", expectedPartSize, newPartSize)
	}

}
