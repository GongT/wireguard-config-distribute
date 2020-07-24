package detect_ip

import (
	"io/ioutil"
	"net/http"
	"os"
)

var envModify = false

func addEnv() {
	if envModify {
		return
	}
	envModify = true

	newEnv := os.Getenv("no_proxy") + ",api.ipify.org,api6.ipify.org"
	os.Setenv("no_proxy", newEnv)
	os.Setenv("NO_PROXY", newEnv)
}

func GetPublicIp() (string, error) {
	addEnv()
	return get("http://api.ipify.org/")
}

func GetPublicIp6() (ret string, err error) {
	addEnv()
	ret, err = get("http://api6.ipify.org/")
	if err != nil {
		return
	}

	if !IsValidIPv6(ret) {
		return "", nil
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

	return
}
