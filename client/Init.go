package client

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/czxichen/Goprograme/parse"
)

type clientConfig struct {
	CheckMd5    bool   `json:"checkmd5"`
	Home        string `json:"home"`
	RequestMode string `json:"requestmode"`
	MasteUrl    string `json:"masteurl"`
	Primary     string `json:"primary"`
	Action      string `json:"action"`
}

const tmp = "tmp/"

var (
	cfgpath string
	CfgName string
	Config  clientConfig
)

var File *os.File

func init() {
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		println("版本:2016-08-15")
		os.Exit(1)
	}
	flag.StringVar(&cfgpath, "f", "", "指定配置文件 -f cfg.json")
	flag.StringVar(&CfgName, "n", "", "只更新指定服务端配置文件 -n system-cofnig.xml 结合-a getcfg使用")
	flag.StringVar(&Config.Action, "a", "", "指定要做的操作 -a install、getcfg")
	flag.StringVar(&Config.MasteUrl, "m", "127.0.0.1:1789", "指定服务端IP端口 -m 127.0.0.1:1789")
	flag.StringVar(&Config.RequestMode, "r", "http", "指定请求的模式 -r https")
	flag.StringVar(&Config.Home, "p", "/test", "指定解压的根目录 -p /test")
	flag.BoolVar(&Config.CheckMd5, "M", true, "是否检查文件md5 -M true")
	flag.StringVar(&Config.Primary, "P", "", "指定请求的关键字 -P 7400001")
	flag.Parse()
	var err error
	File, err = os.Create("client.log")
	if err != nil {
		log.Println("创建日志文件失败:", err)
		os.Exit(1)
	}
	log.SetOutput(File)
	os.MkdirAll(tmp, 0644)
	parseconfig()
	if Config.Action == "" {
		println("必须指定-a参数")
		os.Exit(1)
	}
	Config.Home = strings.Replace(Config.Home, "\\", "/", -1)
	if !strings.HasSuffix(Config.Home, "/") {
		Config.Home += "/"
	}
}

func parseconfig() {
	if cfgpath == "" {
		return
	}
	buf, err := parse.Parse(cfgpath)
	if err != nil {
		log.Println("读取本地配置文件失败:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(buf, &Config)
	if err != nil {
		log.Println("解析本地配置文件失败:", err)
		os.Exit(1)
	}
}
