package dcrs

import (
	"flag"
	"fmt"
)

func main() {
	var initialize = flag.String("init", "", "initialize the repo")
	var commit = flag.String("commit", "", "commit changes")
	
	flag.Parse()

	if *initialize != "" {
		fmt.Println("Initializing repo")
	} else if *commit != "" {
        print("Commiting with message " + *commit)
	} 
}

