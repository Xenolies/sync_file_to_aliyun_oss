package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
	LocalDir        string
	HugoSiteDir     string
}

var LocalMd5List = make(map[string]string)

//var LocalMd5Cache = "./Local_Md5_CACHE.json"

var OSSMd5List = make([]string, 0)

var configPath string = "./config.json"

func main() {

	dir, _ := os.ReadDir(configPath)
	if dir == nil {
		//MakeConfigFile(configDEMO)
		fmt.Println("..............创建配置文件........")
		fmt.Println("未检测到配置文件,生成config.json文件,请配置: ")
		new_userconfig := Config{
			AccessKeyId:     "",
			AccessKeySecret: "",
			Endpoint:        "",
			BucketName:      "",
			LocalDir:        "",
			HugoSiteDir:     "",
		}

		fmt.Println("请输入AccessKeyId(输入回车结束): ")
		reader := bufio.NewReader(os.Stdin) //标准输入流
		AccessKeyId, _ := reader.ReadString('\n')
		new_userconfig.AccessKeyId = strings.Trim(AccessKeyId, "\r\n")
		fmt.Printf("AccessKeyId: %v\n", new_userconfig.AccessKeyId)
		fmt.Println("........................................................")

		fmt.Println("请输入AccessKeySecret(输入回车结束): ")
		AccessKeySecret, _ := reader.ReadString('\n')
		new_userconfig.AccessKeySecret = strings.Trim(AccessKeySecret, "\r\n")
		fmt.Printf("AccessKeySecret: %v\n", new_userconfig.AccessKeySecret)
		fmt.Println("........................................................")

		fmt.Println("请输入Endpoint (输入回车结束 ,Endpoint格式类似这样 oss-cn-hongkong.aliyuncs.com): ")
		Endpoint, _ := reader.ReadString('\n')
		new_userconfig.Endpoint = strings.Trim(Endpoint, "\r\n")
		fmt.Printf("Endpoint: %v\n", new_userconfig.Endpoint)
		fmt.Println("........................................................")

		fmt.Println("请输入BucketName: ")
		BucketName, _ := reader.ReadString('\n')
		new_userconfig.BucketName = strings.Trim(BucketName, "\r\n")
		fmt.Printf("BucketName: %v\n", new_userconfig.BucketName)
		fmt.Println("........................................................")

		fmt.Println("请输入LocalDir (LocalDir是你本地博客目录下的public目录), ")
		fmt.Println("格式类似这样 D:/OdinXu/odinxu-blog/public/")
		fmt.Println("请输入(输入回车结束): ")
		LocalDir, _ := reader.ReadString('\n')
		new_userconfig.LocalDir = strings.Trim(LocalDir, "\r\n")
		fmt.Printf("LocalDir: %v\n", new_userconfig.LocalDir)
		fmt.Println("........................................................")

		fmt.Println("请输入HugoSiteDir (HugoSiteDir是你本地博客所在目录), ")
		fmt.Println("格式类似这样 D:/OdinXu/odinxu-blog/")
		fmt.Println("请输入(输入回车结束): ")
		HugoSiteDir, _ := reader.ReadString('\n')
		new_userconfig.HugoSiteDir = strings.Trim(HugoSiteDir, "\r\n")
		fmt.Printf("HugoSiteDir: %v\n", new_userconfig.HugoSiteDir)
		fmt.Println("........................................................")

		MakeConfigFile(new_userconfig)
		fmt.Println("配置文件创建成功!五秒后将关闭程序")
		// getLocalDirList(new_userconfig.LocalDir, LocalMd5List)
		time.Sleep(5 * time.Second)

	} else {

		userConfig := LoadConfig()
		fmt.Println("当前OSS配置为:")
		fmt.Printf("AccessKeyId = %v\n", userConfig.AccessKeyId)
		fmt.Printf("AccessKeySecret = %v\n", userConfig.AccessKeySecret)
		fmt.Printf("Endpoint = %v\n", userConfig.Endpoint)
		fmt.Printf("BucketName = %v\n", userConfig.BucketName)
		fmt.Printf("LocalDir = %v\n", userConfig.LocalDir)
		fmt.Printf("HugoSiteDir = %v\n", userConfig.HugoSiteDir)
		fmt.Println("........................................................")

		client := getBucketClient(userConfig)

		//生成新 Hugo 文件

		//hugoCommand(userConfig)

		//go getLocalDirList(userConfig.LocalDir, LocalMd5List)
		OSSMd5List = getOSSMd5List(client, userConfig, OSSMd5List)
		//
		//获取 OSS 和 本地文件列表

		//
		//ossDelete(client, userConfig, "robots.txt")

		//getOSSMd5List(client, userConfig, OSSMd5List)
		//
		//fmt.Println(OSSMd5List)
		//time.Sleep(10 * time.Second)

		//ossDelete(client, userConfig, "about/index.html")

		//bucket, _ := client.Bucket(userConfig.BucketName)

		//exist, _ := bucket.IsObjectExist("about/index.html")
		//fmt.Println(exist)

		//isAnagram(client, userConfig, OSSMd5List, LocalMd5List)
		//pause()

	}

}

func getBucketClient(bucketConfig Config) *oss.Client {
	client, err := oss.New(bucketConfig.Endpoint, bucketConfig.AccessKeyId, bucketConfig.AccessKeySecret)
	if err != nil {
		fmt.Printf(" oss.New err: %v", err)
	}

	return client
}

func getOSSMd5List(ossClient *oss.Client, config Config, ossList []string) []string {
	bucket, err := ossClient.Bucket(config.BucketName)

	if err != nil {
		fmt.Printf(" ossClient.Bucke err: %v", err)
	}

	lsRes, err := bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		fmt.Printf(" bucket.ListObjects err: %v", err)
	}

	for _, object := range lsRes.Objects {
		meta, _ := bucket.GetObjectMeta(object.Key)

		//OSSList[strings.Trim(meta.Get("Etag"), "\"")] = object.Key

		fmt.Println(meta)
		//
		//fmt.Println(meta.Get("Etag"))

		ossList = append(ossList, strings.Trim(meta.Get("Etag"), "\""))

	}

	return ossList
}

func getLocalDirList(path string, files map[string]string) {
	dir, _ := os.ReadDir(path)

	for _, fi := range dir {
		if fi.IsDir() {
			fullDir := path + fi.Name() + "/"
			getLocalDirList(fullDir, files)
		} else {
			md5num, _ := getMd5(path + fi.Name())

			files[md5num] = path + fi.Name()

			//jsonFiles, _ := json.Marshal(files)
			//
			//os.WriteFile(LocalMd5Cache, jsonFiles, 0644)
			//
			//return jsonFiles

		}

	}

}

// 比较两个MD5 list(map)的差异 并且上传文件
func isAnagram(ossClient *oss.Client, config Config, OSS []string, Local map[string]string) {

	//if len(OSS) < len(Local) {
	//	fmt.Println("len(OSS) < len(Local) Begin OSS: ", len(OSS))
	//	fmt.Println("len(OSS) < len(Local) Begin Local: ", len(Local))
	//	//OSS 文件数小于 本地 本地上传到OSS
	//	for ossItem, _ := range OSS {
	//		_, ok := Local[OSS[ossItem]]
	//		if ok {
	//			// OSS 中存在相同MD5值的文件 ,跳过
	//			delete(Local, OSS[ossItem])
	//		}
	//	}
	//	fmt.Println("len(OSS) < len(Local) OSS: ", len(OSS))
	//	fmt.Println("len(OSS) < len(Local) Local: ", len(Local))
	//
	//	//遍历上传
	//	//fmt.Println("上传文件中....")
	//	//for _, value := range Local {
	//	//	//fmt.Println("strings.Trim(value, config.LocalDir): ", strings.TrimLeft(value, config.LocalDir))
	//	//	//fmt.Println("value:", value)
	//	//	//
	//	//	index := strings.LastIndex(config.LocalDir, "/")
	//	//	//fmt.Println(index)
	//	//	//str := value[index:]
	//	//	//fmt.Println("value[index:]: ", str)
	//	//
	//	//	//ossUpload(ossClient, config, strings.TrimLeft(value, config.LocalDir), value)
	//	//	ossUpload(ossClient, config, value[index+1:], value)
	//	//}
	//
	//}
	fmt.Println("Begin OSS: ", len(OSS))
	fmt.Println("Begin Local: ", len(Local))
	// OSS 文件数等于本地 查询MD5 删除MD5不同的文件

	for ossItem, _ := range OSS {
		_, ok := Local[OSS[ossItem]]
		//fmt.Println(ok)
		if ok {
			// OSS 中存在相同MD5值的文件 ,跳过
			delete(Local, OSS[ossItem])
		}
	}

	if len(Local) > 0 {

		fmt.Println("len(OSS) == len(Local) OSS: ", len(OSS))
		fmt.Println("len(OSS) == len(Local) Local: ", len(Local))

		//遍历上传
		//fmt.Println("上传文件中....")

		for _, value := range Local {
			//fmt.Println("strings.Trim(value, config.LocalDir): ", strings.TrimLeft(value, config.LocalDir))
			//fmt.Println("value:", value)
			//
			index := strings.LastIndex(config.LocalDir, "/")
			//fmt.Println(index)
			//str := value[index:]
			//fmt.Println("value[index:]: ", str)

			//ossDelete(ossClient, config, value[index+1:])
			ossUpload(ossClient, config, value[index+1:], value)
		}

	} else {
		fmt.Println("OSS 中文件和本地文件相同无需上传")
	}

	//} else {
	//	// OSS 文件数大于本地 以本地为准
	//
	//	for ossItem, _ := range OSS {
	//		_, ok := Local[OSS[ossItem]]
	//		//fmt.Println(ok)
	//		if ok {
	//			// OSS 中存在相同MD5值的文件 ,跳过
	//			delete(Local, OSS[ossItem])
	//		}
	//	}
	//
	//}

	//for key := range Local {
	//
	//}
}

// 计算MD5
func getMd5(filepath string) (string, error) {
	f, _ := os.Open(filepath)
	defer f.Close()

	body, _ := io.ReadAll(f)

	md5sum := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()

	return strings.ToUpper(md5sum), nil
}

// OSS 删除
func ossDelete(client *oss.Client, config Config, object string) {
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		fmt.Printf("client.Bucket ERROR: %v\n", err)
	}

	err = bucket.DeleteObject(object)
	if err != nil {
		fmt.Printf("bucket.DeleteObject ERROR: %v\n", err)
	}
	fmt.Printf("文件 %v 已从 %v 中删除", object, config.BucketName)

}

// OSS上传
func ossUpload(client *oss.Client, config Config, fileName string, filePath string) {
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		fmt.Printf("client.Bucket ERROR: %v\n", err)
	}

	err = bucket.PutObjectFromFile(fileName, filePath)
	fmt.Printf("正在上传 %v , 本地路径: %v\n", fileName, filePath)
	if err != nil {
		fmt.Printf("bucket.PutObjectFromFile ERROR: %v\n", err)
	}

	fmt.Printf("文件 %v 上传完毕\n", fileName)

}

// SaveConfig 保存配置文件
func SaveConfig(config *Config) {

	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configPath, data, 0660)
	if err != nil {
		log.Fatal(err)
	}

}

// LoadConfig 读取配置文件
func LoadConfig() (config Config) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config = Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// MakeConfigFile 生成配置文件
func MakeConfigFile(config Config) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configPath, data, 0660)
	if err != nil {
		log.Fatal(err)
	}
}

// hugo 命令
func hugoCommand(config Config) {

	//删除Public
	err := os.RemoveAll(config.LocalDir)
	if err != nil {
		fmt.Printf("os.RemoveAll ERROR: %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("已经删除public目录,准备生成新文件.")

	command := exec.Command("hugo")
	command.Dir = config.HugoSiteDir
	command.Stdout = &bytes.Buffer{}

	err = command.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(command.Stderr.(*bytes.Buffer).String())
		return
	}

	fmt.Println(command.Stdout.(*bytes.Buffer).String())
	fmt.Println("新文件生成完毕!!")

}

func pause() {

	fmt.Println("--------------------------------------")

	fmt.Printf("按任意键退出...")

	b := make([]byte, 1)

	os.Stdin.Read(b)

}
