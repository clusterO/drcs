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
)

func initialize(directoryPath string) {
	err := os.MkdirAll(directoryPath + "/init", os.ModePerm)
	if err != nil {
    log.Fatal(err)
  }
  
  var username string
  println("enter username: ")
  fmt.Scan("%s", &username)
}

func add(path string) {
  var list []string

  file, err := os.Stat(path)
  if err != nil {
      fmt.Println(err)
      return
  }

  switch mode := file.Mode(); {
    case mode.IsRegular():
      p, err := filepath.Abs(path)
      if err != nil {
        log.Fatal(err)
      }

      files, err := os.Open("files.txt")
      if err != nil {
	      log.Fatal(err)
      }

      defer files.Close()
      scanner := bufio.NewScanner(files)

      fo, err := os.Create("files.txt")
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
        err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
          list = append(list, path)
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
          p, err := filepath.Abs(path)
          if err != nil {
            log.Fatal(err)
          }

          files, err := os.Open("files.txt")
          if err != nil {
            log.Fatal(err)
          }

          defer files.Close()
          scanner := bufio.NewScanner(files)

          fo, err := os.Create("files.txt")
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
          add(f)
      }
    }
  }
}

func commit(message string) (int64, error) {
  username := ""
  dateandtime := time.Now().String()

  hash := sha1.New()
  hash.Write([]byte(username + dateandtime))
  hashmap := hash.Sum(nil)
  
  files, err := os.Open("files.txt")
  if err != nil {
    log.Fatal(err)
  }

  defer files.Close()
  scanner := bufio.NewScanner(files)

  for scanner.Scan() {
    line := scanner.Text()
    path := strings.Fields(line)

    if(path[1] == "commited\n") {
      src := path[0]
      p, err := filepath.Abs("init")
      dst := p + "/object/" + string(hashmap)

      sourceFileStat, err := os.Stat(src)
      if err != nil {
        return 0, err
      }

      if !sourceFileStat.Mode().IsRegular() {
        return 0, fmt.Errorf("%s is not a regular file", src)
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
      nBytes, err := io.Copy(destination, source)
      return nBytes, err
    }
  }

  p, err := filepath.Abs("init")
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

  return 0, nil
}

func renmae(directoryPath string, newName string) {
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

func clone() {}
func logging() {}
func diff() {}

func status() {
  files, err := os.Open("files.txt")
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

func pull(url string) {}
func push(url string) {}
func revert(commitHash string) {}