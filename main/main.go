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

var Local_Md5_CACHE = "./Local_Md5_CACHE.json"

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
		getLocalDirFiles(new_userconfig.LocalDir, LocalMd5List)
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

		//client := GetBucketClient(userConfig)

		hugoCommand(userConfig)
		getLocalDirFiles(userConfig.LocalDir, LocalMd5List)

	}

}

func GetBucketClient(bucketConfig Config) *oss.Client {
	client, err := oss.New(bucketConfig.Endpoint, bucketConfig.AccessKeyId, bucketConfig.AccessKeySecret)
	if err != nil {
		fmt.Printf(" oss.New err: %v", err)
	}

	return client
}

func GetOSSMd5List(ossClient *oss.Client, bucketConfig Config) {
	bucket, err := ossClient.Bucket(bucketConfig.BucketName)
	if err != nil {
		fmt.Printf(" ossClient.Bucke err: %v", err)
	}

	lsRes, err := bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		fmt.Printf(" bucket.ListObjects err: %v", err)
	}

	for _, object := range lsRes.Objects {
		meta, _ := bucket.GetObjectMeta(object.Key)

		OSSMd5List = append(OSSMd5List, strings.Trim(meta.Get("Etag"), "\""))

	}

}

func getLocalDirFiles(path string, files map[string]string) {
	dir, _ := os.ReadDir(path)

	for _, fi := range dir {
		if fi.IsDir() {
			fullDir := path + "/" + fi.Name()
			getLocalDirFiles(fullDir, files)
		} else {
			md5num, _ := getMd5(path + fi.Name())

			files[md5num] = path + fi.Name()

			json_files, _ := json.Marshal(files)

			err := os.WriteFile(Local_Md5_CACHE, json_files, 0660)
			if err != nil {
				log.Fatal(err)
			}

		}

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
func ossDelete(client *oss.Client, bucketConfig Config, object string) {
	bucket, err := client.Bucket(bucketConfig.BucketName)
	if err != nil {
		// HandleError(err)
	}

	err = bucket.DeleteObject(object)
	if err != nil {
		// HandleError(err)
	}

}

// OSS上传
func ossUpload() {
	client, err := oss.New("Endpoint", "AccessKeyId", "AccessKeySecret")
	if err != nil {
		// HandleError(err)
	}

	bucket, err := client.Bucket("my-bucket")
	if err != nil {
		// HandleError(err)
	}

	err = bucket.PutObjectFromFile("my-object", "LocalFile")
	if err != nil {
		// HandleError(err)
	}

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
