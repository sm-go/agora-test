package main

import (
	"fmt"
	"os"
)

func main() {
	os.Setenv("AppName", "this is my app")
	shell, ok := os.LookupEnv("AppName")
	if !ok {
		fmt.Println("the env is not set")
	} else {
		fmt.Println("shell >>", shell)
	}
}
