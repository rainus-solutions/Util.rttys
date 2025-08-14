/* SPDX-License-Identifier: MIT */
/*
 * Author: Jianhui Zhao <zhaojh329@gmail.com>
 */

package main

import (
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

type RttyServer struct {
	mu            sync.RWMutex
	groups        sync.Map
	cfg           Config
	httpProxyPort int
}

type DeviceGroup struct {
	devices sync.Map
	count   atomic.Int32
}

func (srv *RttyServer) Run() error {
	log.Debug().Msgf("%+v", srv.cfg)

	if srv.cfg.PprofAddr != "" {
		go srv.ListenPprof()
	}

	go srv.ListenDevices()
	go srv.ListenHttpProxy()

	return srv.ListenAPI()
}

func (srv *RttyServer) ListenPprof() {
	ln, err := net.Listen("tcp", srv.cfg.PprofAddr)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to start pprof server")
		return
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	log.Info().Msgf("Starting pprof server on: %s", addr)

	host := addr.IP.String()
	if host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}
	log.Info().Msgf("Access pprof at: http://%s:%d/debug/pprof/", host, addr.Port)

	err = http.Serve(ln, nil)
	if err != nil {
		log.Error().Err(err).Msgf("pprof server failed")
	}
}

func (srv *RttyServer) GetDevice(group, id string) *Device {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	g := srv.GetGroup(group, false)
	if g == nil {
		return nil
	}

	if v, ok := g.devices.Load(id); ok {
		return v.(*Device)
	}

	return nil
}

func (srv *RttyServer) AddDevice(dev *Device) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	g := srv.GetGroup(dev.group, true)

	if _, loaded := g.devices.LoadOrStore(dev.id, dev); loaded {
		return false
	}

	g.count.Add(1)

	return true
}

func (srv *RttyServer) DelDevice(dev *Device) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	g := srv.GetGroup(dev.group, false)
	if g == nil {
		return
	}

	if deleted := g.devices.CompareAndDelete(dev.id, dev); deleted {
		if g.count.Add(-1) == 0 {
			srv.groups.Delete(dev.group)
		}
	}
}

func (srv *RttyServer) GetGroup(group string, create bool) *DeviceGroup {
	if create {
		val, _ := srv.groups.LoadOrStore(group, &DeviceGroup{})
		return val.(*DeviceGroup)
	} else {
		val, ok := srv.groups.Load(group)
		if !ok {
			return nil
		}
		return val.(*DeviceGroup)
	}
}
