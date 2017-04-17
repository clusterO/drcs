package tests

import (
	"DCRS/dcrs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "init_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Call the Init function
	dcrs.Init(tempDir, "y")

	// Verify that the necessary files and directories were created
	_, err = os.Stat(filepath.Join(tempDir, ".obj"))
	if err != nil {
		t.Errorf("Expected .obj directory to be created, but got error: %v", err)
	}

	_, err = os.Stat(filepath.Join(tempDir, ".obj", "config.txt"))
	if err != nil {
		t.Errorf("Expected config.txt file to be created, but got error: %v", err)
	}

	_, err = os.Stat(filepath.Join(tempDir, ".obj", "tracker.txt"))
	if err != nil {
		t.Errorf("Expected tracker.txt file to be created, but got error: %v", err)
	}
}
