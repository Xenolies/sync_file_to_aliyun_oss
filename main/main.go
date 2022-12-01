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

var OSSMd5List = make(map[string]string)

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
		hugoCommand(userConfig)

		OSSMd5List = getOSSMd5List(client, userConfig, OSSMd5List)
		getLocalDirList(userConfig.LocalDir, userConfig, LocalMd5List)

		isAnagram(client, userConfig, OSSMd5List, LocalMd5List)
		pause()

	}

}

func getBucketClient(bucketConfig Config) *oss.Client {
	client, err := oss.New(bucketConfig.Endpoint, bucketConfig.AccessKeyId, bucketConfig.AccessKeySecret)
	if err != nil {
		fmt.Printf(" oss.New err: %v", err)
	}

	return client
}

func getOSSMd5List(ossClient *oss.Client, config Config, ossList map[string]string) map[string]string {
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

		objMD5 := strings.Trim(meta.Get("Etag"), "\"")

		ossList[object.Key] = objMD5

	}

	return ossList
}

func getLocalDirList(path string, config Config, files map[string]string) {

	dir, _ := os.ReadDir(path)

	for _, fi := range dir {
		if fi.IsDir() {
			fullDir := path + fi.Name() + "/"
			getLocalDirList(fullDir, config, files)
		} else {

			md5num, _ := getMd5(path + fi.Name())

			files[path+fi.Name()] = md5num

			//jsonFiles, _ := json.Marshal(files)
			//
			//os.WriteFile(LocalMd5Cache, jsonFiles, 0644)
			//
			//return jsonFiles

		}

	}

}

// 比较两个MD5 list(map)的差异 并且上传文件
func isAnagram(ossClient *oss.Client, config Config, OSS map[string]string, Local map[string]string) {

	//fmt.Println("OSS:  ", OSS)

	for localItem, _ := range Local {
		index := strings.LastIndex(config.LocalDir, "/")

		if OSS[(localItem[index+1:])] == Local[localItem] {
			delete(Local, localItem)
			delete(OSS, (localItem[index+1:]))
		}
	}

	var UpTime int
	var DeleteTime int
	if len(OSS) > 0 && len(Local) == 0 {

		for key, _ := range OSS {
			fmt.Println("需要删除: ", key)
			ossDelete(ossClient, config, key)
			DeleteTime++
		}
		fmt.Printf("删除了 %v 个文件\n", DeleteTime)

	} else if len(OSS) > 0 && len(Local) > 0 {

		for key, _ := range OSS {
			fmt.Println("需要删除: ", key)
			ossDelete(ossClient, config, key)
			DeleteTime++
		}

		index := strings.LastIndex(config.LocalDir, "/")
		for localItem, _ := range Local {
			ossUpload(ossClient, config, localItem[index+1:], localItem)
			fmt.Println("需要上传: ", localItem)
			UpTime++
		}
		fmt.Printf("上传了 %v 个文件\n", UpTime)
		fmt.Printf("删除了 %v 个文件\n", DeleteTime)

	} else {
		fmt.Printf("上传了 %v 个文件\n", UpTime)
		fmt.Printf("删除了 %v 个文件\n", DeleteTime)
	}

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
	fmt.Printf("文件 %v 已从 %v 中删除\n", object, config.BucketName)

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
