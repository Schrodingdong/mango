package main

import (
	"os"

	"github.com/schrodingdong/mango/cmd"
)

func init() {
	os.Setenv("TZ", "Africa/Casablanca")
}

func main() {
	cmd.Execute()
}
