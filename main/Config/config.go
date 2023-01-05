package Config

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
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

// ConfigMaker 配置文件生成引导程序
func ConfigMaker() bool {

	dir, _ := os.ReadDir(configPath)
	if dir == nil {
		//MakeConfigFile(configDEMO)
		fmt.Println("未检测到配置文件,生成config.json文件,请配置: ")
		fmt.Println("..............创建配置文件........")

		newUserconfig := Config{
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
		newUserconfig.AccessKeyId = strings.Trim(AccessKeyId, "\r\n")
		fmt.Printf("AccessKeyId: %v\n", newUserconfig.AccessKeyId)
		fmt.Println("........................................................")

		fmt.Println("请输入AccessKeySecret(输入回车结束): ")
		AccessKeySecret, _ := reader.ReadString('\n')
		newUserconfig.AccessKeySecret = strings.Trim(AccessKeySecret, "\r\n")
		fmt.Printf("AccessKeySecret: %v\n", newUserconfig.AccessKeySecret)
		fmt.Println("........................................................")

		fmt.Println("请输入Endpoint (输入回车结束 ,Endpoint格式类似这样 oss-cn-hongkong.aliyuncs.com): ")
		Endpoint, _ := reader.ReadString('\n')
		newUserconfig.Endpoint = strings.Trim(Endpoint, "\r\n")
		fmt.Printf("Endpoint: %v\n", newUserconfig.Endpoint)
		fmt.Println("........................................................")

		fmt.Println("请输入BucketName: ")
		BucketName, _ := reader.ReadString('\n')
		newUserconfig.BucketName = strings.Trim(BucketName, "\r\n")
		fmt.Printf("BucketName: %v\n", newUserconfig.BucketName)
		fmt.Println("........................................................")

		fmt.Println("请输入LocalDir (LocalDir是你本地博客目录下的public目录), ")
		fmt.Println("格式类似这样 D:/OdinXu/odinxu-blog/public/")
		fmt.Println("请输入(输入回车结束): ")
		LocalDir, _ := reader.ReadString('\n')
		newUserconfig.LocalDir = strings.Trim(LocalDir, "\r\n")
		fmt.Printf("LocalDir: %v\n", newUserconfig.LocalDir)
		fmt.Println("........................................................")

		fmt.Println("请输入HugoSiteDir (HugoSiteDir是你本地博客所在目录), ")
		fmt.Println("格式类似这样 D:/OdinXu/odinxu-blog/")
		fmt.Println("请输入(输入回车结束): ")
		HugoSiteDir, _ := reader.ReadString('\n')
		newUserconfig.HugoSiteDir = strings.Trim(HugoSiteDir, "\r\n")
		fmt.Printf("HugoSiteDir: %v\n", newUserconfig.HugoSiteDir)
		fmt.Println("........................................................")

		MakeConfigFile(newUserconfig)
		fmt.Println("配置文件创建成功!五秒后将关闭程序")
		// getLocalDirList(new_userconfig.LocalDir, LocalMd5List)
		time.Sleep(5 * time.Second)

	}
	return true

}

var configPath string = "./config.json"

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

// LogSave 保存日志文件
func LogSave() {

}
