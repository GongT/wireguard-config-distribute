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
	newEnv := os.Getenv("no_proxy") + ",api.ipify.org,api6.ipify.org"
	os.Setenv("no_proxy", newEnv)
	os.Setenv("NO_PROXY", newEnv)
}

func httpGetPublicIp4(url string) (ret net.IP, err error) {
	if len(url) == 0 {
		return
	}
	ret, err = get(url, true)
	if err != nil {
		return
	}

	if !tools.IsIPv4(ret) {
		return nil, errors.New("not valid ipv4: " + ret.String())
	}

	return
}

func httpGetPublicIp6(url string) (ret net.IP, err error) {
	if len(url) == 0 {
		return
	}
	ret, err = get(url, false)
	if err != nil {
		return
	}

	if !tools.IsIPv6(ret) {
		return nil, errors.New("Not valid ipv6: " + ret.String())
	}

	return
}

func resolveAs(host string, ipv4 bool) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if tools.IsIPv4(ip) {
			if ipv4 {
				return ip.String(), nil
			}
		} else {
			if !ipv4 {
				return "[" + ip.String() + "]", nil
			}
		}
	}

	v := 6
	if ipv4 {
		v = 4
	}
	return "", fmt.Errorf("failed resolve ipv%v of host %v", v, host)
}

var zeroDialer net.Dialer

func get(api string, ipv4 bool) (net.IP, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var dial_type string
		if ipv4 {
			dial_type = "tcp4"
		} else {
			dial_type = "tcp6"
		}
		return zeroDialer.DialContext(ctx, dial_type, addr)
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
