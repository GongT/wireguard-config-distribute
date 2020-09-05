package tools

import "fmt"

var build_date string
var build_git_hash string

func ShowVersion() {
	Error("Build Date: %v, Git Hash: %v", build_date, build_git_hash)
}

func GetVersion() string {
	return build_git_hash
}

func GetBuildDate() string {
	return build_date
}

func GetAgent() string {
	return fmt.Sprintf("gongt-wireguard/%v@%v", build_git_hash, build_date)
}
