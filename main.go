package main

import (
	_ "lanvs/routers"
	"github.com/astaxie/beego"
)

const (
	PORT = ":4000"
)

func main() {
	beego.Run(PORT)
}

