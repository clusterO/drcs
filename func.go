package dcrs

import (
	"fmt"
  "os"
  "log"
  "path/filepath"
  "bufio"
  "strings"
  "time"
  "crypto/sha1"
  "io"
  "flag"
  "github.com/udhos/equalfile"
)

func main() {
  var init string
  var add string
  var commit string
  var status string
  var log string
  var diff string

	flag.StringVar(&init, "init", "", "initialize the repo")
  flag.StringVar(&init, "i", "", "initialize the repo (shorthand)")

	flag.StringVar(&add, "add", "", "add files")
  flag.StringVar(&add, "a", "", "add files (shorthand)")

  flag.StringVar(&commit, "commit", "", "commit changes")
  flag.StringVar(&commit, "c", "", "commit changes (shorthand)")

  flag.StringVar(&status, "status", "", "show status")
  flag.StringVar(&status, "s", "", "show status (shorthand)")

  flag.StringVar(&log, "log", "", "list all commits")
  flag.StringVar(&log, "l", "", "list all commits (shorthand)")

  flag.StringVar(&diff, "diff", "", "overview of difference")
  flag.StringVar(&diff, "d", "", "overview of difference (shorthand)")
    
	
	flag.Parse()

	if init != "" {
		Init()
	} else if commit != "" {
    Commit(commit)
	} else if add != "" {
    Add(add)
    } else if status != "" {
        Status()
    } else if log != "" {
        Log()
    } else if log != "" {
        Diff(diff, diff)
    } 
}

// Init initialize repository
func Init() {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

	err = os.MkdirAll(path + "/dcrs", os.ModePerm)
	if err != nil {
    log.Fatal(err)
  }
  
  var username string
  println("enter username: ")
  fmt.Scan("%s", &username)

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
func Add(filename string) {
  var list []string

  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

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
          Add(f)
      }
    }
  }
}

// Commit changes
func Commit(message string) (int64, error) {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

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

    if(path[1] == "notcommited\n" || path[1] == "commited\n") {
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

  return 0, nil
}

// Rename a file
func Rename(directoryPath string, newName string) {
  if _, err := os.Stat(newName)
  err == nil {
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
func Clone() {}

// Log list all commits
func Log() {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

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
func Status() {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

  files, err := os.Open(path + "/dcrs/" + "files.txt")
  if err != nil {
    log.Fatal(err)
  }

  defer files.Close()
  scanner := bufio.NewScanner(files)

  for scanner.Scan() {
    line := scanner.Text()
    path := strings.Fields(line)
    
    if path[1] == "notcommited\n" {
      fmt.Println(line)
    }
  }
}

// Pull changes from repository
func Pull(url string) {}
// Push changes to repository
func Push(url string) {}

// Revert changes
func Revert(commitHash string, hashMap string) {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

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

func merge(directory string) {
  path, err := os.Getwd()
  if err != nil {
    log.Println(err)
  }

  Dircmp(path, directory)
}