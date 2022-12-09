package main

import (
	"flag"
	"fmt"

	fm "wallet/pkg/file_manage"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "storage", "File where the password hash is stored")
	flag.Parse()

	err := fm.Access(fileName)
	if err != nil {
		fmt.Println("You don't have access to the wallet!!!")
		return
	}
	fmt.Println("Got access to the wallet!!!")
}
