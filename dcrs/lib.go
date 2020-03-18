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
	"io/ioutil"
	"dcrs/zip"
	"go-three-way-merge"
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

// GetAllCommits return all commits
func GetAllCommits() []string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	files, err := os.Open(path + "/dcrs/" + "status.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)

	var content []string 
	for scanner.Scan() {
		line := scanner.Text()
		_ = append(content, line)
	}
	
	return content
}

// CompressAndSend returns content
func CompressAndSend(commit string) string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	tempdir, err := ioutil.TempDir("", "tempdir")
	archivename := filepath.Join(tempdir, commit + ".zip")
	zip.RecursiveZip(filepath.Join(path, commit), archivename)

	files, err := os.Open(archivename)
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)

	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
	}

	return content
}

// UncompressAndWrite extract content
func UncompressAndWrite(commit string, content string) {
	tempdir, err := ioutil.TempDir("", "tempdir")
	archivename := filepath.Join(tempdir, commit + ".zip")

	file, err := os.Stat(archivename)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	switch mode := file.Mode(); {
		case mode.IsRegular():
			archive, err := os.Create(archivename)
			if err != nil {
				panic(err)
			}
				
			defer archive.Close()

			_, err = archive.WriteString(content) 
				if err != nil {
					panic(err)
				}
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	extractto := filepath.Join(path, commit)
	os.MkdirAll(extractto, os.ModePerm)
	zip.Unzip(extractto, archivename)
}

// GetCommits return commits
func GetCommits(d string) []string {
    return GetAllCommits()
}

// GetCommitsContent return content from commit file
func GetCommitsContent(d string, c string) string {
    return CompressAndSend(c)
}

// Difference Set A - B
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return
}

// GetFileList return list of files
func GetFileList(commitTag string) []string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	files, err := os.Open(filepath.Join(path, commitTag, "hashmap.txt"))
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)

	var list []string
	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)
		_ = append(list, p[0], p[1])
	}

	return list
}

// Data merge output structure
type Data struct {
    mdContent string
	conflict int
	merged int
}


// MergeMethod use 3 base method to merge files
func MergeMethod(base string, mine string, other string) Data {
    m := Merge(base, mine, other)
    mg := m.merge_groups()
    conflicts := 0
	flag := 0
	
    for g := range mg {
        if g[0] == "conflict" {
            conflicts++
		}

        if g[0] == "a" {
            flag++
		}
	}

    merged := strings.Join(m.merge_lines(`start_marker = "\n!!!--Conflict--!!!\n!--Your version--", mid_marker = "\n!--Other version--", end_marker = "\n!--End conflict--\n`), "")
	r := Data{mdContent: merged, conflict: conflicts, merged: flag}

	return r
}