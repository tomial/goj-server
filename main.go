package main

import (
	_ "goj-server/boot"
	_ "goj-server/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
