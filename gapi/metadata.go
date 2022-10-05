package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	hostname = "grpcgateway-user-agent"
	useragent = "user-agent"
	clientip = "x-forwarded-for"
)

type Metadata struct {
	ClientIP string
	HostName string
}

func (s *server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if resp := md.Get(hostname); len(resp) > 0 {
			mtdt.HostName = resp[0]
		}
		if resp := md.Get(useragent); len(resp) > 0 {
			mtdt.HostName = resp[0]
		}
		if resp := md.Get(clientip); len(resp) > 0 {
			mtdt.ClientIP = resp[0]
		}
	}
	if p,ok := peer.FromContext(ctx); ok{
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
