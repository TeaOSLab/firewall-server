// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package servers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaOSLab/firewall-server/internal/firewalls"
	executils "github.com/TeaOSLab/firewall-server/internal/utils/exec"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	addr     string
	firewall firewalls.FirewallInterface
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (this *Server) Listen() error {
	var mux = &http.ServeMux{}

	this.firewall = firewalls.Firewall()

	mux.HandleFunc("/name", this.handleName)
	mux.HandleFunc("/isReady", this.handleIsReady)
	mux.HandleFunc("/isMock", this.handleIsMock)
	mux.HandleFunc("/allowPort", this.handleAllowPort)
	mux.HandleFunc("/removePort", this.handleRemovePort)
	mux.HandleFunc("/rejectSourceIP", this.handleRejectSourceIP)
	mux.HandleFunc("/dropSourceIP", this.handleDropSourceIP)
	mux.HandleFunc("/removeSourceIP", this.handleRemoveSourceIP)

	var httpServer = &http.Server{
		Addr:    this.addr,
		Handler: mux,
	}
	return httpServer.ListenAndServe()
}

func (this *Server) InstallService() error {
	systemd, err := executils.LookPath("systemctl")
	if err != nil {
		return nil
	}

	shPath, err := executils.LookPath("sh")
	if err != nil {
		return nil
	}

	exePath, _ := os.Executable()
	if len(exePath) == 0 {
		return errors.New("can not find executable path")
	}

	var systemdServiceFile = "/etc/systemd/system/edge-firewall-server.service"


	serviceData, err := os.ReadFile(systemdServiceFile)
	if err == nil && len(serviceData) > 0 && (bytes.Contains(serviceData, []byte("ExecStart="+exePath)) || bytes.Contains(serviceData, []byte("ExecStart=" + shPath + " "+exePath + ".sh"))) {
		return nil
	}

	err = os.WriteFile(exePath + ".sh", []byte(`#!/usr/bin/env bash

` + exePath), 0733)
	if err != nil {
		return fmt.Errorf("create boot shell file failed: %w", err)
	}

	var shortName = "edge-firewall-server"
	var longName = "edge-firewall-server"

	var desc = `### BEGIN INIT INFO
# Provides:          ` + shortName + `
# Required-Start:     $local_fs $network
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:
# Short-Description: ` + longName + ` Service
### END INIT INFO

[Unit]
Description=` + longName + ` Service
Before=shutdown.target
After=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=1s
ExecStart=` + shPath + " " + exePath + `.sh

[Install]
WantedBy=multi-user.target`

	// write file
	err = os.WriteFile(systemdServiceFile, []byte(desc), 0777)
	if err != nil {
		return err
	}

	// stop current systemd service if running
	_ = executils.NewTimeoutCmd(10*time.Second, systemd, "stop", shortName+".service").Start()

	// reload
	_ = executils.NewTimeoutCmd(10*time.Second, systemd, "daemon-reload").Start()

	// enable
	var cmd = executils.NewTimeoutCmd(10*time.Second, systemd, "enable", shortName+".service")
	cmd.WithStderr()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%w: %s", err, cmd.Stderr())
	}
	return nil
}

func (this *Server) handleName(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}
	this.write(writer, http.StatusOK, maps.Map{
		"name": this.firewall.Name(),
	})
}

func (this *Server) handleIsReady(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}
	this.write(writer, http.StatusOK, maps.Map{
		"result": this.firewall.IsReady(),
	})
}

func (this *Server) handleIsMock(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}
	this.write(writer, http.StatusOK, maps.Map{
		"result": this.firewall.IsMock(),
	})
}

func (this *Server) handleAllowPort(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}

	var port = types.Int(req.URL.Query().Get("port"))
	var protocol = req.URL.Query().Get("protocol")
	if port <= 0 {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}
	if protocol != "tcp" && protocol != "udp" {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}
	err := this.firewall.AllowPort(port, protocol)
	if err != nil {
		this.write(writer, http.StatusInternalServerError, nil)
		return
	}

	this.write(writer, http.StatusOK, nil)
}

func (this *Server) handleRemovePort(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}

	var port = types.Int(req.URL.Query().Get("port"))
	var protocol = req.URL.Query().Get("protocol")
	if port <= 0 {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}
	if protocol != "tcp" && protocol != "udp" {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}
	err := this.firewall.RemovePort(port, protocol)
	if err != nil {
		this.write(writer, http.StatusInternalServerError, nil)
		return
	}

	this.write(writer, http.StatusOK, nil)
}

func (this *Server) handleRejectSourceIP(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}

	var ip = req.URL.Query().Get("ip")
	if net.ParseIP(ip) == nil {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}

	var timeoutSeconds = types.Int(req.URL.Query().Get("timeoutSeconds"))
	if timeoutSeconds < 0 {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}

	err := this.firewall.RejectSourceIP(ip, timeoutSeconds)
	if err != nil {
		this.write(writer, http.StatusInternalServerError, nil)
		return
	}

	this.write(writer, http.StatusOK, nil)
}

func (this *Server) handleDropSourceIP(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}

	var ip = req.URL.Query().Get("ip")
	if net.ParseIP(ip) == nil {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}

	var timeoutSeconds = types.Int(req.URL.Query().Get("timeoutSeconds"))
	if timeoutSeconds < 0 {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}

	var async = req.URL.Query().Get("async") == "true"

	err := this.firewall.DropSourceIP(ip, timeoutSeconds, async)
	if err != nil {
		this.write(writer, http.StatusInternalServerError, nil)
		return
	}

	this.write(writer, http.StatusOK, nil)
}

func (this *Server) handleRemoveSourceIP(writer http.ResponseWriter, req *http.Request) {
	if this.firewall == nil {
		this.write(writer, http.StatusServiceUnavailable, nil)
		return
	}

	var ip = req.URL.Query().Get("ip")
	if net.ParseIP(ip) == nil {
		this.write(writer, http.StatusBadRequest, nil)
		return
	}

	err := this.firewall.RemoveSourceIP(ip)
	if err != nil {
		this.write(writer, http.StatusInternalServerError, nil)
		return
	}

	this.write(writer, http.StatusOK, nil)
}

func (this *Server) write(writer http.ResponseWriter, code int, data any) {
	writer.WriteHeader(code)
	if data != nil {
		respJSON, err := json.Marshal(data)
		if err != nil {
			return
		}
		_, _ = writer.Write(respJSON)
	}
}
