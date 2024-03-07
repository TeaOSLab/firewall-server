// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package main

import (
	"flag"
	"fmt"
	"github.com/TeaOSLab/firewall-server/internal/servers"
	"net"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "", "such as 127.0.0.1:2345")
	flag.Parse()

	if len(addr) == 0 {
		fmt.Println("'-addr' option required")
		return
	}

	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("invalid '-addr' option")
		return
	}

	fmt.Println("starting '" + addr + "' ...")

	var server = servers.NewServer(addr)
	err = server.Listen()
	if err != nil {
		fmt.Println("[ERROR]" + err.Error())
		return
	}
}
