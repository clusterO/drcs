package tests

import (
	network "DCRS/net"
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDial(t *testing.T) {
	expected := "localhost:8181"
	result := network.Dial("localhost", "8181")

	if result != expected {
		t.Errorf("Expected server address %s, but got %s", expected, result)
	}
}

func TestConnect(t *testing.T) {
	// Set up a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method == http.MethodGet && r.URL.Path == "/testpkg.zip" {
			// Create a test package file
			testPkg := []byte("Test package content")

			// Write the test package file to the response
			w.Write(testPkg)
		} else if r.Method == http.MethodPost && r.URL.Path == "/testpkg" {
			// Check the request headers
			fileName := r.Header.Get("File-Name")
			if fileName != "file1.txt" && fileName != "file2.txt" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Read the uploaded file contents
			fileContents, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Do something with the file contents (e.g., save it to disk)
			_ = ioutil.WriteFile(fileName, fileContents, 0644)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Set up a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "connect-test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Test pulling a package
	network.Connect(server.URL, "8181", tempDir, "testpkg", true)

	// Check if the package file was downloaded
	packagePath := filepath.Join(tempDir, "testpkg")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		t.Error("Failed to download the package file:", err)
	}

	// Test pushing a package
	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")
	if err := ioutil.WriteFile(file1Path, []byte("File 1 contents"), 0644); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	if err := ioutil.WriteFile(file2Path, []byte("File 2 contents"), 0644); err != nil {
		t.Fatal("Failed to create test file:", err)
	}

	network.Connect(server.URL, "8181", tempDir, "testpkg", false)

	// Check if the files were uploaded
	uploadedFile1Path := filepath.Join(tempDir, "testpkg", "file1.txt")
	if _, err := os.Stat(uploadedFile1Path); os.IsNotExist(err) {
		t.Error("Failed to upload file1.txt:", err)
	}
	uploadedFile2Path := filepath.Join(tempDir, "testpkg", "file2.txt")
	if _, err := os.Stat(uploadedFile2Path); os.IsNotExist(err) {
		t.Error("Failed to upload file2.txt:", err)
	}
}

func TestListen(t *testing.T) {
	// Start a test server
	go func() {
		err := network.Listen()
		if err != nil {
			t.Errorf("Listen error: %s", err)
		}
	}()

	// Connect to the test server
	conn, err := net.Dial("tcp4", "127.0.0.1:8181")
	if err != nil {
		t.Errorf("Failed to connect to the test server: %s", err)
	}

	// Send a sample message
	message := "Hello, server!"
	_, err = fmt.Fprintf(conn, message)
	if err != nil {
		t.Errorf("Failed to send message: %s", err)
	}

	// Receive the response from the server
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("Failed to read response: %s", err)
	}

	// Check if the response is as expected
	expectedResponse := "Received: Hello, server!\n"
	if response != expectedResponse {
		t.Errorf("Unexpected response from server. Expected: %s, Got: %s", expectedResponse, response)
	}

	// Close the connection
	conn.Close()
}
