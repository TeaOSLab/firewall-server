// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	teaconst "github.com/TeaOSLab/firewall-server/internal/const"
	"github.com/TeaOSLab/firewall-server/internal/servers"
	"github.com/iwind/TeaGo/maps"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show version")

	var addr string
	flag.StringVar(&addr, "addr", "", "such as 127.0.0.1:12345")

	flag.Parse()

	// version
	if showVersion {
		fmt.Println("version: " + teaconst.Version)
		return
	}

	exe, _ := os.Executable()

	// restore from local file
	var dataFile = filepath.Dir(exe) + "/firewall.config"
	cachedData, _ := os.ReadFile(dataFile)
	var cacheMap = maps.Map{}
	if len(cachedData) > 0 {
		err := json.Unmarshal(cachedData, &cacheMap)
		if err == nil {

		}
	}

	if len(addr) == 0 {
		var cachedAddr = cacheMap.GetString("addr")
		if len(cachedAddr) > 0 {
			addr = cachedAddr
		}

		// check again
		if len(addr) == 0 {
			fmt.Println("'-addr' option required")
			return
		}
	} else {
		cacheMap["addr"] = addr
		_ = os.WriteFile(dataFile, cacheMap.AsJSON(), 0666)
	}

	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("invalid '-addr' option")
		return
	}

	fmt.Println("starting '" + addr + "' ...")

	var server = servers.NewServer(addr)
	go func() {
		time.Sleep(3 * time.Second)
		_ = server.InstallService()
	}()
	err = server.Listen()
	if err != nil {
		fmt.Println("[ERROR]" + err.Error())
		return
	}
}
