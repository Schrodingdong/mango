package main

import (
	"github.com/schrodingdong/mango/cmd"
	"github.com/schrodingdong/mango/utils"
)

func init() {
	utils.CreateConfigFile()
}

func main() {
	cmd.Execute()
}
