package tests

import (
	"DCRS/dcrs"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func TestAdd(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = ioutil.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}

	// Invoke the Add function
	dcrs.Add("test.txt", tempDir)

	// Assert that the file was added
	trackerFile := filepath.Join(tempDir, ".obj", "tracker.txt")
	trackerData, err := ioutil.ReadFile(trackerFile)
	if err != nil {
		t.Fatalf("Failed to read tracker file: %s", err)
	}
	trackerContent := string(trackerData)
	expectedLine := fmt.Sprintf("%s uncommitted", filepath.Join(tempDir, ".obj", "test.txt"))
	if !strings.Contains(trackerContent, expectedLine) {
		t.Errorf("Expected line not found in tracker file: %s", expectedLine)
	}
}