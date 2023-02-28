package main

import (
	"GolangAliOSSUpload/main/BucketOperation"
	"GolangAliOSSUpload/main/Command"
	"GolangAliOSSUpload/main/Config"
	"GolangAliOSSUpload/main/LocalOperation"
	"fmt"
)

var LocalMd5List = make(map[string]string)

//var LocalMd5Cache = "./Local_Md5_CACHE.json"

var OSSMd5List = make(map[string]string)

var configPath string = "./config.json"

func main() {
	if Config.ConfigMaker() {
		userConfig := Config.LoadConfig()
		fmt.Println("当前OSS配置为:")
		fmt.Printf("AccessKeyId = %v\n", userConfig.AccessKeyId)
		fmt.Printf("AccessKeySecret = %v\n", userConfig.AccessKeySecret)
		fmt.Printf("Endpoint = %v\n", userConfig.Endpoint)
		fmt.Printf("BucketName = %v\n", userConfig.BucketName)
		fmt.Printf("LocalDir = %v\n", userConfig.LocalDir)
		fmt.Printf("HugoSiteDir = %v\n", userConfig.HugoSiteDir)
		fmt.Println("........................................................")

		client := BucketOperation.GetBucketClient(userConfig)

		//生成新 Hugo 文件
		Command.Command_hugo(userConfig)

		OSSMd5List = BucketOperation.GetOSSMd5List(client, userConfig, OSSMd5List)
		LocalOperation.GetLocalDirList(userConfig.LocalDir, userConfig, LocalMd5List)

		LocalOperation.Comparison(client, userConfig, OSSMd5List, LocalMd5List)
		Command.Pause()
	} else {
		Config.ConfigMaker()
	}

}
