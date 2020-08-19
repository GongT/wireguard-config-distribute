package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	os.MkdirAll(filepath.Dir(dst), os.FileMode(0755))
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error read file [%v]: %v", src, err)
	}

	err = ioutil.WriteFile(dst, input, 0755)
	if err != nil {
		return fmt.Errorf("error copy file [%v]: %v", dst, err)
	}

	return nil
}
