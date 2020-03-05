package dcrs

import (
	"fmt"
  "os"
  "log"
  "path/filepath"
)

var files []string

func initialize(directoryPath string) {
	fmt.Print("local functionalities of rcs")

	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
    log.Fatal(err)
	}
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
      files = append(files, p, "not commited")
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
          files = append(files, p, "not commited")
        case mode.IsDir():
          add(f)
      }
    }
  }
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
func status() {}
func pull(url string) {}
func push(url string) {}
func revert(commitHash string) {}