package main

import (
	_ "goj/boot"
	_ "goj/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
