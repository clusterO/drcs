package dcrs

import (
	network "DCRS/net"
	equalfile "EqualFile"
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Cli() {
	var (
		init    bool
		logging bool
		status  bool
		add     string
		commit  string
		diff    string
		pull    string
		revert  string
		push    string
		clone   string
	)

	flag.BoolVar(&init, "init", false, "Create an empty repository or reinitialize an existing one")
	flag.BoolVar(&logging, "log", false, "Show commit logs")
	flag.BoolVar(&status, "status", false, "Show the working tree status")
	flag.StringVar(&add, "add", "", "Add file contents to the index")
	flag.StringVar(&commit, "commit", "", "Record changes to the repository")
	flag.StringVar(&diff, "diff", "", "Show changes between commits")
	flag.StringVar(&pull, "pull", "", "Fetch and merge commits and files")
	flag.StringVar(&revert, "revert", "", "Revert current directory to an old commit")
	flag.StringVar(&push, "push", "", "Update remote refs along with associated objects")
	flag.StringVar(&clone, "clone", "", "Clone remote repository")
	flag.Parse()

	path, err := os.Getwd(); if err != nil {
		log.Fatal(err)
	}

	dir, err := filepath.Abs(path); if err != nil {
		log.Fatal(err)
	}

	switch {
		case init:
			Init(dir, flag.Arg(0))
		case add != "":
			Add(add, dir)
		case commit != "":
			Commit(commit, dir)
		case status:
			Status(dir)
		case logging:
			Log(path)
		case diff != "":
			diffs := strings.Fields(diff)
			Diff(dir, diffs[0], diffs[1])
		case pull != "":
			Pull(pull, dir)
		case revert != "":
			commitHash := revert
			hashMap := revert
			Revert(commitHash, hashMap, path)
		case push != "":
			Push(push, dir)
		case clone != "":
			Clone(clone, dir)
	}
}

func Init(path string, arg string) {
	dir := filepath.Join(path, ".obj")
	err := os.MkdirAll(dir, os.ModePerm); if err != nil {
		log.Fatal(err)
	}

	config, err := os.Create(filepath.Join(dir, "config.txt")); if err != nil {
		panic(err)
	};
	
	defer config.Close()

	if arg == "y" {
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
	}; 
	
	defer files.Close()
}

func Add(filename string, dir string) {
	// works only with a given filename
	// should consider multuple files and . for all files 
	var list []string
	file, err := os.Stat(filepath.Join(dir, filename)); if err != nil {
		fmt.Println(err)
		return
	}

	tracker := filepath.Join(dir, ".obj", "tracker.txt")
	source, err := os.Open(tracker); if err != nil {
		fmt.Println(err)
		return
	}

	data, _ := os.ReadFile(tracker)
	hash := sha1.New()
	hash.Write(data)
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	destination, err := os.Create(filepath.Join(dir, ".obj", sha)); if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(destination, source)
	source.Close()
	destination.Close()

	switch mode := file.Mode(); {
		case mode.IsRegular():
			isFound := false
			p := filepath.Join(dir, ".obj", filename)
			files, err := os.Open(filepath.Join(dir, ".obj", sha)); if err != nil {
				log.Fatal(err)
			};

			defer files.Close()

			scanner := bufio.NewScanner(files)
			fo, err := os.Create(tracker); if err != nil {
				panic(err)
			}; 
			
			defer fo.Close()

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
		err := filepath.Walk(filepath.Join(dir, filename), func(p string, _ os.FileInfo, _ error) error {
			if p != filepath.Join(dir, filename) {
				list = append(list, p)
			}
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
				case mode.IsRegular(): // DRY
					isFound := false
					p, err := filepath.Abs(fi.Name()); if err != nil {
						log.Fatal(err)
					}

					files, err := os.Open(filepath.Join(dir, ".obj", sha)); if err != nil {
						log.Fatal(err)
					}; 
					
					defer files.Close()
					scanner := bufio.NewScanner(files)
					fo, err := os.Create(tracker); if err != nil {
						panic(err)
					}; 
					
					defer fo.Close()

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

func Commit(message string, dir string) (int64, error) {
	files, err := os.Open(filepath.Join(dir, ".obj", "config.txt")); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()
	date := time.Now().String()
	hash := sha1.New()
	hash.Write([]byte(username + date))
	hashmap := hash.Sum(nil)
	tracker := filepath.Join(dir, ".obj", "tracker.txt")
	files, err = os.Open(tracker); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()

	scanner = bufio.NewScanner(files)
	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "uncommitted" || path[1] == "committed" {
			src := path[0]
			dst := filepath.Join(dir, ".obj", base64.URLEncoding.EncodeToString(hashmap))

			_, err = os.Stat(src); if err != nil {
				return 0, err
			}

			source, err := os.Open(src); if err != nil {
				return 0, err
			}; 
			
			defer source.Close()
			
			destination, err := os.Create(dst); if err != nil {
				return 0, err
			}; 
			
			defer destination.Close()

			_, _ = io.Copy(destination, source)
		}
	}

	filePath := filepath.Join(dir, ".obj", string(hashmap), "message.txt")
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); if err != nil {
    	panic(err)
	}

	fo, err := os.Create(filePath); if err != nil {
    	panic(err)
	}
	
	defer fo.Close()

	r, err := fo.WriteString(message); if err != nil {
		fmt.Println(r)
		panic(err)
	}

	fo, err = os.Create(tracker); if err != nil {
		panic(err)
	}; 
	
	defer fo.Close()

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

func Rename(directoryPath string, newName string) {
	if _, err := os.Stat(newName); err == nil {
		fmt.Println(newName, " directory exists")
	} else if os.IsNotExist(err) {
		err = os.Rename(directoryPath, newName); if err != nil {
			log.Fatal(err)
		}
	}
}

func Clone(target string, dir string) {
	directory := filepath.Base(target)
	Init(directory, "")
	Pull(target, dir)
}

func Log(dir string) {
	files, err := os.Open(filepath.Join(dir, ".obj", "config.txt")); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	scanner := bufio.NewScanner(files)
	username := scanner.Text()
	dst := filepath.Join(dir, ".obj")
	var list []string
	err = filepath.Walk(dst, func(p string, _ os.FileInfo, _ error) error {
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

func Diff(dir string, commitx string, commity string) bool {
	filex := filepath.Join(dir, ".obj", commitx)
	filey := filepath.Join(dir, ".obj", commity)
	cmp := equalfile.New(nil, equalfile.Options{})
	equal, _ := cmp.CompareFile(filex, filey)

	return equal
}

func Status(dir string) {
	files, err := os.Open(filepath.Join(dir, ".obj", "tracker.txt")); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	for scanner.Scan() {
		line := scanner.Text()
		path := strings.Fields(line)

		if path[1] == "uncommitted" {
			fmt.Println(line)
		}
	}
}

func Pull(url string, dir string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	packageName := strings.Split(strings.Split(url, ":")[1], "/")[1]
	network.Connect(ip, port, dir, packageName, true)
}

func Push(url string, dir string) {
	ip := strings.Split(url, ":")[0]
	port := strings.Split(strings.Split(url, ":")[1], "/")[0]
	packageName := strings.Split(strings.Split(url, ":")[1], "/")[1]
	network.Connect(ip, port, dir, packageName, true)
}

func Revert(commitHash string, hashMap string, dir string) {
	// verify the params
	files, err := os.Open(filepath.Join(dir, ".obj", commitHash, hashMap)); if err != nil {
		log.Fatal(err)
	}; 
	
	defer files.Close()
	
	scanner := bufio.NewScanner(files)
	for scanner.Scan() {
		line := scanner.Text()
		p := strings.Fields(line)
		filename := p[0]
		content := GetFile(dir, commitHash, filename)
		dumpFile := filepath.Join(dir, ".obj", "dump.txt")
		fo, err := os.Create(dumpFile); if err != nil {
			panic(err)
		}; 
		
		defer fo.Close()
		
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
		}; 
		
		defer source.Close()
		
		destination, err := os.Create(dst); if err != nil {
			return
		}; 
		
		defer destination.Close()

		io.Copy(destination, source)
	}
}

func Merge(directory string, dir string) { 	// not processed
	commits := GetAllCommits(directory)
	mycommits := GetAllCommits(dir)
	commitsToFetch := Difference(commits, mycommits)
	f, err := os.OpenFile(filepath.Join(dir, "status.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); if err != nil {
		log.Println(err)
	};
	
	defer f.Close()

	for i := 0; i < len(commitsToFetch); i++ {
		k := strings.Fields(commitsToFetch[i])
		committag := k[1]
		print(committag)
		content := GetCommitsContent(directory, committag)
		UncompressAndWrite(dir, committag, content)
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
	myFileList := GetFileList(dir, strings.Fields(mycommits[0])[1])
	otherFileList := GetFileList(dir, strings.Fields(commits[0])[1])
	flag := 0

	for elem := range myFileList {
		for temp := range otherFileList {
			if myFileList[elem] == otherFileList[temp] {
				dicts := MergeMethod(GetFile(dir, strings.Fields(parentCommit)[1], myFileList[elem]), GetFile(dir, strings.Fields(mycommits[0])[1], myFileList[elem]), GetFile(dir, strings.Fields(commits[0])[1], myFileList[elem]))
				fo, err := os.Create(filepath.Join(dir, myFileList[elem])); if err != nil {
					panic(err)
				}; 
				
				defer fo.Close()
				
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
				}; 
				
				defer files.Close()
				
				content := GetFile(dir, commits[0], myFileList[elem])
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
