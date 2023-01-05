package BucketOperation

import (
	"GolangAliOSSUpload/main/Config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

func GetBucketClient(bucketConfig Config.Config) *oss.Client {
	client, err := oss.New(bucketConfig.Endpoint, bucketConfig.AccessKeyId, bucketConfig.AccessKeySecret)
	if err != nil {
		fmt.Printf(" oss.New err: %v", err)
	}

	return client
}

func GetOSSMd5List(ossClient *oss.Client, config Config.Config, ossList map[string]string) map[string]string {
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

// OssDelete 删除OSS文件
func OssDelete(client *oss.Client, config Config.Config, object string) {
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

// OssUpload OSS上传
func OssUpload(client *oss.Client, config Config.Config, fileName string, filePath string) {
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
