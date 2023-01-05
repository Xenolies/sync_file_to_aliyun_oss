package LocalOperation

import (
	"GolangAliOSSUpload/main/BucketOperation"
	"GolangAliOSSUpload/main/Config"
	"GolangAliOSSUpload/main/Utils"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"strings"
)

func GetLocalDirList(path string, config Config.Config, files map[string]string) {

	dir, _ := os.ReadDir(path)

	for _, fi := range dir {
		if fi.IsDir() {
			fullDir := path + fi.Name() + "/"
			GetLocalDirList(fullDir, config, files)
		} else {

			md5num, _ := Utils.GetMd5(path + fi.Name())

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
func Comparison(ossClient *oss.Client, config Config.Config, OSS map[string]string, Local map[string]string) {

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
			BucketOperation.OssDelete(ossClient, config, key)
			DeleteTime++
		}
		fmt.Printf("删除了 %v 个文件\n", DeleteTime)

	} else if len(OSS) > 0 && len(Local) > 0 {

		for key, _ := range OSS {
			fmt.Println("需要删除: ", key)
			BucketOperation.OssDelete(ossClient, config, key)
			DeleteTime++
		}

		index := strings.LastIndex(config.LocalDir, "/")
		for localItem, _ := range Local {
			BucketOperation.OssUpload(ossClient, config, localItem[index+1:], localItem)
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
