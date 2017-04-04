package dcrs

import (
	"bufio"
	"crypto/sha1"
	network "dcrs/net"
	"equalfile"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main(dirc string) {
	var init string
	var add string
	var commit string
	var status string
	var logging string
	var diff string
	var pull string
	var revert string
	var push string
	var clone string

	flag.StringVar(&init, "init", "", "initialize the repo")
	flag.StringVar(&init, "i", "", "initialize the repo (shorthand)")

	flag.StringVar(&add, "add", "", "add files")
	flag.StringVar(&add, "a", "", "add files (shorthand)")

	flag.StringVar(&commit, "commit", "", "commit changes")
	flag.StringVar(&commit, "c", "", "commit changes (shorthand)")

	flag.StringVar(&status, "status", "", "show status")
	flag.StringVar(&status, "s", "", "show status (shorthand)")

	flag.StringVar(&logging, "log", "", "list all commits")
	flag.StringVar(&logging, "l", "", "list all commits (shorthand)")

	flag.StringVar(&diff, "diff", "", "overview of difference")
	flag.StringVar(&diff, "d", "", "overview of difference (shorthand)")

	flag.StringVar(&pull, "pull", "", "pull and merge commits and files")
	flag.StringVar(&pull, "p", "", "pull and merge commits and files (shorthand)")

	flag.StringVar(&revert, "revert", "", "revert current directory to an old commit")
	flag.StringVar(&revert, "r", "", "revert current directory to an old commit (shorthand)")

	flag.StringVar(&push, "push", "", "pushes the commits")
	flag.StringVar(&push, "ps", "", "pushes the commits (shorthand)")

	flag.StringVar(&clone, "clone", "", "clone remote repository")
	flag.StringVar(&clone, "c", "", "clone remote repository (shorthand)")

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	directory := dirc
	statusfile, _ := filepath.Abs(filepath.Join(directory, "dcrs", "status.txt"))
	userfile, _  := filepath.Abs(filepath.Join(directory, "dcrs", "username.txt"))
	trackingfile, _ := filepath.Abs(filepath.Join(directory, "dcrs", "files.txt"))
	objectdir, _  := filepath.Abs(filepath.Join(directory, "dcrs", "object"))
	dcrs, _  := filepath.Abs(filepath.Join(directory, "dcrs"))
	hashmap, _  := filepath.Abs(filepath.Join(directory, "dcrs", "object", "commitfiles"))
	commitfiles, _  := filepath.Abs(filepath.Join(directory, "dcrs", "object", "commitfiles"))

	flag.Parse()

	if init != "" {
		Init(userfile)
	} else if commit != "" {
		Commit(commit, trackingfile)
	} else if add != "" {
		Add(add, path)
	} else if status != "" {
		Status(statusfile)
	} else if logging != "" {
		Log(path)
	} else if diff != "" {
		Diff(diff, path)
	} else if pull != "" {
		Pull(objectdir)
	} else if revert != "" {
		Revert(dcrs, hashmap, path)
	} else if push != "" {
		Push(commitfiles)
	} else if clone != "" {
		Clone(clone)
	}
}

// Init initialize repository
func Init(path string) {
	err := os.MkdirAll(path+"/dcrs", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	var username string
	var email string
	println("enter username: ")
	fmt.Scan("%s", &username)
	println("enter email: ")
	fmt.Scan("%s", &email)

	config, err := os.Create(path + "/dcrs/" + "config.txt")
	if err != nil {
		panic(err)
	}

	defer config.Close()

	_, err = config.WriteString(username)
	if err != nil {
		panic(err)
	}

	files, err := os.Create(path + "/dcrs/" + "files.txt")
	if err != nil {
		panic(err)
	}

	defer files.Close()
}

// Add files to tracking system
func Add(filename string, path string) {
	var list []string

	file, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mode := file.Mode(); {
	case mode.IsRegular():
		p, err := filepath.Abs(filename)
		if err != nil {
			log.Fatal(err)
		}

		files, err := os.Open(path + "/dcrs/" + "files.txt")
		if err != nil {
			log.Fatal(err)
		}

		defer files.Close()
		scanner := bufio.NewScanner(files)

		fo, err := os.Create(path + "/dcrs/" + "files.txt")
		if err != nil {
			panic(err)
		}

		defer fo.Close()

		for scanner.Scan() {
			line := scanner.Text()
			path := strings.Fields(line)

			if path[0] == p {
				_, err := fo.WriteString(path[0] + " notcommited\n")
				if err != nil {
					panic(err)
				}
			} else {
				_, err := fo.WriteString(line)
				if err != nil {
					panic(err)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	case mode.IsDir():
		err := filepath.Walk(filename, func(p string, info os.FileInfo, err error) error {
			list = append(list, p)
			return nil
		})

		if err != nil {
			panic(err)
		}

		for _, f := range list {
			fi, err := os.Stat(f)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch mode := fi.Mode(); {
			case mode.IsRegular():
				p, err := filepath.Abs(fi.Name())
				if err != nil {
					log.Fatal(err)
				}

				files, err := os.Open(path + "/dcrs/" + "files.txt")
				if err != nil {
					log.Fatal(err)
				}

				defer files.Close()
				scanner := bufio.NewScanner(files)

				fo, err := os.Create(path + "/dcrs/" + "files.txt")
				if err != nil {
					panic(err)
				}

				defer fo.Close()

				for scanner.Scan() {
					line := scanner.Text()
					path := strings.Fields(line)

					if path[0] == p {
						_, err := fo.WriteString(path[0] + " notcommited\n")
						if err != nil {
							panic(err)
						}
					} else {
						_, err := fo.WriteString(line)
						if err != nil {
							panic(err)
						}
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			case mode.IsDir():
				Add(f, path)
			}
		}
	}

	UpdateModifyTime("")
}

// Commit changes
func Commit(message string, path string) (int64, error) {
	files, err := os.Open(path + "/dcrs/" + "config.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()
	dateandtime := time.Now().String()

	hash := sha1.New()
	hash.Write([]byte(username + dateandtime))
	hashmap := hash.Sum(nil)

	files, err = os.Open(path + "/dcrs/" + "files.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner = bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "notcommited" || path[1] == "commited" {
			src := path[0]
			p, err := filepath.Abs("dcrs")
			dst := p + "/object/" + string(hashmap)

			_, err = os.Stat(src)
			if err != nil {
				return 0, err
			}

			source, err := os.Open(src)
			if err != nil {
				return 0, err
			}

			defer source.Close()

			destination, err := os.Create(dst)
			if err != nil {
				return 0, err
			}

			defer destination.Close()
			_, err = io.Copy(destination, source)
		}
	}

	p, err := filepath.Abs("dcrs")
	fo, err := os.Create(p + "/object/" + string(hashmap) + "/message.txt")
	if err != nil {
		panic(err)
	}

	defer fo.Close()

	r, err := fo.WriteString(message)
	if err != nil {
		fmt.Println(r)
		panic(err)
	}

	fo, err = os.Create(path + "/dcrs/" + "files.txt")
	if err != nil {
		panic(err)
	}

	defer fo.Close()

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		_, err := fo.WriteString(path[0] + "commited\n")
		if err != nil {
			panic(err)
		}
	}

	UpdateModifyTime("")
	return 0, nil
}

// Rename a file
func Rename(directoryPath string, newName string) {
	if _, err := os.Stat(newName); err == nil {
		fmt.Println(newName, " directory exists")
	} else if os.IsNotExist(err) {
		err = os.Rename(directoryPath, newName)
		if err != nil {
			log.Fatal(err)
		}
	} else {

	}
}

// Clone a repository
func Clone(target string) {
	directory := filepath.Base(target)
	Init(directory)
	Pull(target)
}

// Log list all commits
func Log(path string) {
	files, err := os.Open(path + "/dcrs/" + "config.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()

	p, err := filepath.Abs("dcrs")
	dst := p + "/object/"
	var list []string
	err = filepath.Walk(dst, func(p string, info os.FileInfo, err error) error {
		list = append(list, p)
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, f := range list {
		fi, err := os.Stat(f)
		if err != nil {
			fmt.Println(err)
			return
		}

		moddate := fi.ModTime()
		fmt.Println("Commit tag: ", fi.Name())
		fmt.Println("Author: ", username)
		fmt.Println("Time Stamp: ", moddate)
		fmt.Println()
	}
}

// Diff show difference between two versions
func Diff(oldCommit string, newCommit string) bool {
	oldFile, _ := filepath.Abs("dcrs/object/" + oldCommit)
	newFile, _ := filepath.Abs("dcrs/object/" + newCommit)
	cmp := equalfile.New(nil, equalfile.Options{})
	equal, _ := cmp.CompareFile(oldFile, newFile)

	return equal
}

// Status show project status
func Status(path string) {
	files, err := os.Open(path + "/dcrs/" + "files.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "notcommited" {
			fmt.Println(line)
		}
	}
}

// Pull changes from repository
func Pull(url string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	directory := strings.Split(url, ":")[1][4:len(url)]
	network.Connect(ip, port, directory, true)
}

// Push changes to repository
func Push(url string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	directory := strings.Split(url, ":")[1][4:len(url)]
	network.Connect(ip, port, directory, false)
}

// Revert changes
func Revert(commitHash string, hashMap string, path string) {
	files, err := os.Open(path + "/dcrs/object/" + commitHash + "/" + hashMap)
	if err != nil {
		log.Fatal(err)
	}

	defer files.Close()
	scanner := bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)
		filename := p[0]

		content := GetFile(commitHash, filename)
		dumpFile, err := filepath.Abs("dump.txt")
		fo, err := os.Create(dumpFile)
		if err != nil {
			panic(err)
		}

		defer fo.Close()

		_, err = fo.WriteString(content)
		if err != nil {
			panic(err)
		}

		src := dumpFile
		dst := path

		_, err = os.Stat(src)
		if err != nil {
			return
		}

		source, err := os.Open(src)
		if err != nil {
			return
		}

		defer source.Close()

		destination, err := os.Create(dst)
		if err != nil {
			return
		}

		defer destination.Close()
		_, err = io.Copy(destination, source)
	}
}

// Merge files
func Merge(directory string) {
	commits := GetCommits(directory)
	print(commits)
	mycommits := GetAllCommits()
	commitsToFetch := Difference(commits, mycommits)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	f, err := os.OpenFile(path+"/dcrs/"+"status.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for i := 0; i < len(commitsToFetch); i++ {
		k := strings.Fields(commitsToFetch[i])
		committag := k[1]
		print(committag)
		content := GetCommitsContent(directory, committag)
		UncompressAndWrite(committag, content)
		print(commitsToFetch[i])

		if _, err := f.WriteString(commitsToFetch[i] + "\n"); err != nil {
			log.Println(err)
		}
	}

	parentCommit := ""
	for i := 0; i < len(commits); i++ {
		for _, v := range mycommits {
			if v == commits[i] {
				parentCommit = commits[i]
			}
		}
	}

	// parentFileList := GetFileList(strings.Fields(parentCommit)[1])
	myFileList := GetFileList(strings.Fields(mycommits[0])[1])
	otherFileList := GetFileList(strings.Fields(commits[0])[1])
	flag := 0

	for elem := range myFileList {
		for temp := range otherFileList {
			if myFileList[elem] == otherFileList[temp] {
				dicts := MergeMethod(GetFile(strings.Fields(parentCommit)[1], myFileList[elem]), GetFile(strings.Fields(mycommits[0])[1], myFileList[elem]), GetFile(strings.Fields(commits[0])[1], myFileList[elem]))

				fo, err := os.Create(path + "/dcrs/" + myFileList[elem])
				if err != nil {
					panic(err)
				}

				defer fo.Close()

				_, err = fo.WriteString(dicts.mdContent)
				if err != nil {
					panic(err)
				}

				Add(myFileList[elem], path)

				if dicts.conflict == 1 {
					print("Merged with conflicts in " + myFileList[elem] + " not commiting.Please commit after manually changing")
					flag = 1
				} else {
					print("Merged " + myFileList[elem] + " successfully\n")
				}
			} else {
				files, err := os.Create(myFileList[elem])
				if err != nil {
					panic(err)
				}

				defer files.Close()

				content := GetFile(commits[0], myFileList[elem])

				_, err = files.WriteString(content)
				if err != nil {
					panic(err)
				}

				Add(myFileList[elem], path)
			}
		}

	}

	if flag == 0 {
		Commit("auto-merged successfull", path)
	}
}
