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

func TestCommit(t *testing.T) {
	dir := "/path/to/directory" // Replace with the actual directory path

	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the temporary file
	content := "Test content"
	err = ioutil.WriteFile(tempFile.Name(), []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write content to temporary file: %v", err)
	}

	// Call the Commit function
	message := "Test commit"
	_, err = dcrs.Commit(message, dir)
	if err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}

	// Perform assertions to verify the expected changes after committing

	// Check if the temporary file exists in the committed directory
	committedFilePath := filepath.Join(dir, ".obj", "committed_directory", tempFile.Name())
	_, err = os.Stat(committedFilePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist in committed directory, but it does not", committedFilePath)
	}

	// Check if the message file exists in the committed directory
	messageFilePath := filepath.Join(dir, ".obj", "committed_directory", "message.txt")
	_, err = os.Stat(messageFilePath)
	if os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist in committed directory, but it does not", messageFilePath)
	}

	// Read the content of the message file and compare with the expected message
	messageContent, err := ioutil.ReadFile(messageFilePath)
	if err != nil {
		t.Fatalf("Failed to read message file: %v", err)
	}
	if string(messageContent) != message {
		t.Errorf("Expected message content '%s', but got '%s'", message, string(messageContent))
	}
}

func TestPull(t *testing.T) {
	dir := "/path/to/directory"
	url := "127.0.0.1:8181/packageName"

	// Perform the pull operation
	dcrs.Pull(url, dir)

	// Assert the expected results
	// ...
}

func TestPush(t *testing.T) {
	dir := "/path/to/directory"
	url := "127.0.0.1:8181/packageName"

	// Perform the push operation
	dcrs.Push(url, dir)

	// Assert the expected results
	// ...
}
