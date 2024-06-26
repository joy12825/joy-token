// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/joy12825/gf.

package ghttp

import (
	"context"
	"fmt"

	"github.com/joy12825/gf/v2/net/gipv4"
	"github.com/joy12825/gf/v2/net/gsvc"
	"github.com/joy12825/gf/v2/os/gctx"
	"github.com/joy12825/gf/v2/text/gstr"
	"github.com/joy12825/gf/v2/util/gconv"
)

// doServiceRegister registers current service to Registry.
func (s *Server) doServiceRegister() {
	if s.registrar == nil {
		return
	}
	s.serviceMu.Lock()
	defer s.serviceMu.Unlock()
	var (
		ctx      = gctx.GetInitCtx()
		protocol = gsvc.DefaultProtocol
		insecure = true
		err      error
	)
	if s.config.TLSConfig != nil {
		protocol = `https`
		insecure = false
	}
	metadata := gsvc.Metadata{
		gsvc.MDProtocol: protocol,
		gsvc.MDInsecure: insecure,
	}
	s.service = &gsvc.LocalService{
		Name:      s.GetName(),
		Endpoints: s.calculateListenedEndpoints(ctx),
		Metadata:  metadata,
	}
	s.Logger().Debugf(ctx, `service register: %+v`, s.service)
	if len(s.service.GetEndpoints()) == 0 {
		s.Logger().Warningf(ctx, `no endpoints found to register service, abort service registering`)
		return
	}
	if s.service, err = s.registrar.Register(ctx, s.service); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
}

// doServiceDeregister de-registers current service from Registry.
func (s *Server) doServiceDeregister() {
	if s.registrar == nil {
		return
	}
	s.serviceMu.Lock()
	defer s.serviceMu.Unlock()
	if s.service == nil {
		return
	}
	var ctx = gctx.GetInitCtx()
	s.Logger().Debugf(ctx, `service deregister: %+v`, s.service)
	if err := s.registrar.Deregister(ctx, s.service); err != nil {
		s.Logger().Errorf(ctx, `%+v`, err)
	}
	s.service = nil
}

func (s *Server) calculateListenedEndpoints(ctx context.Context) gsvc.Endpoints {
	var (
		configAddr = s.config.Address
		endpoints  = make(gsvc.Endpoints, 0)
		addresses  = s.config.Endpoints
	)
	if configAddr == "" {
		configAddr = s.config.HTTPSAddr
	}
	if len(addresses) == 0 {
		addresses = gstr.SplitAndTrim(configAddr, ",")
	}
	for _, address := range addresses {
		var (
			addrArray     = gstr.Split(address, ":")
			listenedIps   []string
			listenedPorts []int
		)
		if len(addrArray) == 1 {
			addrArray = append(addrArray, gconv.String(defaultEndpointPort))
		}
		// IPs.
		switch addrArray[0] {
		case "127.0.0.1":
			// Nothing to do.
		case "0.0.0.0", "":
			intranetIps, err := gipv4.GetIntranetIpArray()
			if err != nil {
				s.Logger().Errorf(ctx, `error retrieving intranet ip: %+v`, err)
				return nil
			}
			// If no intranet ips found, it uses all ips that can be retrieved,
			// it may include internet ip.
			if len(intranetIps) == 0 {
				allIps, err := gipv4.GetIpArray()
				if err != nil {
					s.Logger().Errorf(ctx, `error retrieving ip from current node: %+v`, err)
					return nil
				}
				s.Logger().Noticef(
					ctx,
					`no intranet ip found, using internet ip to register service: %v`,
					allIps,
				)
				listenedIps = allIps
				break
			}
			listenedIps = intranetIps
		default:
			listenedIps = []string{addrArray[0]}
		}
		// Ports.
		switch addrArray[1] {
		case "0":
			listenedPorts = s.GetListenedPorts()
		default:
			listenedPorts = []int{gconv.Int(addrArray[1])}
		}
		for _, ip := range listenedIps {
			for _, port := range listenedPorts {
				endpoints = append(endpoints, gsvc.NewEndpoint(fmt.Sprintf(`%s:%d`, ip, port)))
			}
		}
	}
	return endpoints
}
