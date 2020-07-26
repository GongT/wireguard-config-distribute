package detect_ip

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func init() {
	newEnv := os.Getenv("no_proxy") + ",api.ipify.org,api6.ipify.org"
	os.Setenv("no_proxy", newEnv)
	os.Setenv("NO_PROXY", newEnv)
}

func httpGetPublicIp4() (ret string, err error) {
	ret, err = get("https://api.ipify.org/")
	if err != nil {
		return
	}

	if !IsValidIPv4(ret) {
		return "", errors.New("Not valid ipv4: " + ret)
	}

	return
}

func httpGetPublicIp6() (ret string, err error) {
	ret, err = get("https://api6.ipify.org/")
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
