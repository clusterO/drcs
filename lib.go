package dcrs

import (
	"os"
	"log"
	"bytes"
	"compress/zlib"
	"path/filepath"
	"fmt"
	"io"
	"bufio"
	"strings"
)

// GetFile return file path
func GetFile(commitTag string, filename string) string {
    path, err := os.Getwd()
    if err != nil {
        log.Println(err)
    }

    files := filepath.Join(path, "/dcrs/object/", commitTag)
    hashmap := filepath.Join(files, "hashmap")
    h := GetHashNameFromHashmap(hashmap, filename)
    
    var b bytes.Buffer
    content := zlib.NewWriter(&b)
    p := filepath.Join(files, h)
    fmt.Println(p)
    r, err := zlib.NewReader(&b)
    io.Copy(content, r)
    r.Close() 

    return ""
}

// GetHashNameFromHashmap returns name
func GetHashNameFromHashmap(hashfile string, name string) string {
    files, err := os.Open(hashfile)
    if err != nil {
        log.Fatal(err)
    }

    defer files.Close()
    scanner := bufio.NewScanner(files)

    for scanner.Scan() { 
        line := scanner.Text()
        p := strings.Fields(line)
        
        if p[0] == name {
            return p[1]
        }
    }

    return ""
}

// Fileattr file attribute
type Fileattr struct {
	fileinfo os.FileInfo
	path string
}

var leftList []Fileattr
var rightList []Fileattr

// Dircmp compare two direcotories
func Dircmp(leftDir string, rightDir string) {
	filepath.Walk(leftDir, LeftVisit)
	filepath.Walk(rightDir, RightVisit)
	
	cmp(leftDir, rightDir, leftList, rightList)
}

func cmp(dir1 string, dir2 string, list1 []Fileattr, list2 []Fileattr) {
	if (len(list1) > len(list2)) {
			fmt.Printf("Directory: %s supersedes in total files %d over Directory: %s\n", dir1, (len(list1) - len(list2)), dir2)
	} else {
			fmt.Printf("Directory %s supersedes in total files %d over Directory: %s\n", dir2, (len(list2) - len(list1)), dir1)
	}
}

// LeftVisit left directory
func LeftVisit(path string, f os.FileInfo, _ error) error {
	var attr Fileattr

	attr.path = path
	attr.fileinfo = f

	leftList = Append(leftList, attr)
	return nil
}

// RightVisit right directory
func RightVisit(path string, f os.FileInfo, _ error) error {
	var attr Fileattr

	attr.path = path
	attr.fileinfo = f

	rightList = Append(rightList, attr)
	return nil
}

// Append to
func Append(slice []Fileattr, data Fileattr) []Fileattr {
	l := len(slice)
	if data.path != "" {
			newSlice := make([]Fileattr, l+1)
			copy(newSlice, slice)
			slice = newSlice
			slice[l] = data

	}
	return slice
}