package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Config.Indent = "    "

	filePath := os.Getenv("GOFILE")
	pkgName := os.Getenv("GOPACKAGE")

	oData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	if loc := regexp.MustCompile("(?m)^package\\s").FindIndex(oData); loc == nil {
		panic(errors.New("Can not find `package ???` from file"))
	} else {
		oData = oData[:loc[0]]
	}

	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("got argument: (", len(os.Args), ")", os.Args)
		panic(errors.New("Requires 1 or 2 arguments, Usage: $0 <TYPE> <import from>"))
	}
	cType := os.Args[1]
	dType := cType

	imports := ""
	if len(os.Args) == 3 {
		dType = "t." + cType
		imports = "t \"" + os.Args[2] + "\""
	}

	cType = strings.ToUpper(cType[:1]) + cType[1:]

	data := template
	data = strings.ReplaceAll(data, "THIS_PACKAGE", pkgName)
	data = strings.ReplaceAll(data, "X_IMPORT", imports)
	data = strings.ReplaceAll(data, "tttt", dType)
	data = strings.ReplaceAll(data, "TTTT", cType)

	err = ioutil.WriteFile(filePath, append(oData, []byte(data)...), os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}

var template = `package THIS_PACKAGE

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/net/context"
	X_IMPORT
)

type AsyncChanTTTT struct {
	ch     chan tttt
	closed bool
}

func NewChanTTTT() *AsyncChanTTTT {
	ch := make(chan tttt)
	return &AsyncChanTTTT{
		ch:     ch,
		closed: false,
	}
}
func (cc *AsyncChanTTTT) Close() {
	cc.closed = true
	close(cc.ch)
}

func (cc *AsyncChanTTTT) Read() <-chan tttt {
	if cc.closed {
		tools.Error("Read channel after closed!")
	}
	return cc.ch
}
func (cc *AsyncChanTTTT) Write(data tttt) {
	if cc.closed {
		tools.Error("Write channel after closed!")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			tools.Error("Data not consumed within 1s! (data is %v)", data)
		case cc.ch <- data:
		}
	}()
}
`
