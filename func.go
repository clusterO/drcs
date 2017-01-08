package dcrs

import (
	"fmt"
  "os"
  "log"
)

func initialize(directoryPath string) {
	fmt.Print("local functionalities of rcs")

	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
    log.Fatal(err)
	}
}

func add(filename string) {}

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