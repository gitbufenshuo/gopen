package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gitbufenshuo/gopen/help"
)

func main() {
	//
	for idx := 1; idx != len(os.Args); idx++ {
		convertOneFile(os.Args[idx])
	}
}

func convertOneFile(filename string) {
	fmt.Printf("转换%s为png格式\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("读取文件%s失败%v\n", filename, err)
		return
	}
	outdata, err := help.ToPng(data)
	if err != nil {
		fmt.Printf("转换文件[%s]失败[%v]\n", filename, err)
		return
	}
	segs := strings.Split(filename, ".")
	if segs[len(segs)-1] != "jpg" {
		fmt.Printf("后缀名[%s]错误\n", filename)
		return
	}
	segs[len(segs)-1] = "png"
	ioutil.WriteFile(strings.Join(segs, "."), outdata, 0644)
}
