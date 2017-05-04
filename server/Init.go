package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	Cfg          Config
	cpath        string
	serverpkg    string
	Templatedir  string
	Variables    map[string][]string
	Pathrelation map[string][]string
	FlagSet      = flag.FlagSet{Usage: func() {}}
)

func init() {
	FlagSet.StringVar(&cpath, "c", "", "-c cfg.json 从文件读取配置")
	FlagSet.StringVar(&Cfg.IP, "l", ":1789", "-l 127.0.0.1;1789 指定监听的IP端口")
	FlagSet.StringVar(&Cfg.Proto, "p", "http", "-p http 指定协议类型,http,https")
	FlagSet.StringVar(&Cfg.CrtPath, "crt", "", "-crt tool.crt 指定https的crt文件")
	FlagSet.StringVar(&Cfg.Keypath, "key", "", "-key 指定https的key文件")
	FlagSet.StringVar(&Cfg.Logname, "log", "run.log", "-log run.log 指定log名称路径")
	FlagSet.StringVar(&Cfg.Download, "d", "download", "-d ./download 指定下载文件所在的目录")
}

func readconfig(cfgpath string, cfg *Config) {
	File, err := os.Open(cfgpath)
	if err != nil {
		log.Fatalln(err)
	}
	defer File.Close()
	var buf []byte
	b := bufio.NewReader(File)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalln(err)
		}
		line = bytes.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}
		index := bytes.Index(line, []byte("#"))
		if index == 0 {
			continue
		}
		if index > 0 {
			line = line[:index]
		}
		buf = append(buf, line...)
	}
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		log.Println(string(buf))
		log.Fatalln("解析配置文件失败:", err)
	}
}

type serverInfo struct {
	Variable map[string]string
	Relation map[string][]string
}

type Config struct {
	IP       string `json:"ip"`
	Proto    string `json:"proto"`
	CrtPath  string `json:"crtpath"`
	Keypath  string `json:"keypath"`
	Logname  string `json:"logname"`
	Download string `json:"download"`
}

func Parseconfig() {
	log.Println("开始读取变量和模版文件")
	var err error
	Pathrelation, Variables, err = relationConfig(Templatedir + "server.xlsx")
	if err != nil {
		log.Fatalln("读取关系表失败:", err)
	}

	log.Printf("server.xlsx解析成功\n")
	log.Printf("路径关系:\n")
	for key, value := range Pathrelation {
		log.Printf("文件名:%s\t路径:%v\n", key, value)
	}
	log.Printf("变量:\n")
	for key, value := range Variables {
		log.Printf("索引:%s\t值:%v\n", key, value)
	}
	var list []string
	for k, _ := range Pathrelation {
		list = append(list, k)
	}
	log.Println("开始打包模版文件")
	err = Zip(Cfg.Download, list)
	if err != nil {
		log.Fatalln("打包模版文件失败:", err)
	}
	log.Println("开始获取最新的服务端包")
	err = getsrvpkg()
	if err != nil {
		log.Fatalln("检查server包出错:", err)
	}
	log.Printf("服务端包名称:%s\n", serverpkg)
}
