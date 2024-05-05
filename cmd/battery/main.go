package main

import (
	"github.com/JeffreySmith/battery"
	"os"
)

func main(){
	os.Exit(battery.OpenBSDMain())
}
