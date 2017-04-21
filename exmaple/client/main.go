package main

import "github.com/czxichen/configmanage/client"

func main() {
	defer client.File.Close()
	url := client.Config.RequestMode + "://" + client.Config.MasteUrl

	switch client.Config.Action {
	case "install":
		path := url + "/serverpackage"
		client.InitPackage(path)
		client.CreateConfig(client.CfgName, url, client.Config.Primary)
	case "getcfg":
		client.CreateConfig(client.CfgName, url, client.Config.Primary)
	default:
		println("-a 参数无效,-h 查看帮助命令.")
	}
}
