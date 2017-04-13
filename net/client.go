package network

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func Dial(ip string, port string) string {
	server := ip + ":" + port
	conn, err := net.Dial("tcp4", server); if err != nil {
		println("error: ", err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, _ := bufio.NewReader(conn).ReadString('\n')
	println("status: ", status)

	return server
}

func Connect(ip string, port string, dir string, packageName string, op bool) {
	server := Dial(ip, port)

	if op {
		err := os.MkdirAll(packageName, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		client := &http.Client{}
		pkg := server + "/" + packageName + ".zip"
		response, err := client.Get(pkg)
		if err != nil {
			fmt.Println("Error downloading package:", err)
			return
		}

		defer response.Body.Close()

		filePath := filepath.Join(dir, packageName)
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating package file:", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println("Error saving package file:", err)
			return
		}

		// unzip and remove zip file

		fmt.Println("Package pulled successfully")
	} else {
		dir, err := os.Open(packageName)
		if err != nil {
			fmt.Println("Error opening package directory:", err)
			return
		}

		defer dir.Close()

		fileInfos, err := dir.Readdir(-1)
		if err != nil {
			fmt.Println("Error reading package directory:", err)
			return
		}

		client := &http.Client{}
		for _, fileInfo := range fileInfos {
			filePath := filepath.Join(packageName, fileInfo.Name())
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("Error opening package file:", err)
				return
			}
			defer file.Close()

			request, err := http.NewRequest("POST", server + "/" + packageName, file)
			if err != nil {
				fmt.Println("Error creating upload request:", err)
				return
			}

			request.Header.Set("File-Name", fileInfo.Name())
			response, err := client.Do(request)
			if err != nil {
				fmt.Println("Error uploading package file:", err)
				return
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				fmt.Println("Error uploading package file: Unexpected response status", response.StatusCode)
				return
			}

			fmt.Println("Package pushed successfully")
		}
	}
}
