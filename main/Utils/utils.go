package Utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// GetMd5 计算MD5
func GetMd5(filepath string) (string, error) {
	f, _ := os.Open(filepath)
	defer f.Close()

	body, _ := io.ReadAll(f)

	md5sum := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()

	return strings.ToUpper(md5sum), nil
}
