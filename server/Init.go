package server

import (
	"log"
)

var (
	Cfg          Config
	serverpkg    string
	Templatedir  string
	Variables    map[string][]string
	Pathrelation map[string][]string
)

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
