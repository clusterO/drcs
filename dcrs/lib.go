package dcrs

import (
	"DCRS/zip"
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	merge "go-three-way-merge"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Fileattr struct {
	fileinfo os.FileInfo
	path string
}

// Merge output structure
type Data struct {
    mdContent string
	conflict int
	merged int
}

var leftList []Fileattr
var rightList []Fileattr

// GetFile return hashed file path
func GetFile(dir string, commitTag string, filename string) string {
    files := filepath.Join(dir, ".obj", commitTag)
    hashmap := filepath.Join(files, "hashmap")
    h := GetHashNameFromHashmap(hashmap, filename)
    
    var b bytes.Buffer
    content := zlib.NewWriter(&b)
    p := filepath.Join(files, h)
    r, _ := zlib.NewReader(&b)
    io.Copy(content, r)
    r.Close() 

    return p
}

func GetHashNameFromHashmap(hashfile string, name string) string {
    files, err := os.Open(hashfile); if err != nil {
        log.Fatal(err)
    }; 
	
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

func LeftVisit(path string, f os.FileInfo, _ error) error {
	var attr Fileattr

	attr.path = path
	attr.fileinfo = f

	leftList = Append(leftList, attr)

	return nil
}

func RightVisit(path string, f os.FileInfo, _ error) error {
	var attr Fileattr

	attr.path = path
	attr.fileinfo = f

	rightList = Append(rightList, attr)

	return nil
}

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

func GetAllCommits(dir string) []string {
	files, err := os.Open(filepath.Join(dir, ".obj", "status.txt")); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	var content []string 
	for scanner.Scan() {
		line := scanner.Text()
		_ = append(content, line)
	}
	
	return content
}

func CompressAndSend(dir string, commit string) string {
	archivename := CompressAll(commit, dir)
	files, err := os.Open(archivename); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	content := ""
	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
	}

	return content
}

func CompressAll(commits string, commitDir string) string {
	tempdir, _ := os.MkdirTemp("", "tempdir")
	commitdir := filepath.Join(tempdir, commitDir)
	err := os.Mkdir(commitdir, os.ModePerm); if err != nil {
		log.Fatal(err)
	}

	archivename := filepath.Join(tempdir, commits + ".zip")
	file, err := os.Stat(archivename); if err != nil {
		fmt.Println(err)
		return ""
	}
	
	switch mode := file.Mode(); {
		case !mode.IsRegular():
			fp, err := os.Create(archivename); if err != nil {
				panic(err)
			}; 
			
			defer fp.Close()
			
			_, err = fp.WriteString(commits); if err != nil {
				panic(err)
			}
	}

	extractto := filepath.Join(commitDir, commits)
	os.MkdirAll(extractto, os.ModePerm)
	zip.Unzip(extractto, archivename)

	for range commits {
		filenames := GetFileName(commitDir, commits)
		// file := filepath.Join(commitDir, commits)
		// _, err = io.Copy(file, tempdir)

		for fn := range filenames {
			h := GetFileLocation(commitDir, commits, filenames[fn])
			_ = filepath.Join(commits, h)
			// _, err = io.Copy(floc, commitdir)
		}
	}

	_, _ = os.CreateTemp("", ".zip")
	zip.RecursiveZip(tempdir, archivename)

	return archivename
}

func UncompressAndWrite(dir string, commit string, content string) {
	tempdir, _ := os.MkdirTemp("", "tempdir")
	archivename := filepath.Join(tempdir, commit + ".zip")
	file, err := os.Stat(archivename); if err != nil {
		fmt.Println(err)
		return
	}
	
	switch mode := file.Mode(); {
		case mode.IsRegular():
			archive, err := os.Create(archivename); if err != nil {
				panic(err)
			}; defer archive.Close()
			_, err = archive.WriteString(content); if err != nil {
					panic(err)
				}
	}

	extractto := filepath.Join(dir, commit)
	os.MkdirAll(extractto, os.ModePerm)
	zip.Unzip(extractto, archivename)
}
// Why not called directly
func GetCommitsContent(dir string, c string) string {
    return CompressAndSend(dir, c)
}

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

func GetFileList(dir string, commitTag string) []string {
	files, err := os.Open(filepath.Join(dir, ".obj", commitTag, "hashmap.txt")); if err != nil {
		log.Fatal(err)
	}; 
	
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

func MergeMethod(base string, mine string, other string) Data {
    m, _, _ := merge.Merge(base, mine, other)
    conflicts := 0
	flag := 0

	if m == "conflict" {
		conflicts++
	}

	if m == "a" {
		flag++
	}

    // merged := strings.Join(m.merge_lines(`start_marker = "\n!!!--Conflict--!!!\n!--Your version--", mid_marker = "\n!--Other version--", end_marker = "\n!--End conflict--\n`), "")
	merged := ""
    r := Data{mdContent: merged, conflict: conflicts, merged: flag}

	return r
}

func GetFileName(objectdir string, commitTag string) []string {
	files, err := os.Open(filepath.Join(objectdir, commitTag)); if err != nil {
        log.Fatal(err)
    }; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	var list []string
	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)
		_ = append(list, p[0])
	}
	
	return list
}

func GetFileLocation(objectdir string, commitTag string, filename string) string {
	files, err := os.Open(filepath.Join(objectdir, commitTag)); if err != nil {
        log.Fatal(err)
    }; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)

		if p[1] == filename {
			return p[0]
		}
	}

	return ""
}

func UpdateModifyTime(trackingFile string) {
	files, err := os.Open(trackingFile); if err != nil {
        log.Fatal(err)
    }; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	fp, err := os.Create(trackingFile); if err != nil {
		panic(err)
	}; 
	
	defer fp.Close()

	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)

		name := p[0]
		status := p[1]

		_, err = fp.WriteString(name + " " + status + " " + time.Now().String() + "\n"); if err != nil {
			panic(err)
		}
	}
}

func GetLastCommit(commits []string) string { // not implemented
	return ""
}