//go:build !android
// +build !android

package detect_ip

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func init() {
	newEnv := os.Getenv("no_proxy") + ",api.ipify.org"
	os.Setenv("no_proxy", newEnv)
	os.Setenv("NO_PROXY", newEnv)
}

func httpGetPublicIp(url string) (ret net.IP, err error) {
	if len(url) == 0 {
		return
	}
	ret, err = get(url)
	if err != nil {
		return
	}

	if !tools.IsIPv4(ret) {
		return nil, errors.New("not valid ipv4: " + ret.String())
	}

	return
}

func resolveAs(host string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if tools.IsIPv4(ip) {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("failed resolve ipv4 of host %v", host)
}

var zeroDialer net.Dialer

func get(api string) (net.IP, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}
	client.Transport = transport

	req, _ := http.NewRequest("GET", api, nil)

	tools.Debug("%s %s %s\n", req.Proto, req.Method, req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	retBytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	ret := string(retBytes)
	ret = strings.TrimSpace(ret)

	return net.ParseIP(ret), nil
}
