package dcrs

import (
	"bufio"
	"crypto/sha1"
	network "dcrs/net"
	"encoding/base64"
	"equalfile"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Cli() {
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
	flag.StringVar(&clone, "n", "", "clone remote repository (shorthand)")

	path, err := os.Getwd(); if err != nil {
		fmt.Println(err)
	}
	dir, err := filepath.Abs(path); if err != nil {
		log.Fatal(err)
	}
	dcrs := filepath.Join(dir, "dcrs")

	flag.Parse()

	if init != "" {
		Init(dir, init)
	} else if commit != "" {
		Commit(commit, dir)
	} else if add != "" {
		Add(add, dir)
	} else if status != "" {
		Status(status)
	} else if logging != "" {
		Log(path)
	} else if diff != "" {
		Diff(dcrs, strings.Fields(diff)[0], strings.Fields(diff)[1])
	} else if pull != "" {
		Pull(pull)
	} else if revert != "" {
		Revert(dcrs, revert, path)
	} else if push != "" {
		Push(push)
	} else if clone != "" {
		Clone(clone)
	}
}

// Init initialize repository
func Init(path string, flag  string) {
	dir := filepath.Join(path, "obj")
	err := os.MkdirAll(dir, os.ModePerm); if err != nil {
		log.Fatal(err)
	}

	config, err := os.Create(filepath.Join(dir, "config.txt")); if err != nil {
		panic(err)
	}; defer config.Close()

	if(flag == "y") {
		var username string
		var email string
		var repository string

		reader := bufio.NewReader(os.Stdin)
		println("Enter username: "); username, _ = reader.ReadString('\n')
		println("Enter email: "); email, _ = reader.ReadString('\n')
		println("Enter repository url: "); repository, _ = reader.ReadString('\n')

		_, err = config.WriteString("username: " + username); if err != nil {
			panic(err)
		}
		_, err = config.WriteString("email: " + email); if err != nil {
			panic(err)
		}
		_, err = config.WriteString("repository: " + repository); if err != nil {
			panic(err)
		}
	}

	files, err := os.Create(filepath.Join(dir, "tracker.txt")); if err != nil {
		panic(err)
	}; defer files.Close()
}

// Add files to tracking system
func Add(filename string, dir string) {
	var list []string
	file, err := os.Stat(filepath.Join(dir, "obj", filename)); if err != nil {
		fmt.Println(err)
		return
	}
	tracker := filepath.Join(dir, "obj", "tracker.txt")

	source, err := os.Open(tracker); if err != nil {
		fmt.Println(err)
		return
	}
	data, err := ioutil.ReadFile(tracker)
	hash := sha1.New()
	hash.Write(data)
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	destination, err := os.Create(filepath.Join(dir, "obj", sha)); if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(destination, source)
	source.Close()
	destination.Close()

	switch mode := file.Mode(); {
	case mode.IsRegular():
		isFound := false
		p := filepath.Join(dir, "obj", filename)

		files, err := os.Open(filepath.Join(dir, "obj", sha)); if err != nil {
			log.Fatal(err)
		}; defer files.Close()
		scanner := bufio.NewScanner(files)

		fo, err := os.Create(tracker); if err != nil {
			panic(err)
		}; defer fo.Close()

		for scanner.Scan() {
			line := scanner.Text()
			path := strings.Fields(line)

			if path[0] == p {
				isFound = true
				_, err := fo.WriteString("\n" + path[0] + " uncommitted\n"); if err != nil {
					panic(err)
				}
			} else {
				_, err := fo.WriteString(line); if err != nil {
					panic(err)
				}
			}
		}

		if(!isFound) {
			_, err := fo.WriteString("\n" + p + " uncommitted\n"); if err != nil {
				panic(err)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	case mode.IsDir():
		err := filepath.Walk(filepath.Join(dir, "obj", filename), func(p string, info os.FileInfo, err error) error {
			list = append(list, p)
			return nil
		}); if err != nil {
			panic(err)
		}

		for _, f := range list {
			fi, err := os.Stat(f); if err != nil {
				fmt.Println(err)
				return
			}

			switch mode := fi.Mode(); {
			case mode.IsRegular():
				isFound := false
				p, err := filepath.Abs(fi.Name()); if err != nil {
					log.Fatal(err)
				}

				files, err := os.Open(filepath.Join(dir, "obj", sha)); if err != nil {
					log.Fatal(err)
				}; defer files.Close()
				scanner := bufio.NewScanner(files)

				fo, err := os.Create(tracker); if err != nil {
					panic(err)
				}; defer fo.Close()

				for scanner.Scan() {
					line := scanner.Text()
					path := strings.Fields(line)

					if path[0] == p {
						isFound = true
						_, err := fo.WriteString("\n" + path[0] + " uncommitted\n\r"); if err != nil {
							panic(err)
						}
					} else {
						_, err := fo.WriteString(line); if err != nil {
							panic(err)
						}
					}
				}

				if(!isFound) {
					_, err := fo.WriteString("\n" + p + " uncommitted\n\r"); if err != nil {
						panic(err)
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}

			case mode.IsDir():
				Add(f, dir)
			}
		}
	}

	UpdateModifyTime(tracker)
}

// Commit changes
func Commit(message string, dir string) (int64, error) {
	files, err := os.Open(filepath.Join(dir, "obj", "config.txt")); if err != nil {
		log.Fatal(err)
	}; defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()
	date := time.Now().String()

	hash := sha1.New()
	hash.Write([]byte(username + date))
	hashmap := hash.Sum(nil)

	tracker := filepath.Join(dir, "obj", "tracker.txt")
	files, err = os.Open(tracker); if err != nil {
		log.Fatal(err)
	}; defer files.Close()
	scanner = bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "uncommitted" || path[1] == "committed" {
			src := path[0]
			dst := filepath.Join(dir, "obj", base64.URLEncoding.EncodeToString(hashmap))

			_, err = os.Stat(src); if err != nil {
				return 0, err
			}
			source, err := os.Open(src); if err != nil {
				return 0, err
			}; defer source.Close()
			destination, err := os.Create(dst); if err != nil {
				return 0, err
			}; defer destination.Close()

			_, err = io.Copy(destination, source)
		}
	}

	fo, err := os.Create(filepath.Join(dir, string(hashmap), "message.txt")); if err != nil {
		panic(err)
	}; defer fo.Close()
	r, err := fo.WriteString(message); if err != nil {
		fmt.Println(r)
		panic(err)
	}
	fo, err = os.Create(tracker); if err != nil {
		panic(err)
	}; defer fo.Close()

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)
		_, err := fo.WriteString(path[0] + "committed\n"); if err != nil {
			 panic(err)
		}
	}

	UpdateModifyTime(tracker)
	return 0, nil
}

// Rename a file
func Rename(directoryPath string, newName string) {
	if _, err := os.Stat(newName); err == nil {
		fmt.Println(newName, " directory exists")
	} else if os.IsNotExist(err) {
		err = os.Rename(directoryPath, newName); if err != nil {
			log.Fatal(err)
		}
	}
}

// Clone a repository
func Clone(target string) {
	directory := filepath.Base(target)
	Init(directory, "")
	Pull(target)
}

// Log list all commits
func Log(dir string) {
	files, err := os.Open(filepath.Join(dir, "obj", "config.txt")); if err != nil {
		log.Fatal(err)
	}; defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()

	dst := filepath.Join(dir, "obj")
	var list []string
	err = filepath.Walk(dst, func(p string, info os.FileInfo, err error) error {
		list = append(list, p)
		return nil
	}); if err != nil {
		panic(err)
	}

	for _, f := range list {
		fi, err := os.Stat(f); if err != nil {
			fmt.Println(err)
			return
		}

		moddate := fi.ModTime()
		fmt.Println("Commit tag: ", fi.Name())
		fmt.Println("Author: ", username)
		fmt.Println("Timestamp: ", moddate)
		fmt.Println()
	}
}

// Diff show difference between two versions
func Diff(dir string, commitx string, commity string) bool {
	filex := filepath.Join(dir, "obj", commitx)
	filey := filepath.Join(dir, "obj", commity)
	cmp := equalfile.New(nil, equalfile.Options{})
	equal, _ := cmp.CompareFile(filex, filey)

	return equal
}

// Status show project status
func Status(dir string) {
	files, err := os.Open(filepath.Join(dir, "tracker.txt")); if err != nil {
		log.Fatal(err)
	}; defer files.Close()
	scanner := bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "uncommitted" {
			fmt.Println(line)
		}
	}
}

// Pull changes from repository
func Pull(url string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	directory := strings.Split(strings.Split(url, ":")[1], "/")[1]
	network.Connect(ip, port, directory, true)
}

// Push changes to repository
func Push(url string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	directory := strings.Split(strings.Split(url, ":")[1], "/")[1]
	network.Connect(ip, port, directory, false)
}

// Revert changes
func Revert(commitHash string, hashMap string, dir string) {
	files, err := os.Open(filepath.Join(dir, "obj", commitHash, hashMap)); if err != nil {
		log.Fatal(err)
	}; defer files.Close()
	scanner := bufio.NewScanner(files)

	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)
		filename := p[0]

		content := GetFile(commitHash, filename)
		dumpFile := filepath.Join(dir, "dump.txt")
		fo, err := os.Create(dumpFile); if err != nil {
			panic(err)
		}; defer fo.Close()
		_, err = fo.WriteString(content); if err != nil {
			panic(err)
		}

		src := dumpFile
		dst := dir

		_, err = os.Stat(src); if err != nil {
			return
		}
		source, err := os.Open(src); if err != nil {
			return
		}; defer source.Close()
		destination, err := os.Create(dst); if err != nil {
			return
		}; defer destination.Close()

		io.Copy(destination, source)
	}
}

// Merge files
func Merge(directory string, dir string) {
	commits := GetCommits(directory)
	print(commits)
	mycommits := GetAllCommits()
	commitsToFetch := Difference(commits, mycommits)

	f, err := os.OpenFile(filepath.Join(dir, "status.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); if err != nil {
		log.Println(err)
	}; defer f.Close()

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
				fo, err := os.Create(filepath.Join(dir, myFileList[elem])); if err != nil {
					panic(err)
				}; defer fo.Close()
				_, err = fo.WriteString(dicts.mdContent); if err != nil {
					panic(err)
				}
				Add(myFileList[elem], dir)

				if dicts.conflict == 1 {
					print("Merged with conflicts in " + myFileList[elem] + " not commiting.Please commit after manually changing")
					flag = 1
				} else {
					print("Merged " + myFileList[elem] + " successfully\n")
				}
			} else {
				files, err := os.Create(myFileList[elem]); if err != nil {
					panic(err)
				}; defer files.Close()
				content := GetFile(commits[0], myFileList[elem])
				_, err = files.WriteString(content); if err != nil {
					panic(err)
				}
				Add(myFileList[elem], dir)
			}
		}
	}

	if flag == 0 {
		Commit("auto-merged successfull", dir)
	}
}
