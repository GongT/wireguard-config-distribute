package tools

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func GetRemoteFromContext(ctx context.Context) string {
	p, ok1 := peer.FromContext(ctx)
	if !ok1 {
		return ""
	}

	md, ok2 := metadata.FromIncomingContext(ctx)
	if !ok2 {
		return p.Addr.String()
	}

	var ip string

	ip = findValid(md.Get("x-real-ip"))
	if len(ip) > 0 {
		return ip
	}

	ip = findValid2(md.Get("x-forwarded-for"))
	if len(ip) > 0 {
		return ip
	}

	return p.Addr.String()
}

func findValid(list []string) string {
	for _, s := range list {
		if ip := net.ParseIP(s); ip != nil {
			return ip.String()
		}
	}
	return ""
}

func findValid2(list []string) string {
	for _, s := range list {
		sa := strings.Split(s, ",")
		if len(sa) == 0 {
			continue
		}
		if ip := net.ParseIP(sa[0]); ip != nil {
			return ip.String()
		}
	}
	return ""
}
