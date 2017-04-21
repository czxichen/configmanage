package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/czxichen/configmanage/server"
)

var (
	version string = "2017-02-21"
)

func init() {
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		fmt.Printf("Version:%s\n" + version)
		os.Exit(0)
	}

	cpath := flag.String("c", "", "-c cfg.json 从文件读取配置")
	flag.StringVar(&server.Cfg.IP, "l", ":1789", "-l 127.0.0.1;1789 指定监听的IP端口")
	flag.StringVar(&server.Cfg.Proto, "p", "http", "-p http 指定协议类型")
	flag.StringVar(&server.Cfg.CrtPath, "ctr", "", "-crt tool.crt 指定https的crt文件")
	flag.StringVar(&server.Cfg.Keypath, "key", "", "-key 指定https的key文件")
	flag.StringVar(&server.Cfg.Logname, "log", "run.log", "-log run.log 指定log名称路径")
	flag.StringVar(&server.Cfg.Download, "d", "download", "-d ./download 指定下载文件所在的目录")
	flag.Parse()

	if *cpath != "" {
		readconfig(*cpath, &server.Cfg)
	}

	server.Cfg.Download = strings.Replace(server.Cfg.Download, "\\", "/", -1)
	if !strings.HasSuffix(server.Cfg.Download, "/") {
		server.Cfg.Download += "/"
	}

	File, err := os.Create(server.Cfg.Logname)
	if err != nil {
		log.Fatalln("初始化日志文件失败:", err)
	}
	log.SetOutput(File)

	server.Templatedir = server.Cfg.Download + "template/"
	server.Parseconfig()

	go server.Notify(server.Templatedir, 10, func() {
		server.Parseconfig()
	})
}

func readconfig(cfgpath string, cfg *server.Config) {
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

func main() {
	err := server.Server()
	if err != nil {
		fmt.Println(err)
	}
}
