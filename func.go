package dcrs

import (
	"fmt"
)

func init() {
	fmt.Print("local functionalities of rcs")
}

func add() {}
func clone() {}
func log() {}
func diff() {}
func status() {}
func pull(url string) {}
func push(url string) {}
func revert(commitHash string) {}