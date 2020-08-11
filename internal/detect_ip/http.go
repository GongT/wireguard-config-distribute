package detect_ip

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func init() {
	newEnv := os.Getenv("no_proxy") + ",api.ipify.org,api6.ipify.org"
	os.Setenv("no_proxy", newEnv)
	os.Setenv("NO_PROXY", newEnv)
}

func httpGetPublicIp4(url string) (ret string, err error) {
	if len(url) == 0 {
		url = "https://api.ipify.org/"
	}
	ret, err = get(url)
	if err != nil {
		return
	}

	if !IsValidIPv4(ret) {
		return "", errors.New("Not valid ipv4: " + ret)
	}

	return
}

func httpGetPublicIp6(url string) (ret string, err error) {
	if len(url) == 0 {
		url = "https://api.ipify.org/"
	}
	for i := 0; i < 3; i++ {
		ret, err = get(url)
		if err == nil {
			break
		}
		tools.Error("failed get ip: %v", err)
	}
	if err != nil {
		return
	}

	if !IsValidIPv6(ret) {
		return "", errors.New("Not valid ipv6: " + ret)
	}

	return
}

func get(url string) (ret string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}

	retBytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}

	ret = string(retBytes)
	ret = strings.TrimSpace(ret)

	return
}
