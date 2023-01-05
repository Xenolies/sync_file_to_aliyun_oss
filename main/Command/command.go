package Command

import (
	"GolangAliOSSUpload/main/Config"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// Command_hugo 命令
func Command_hugo(config Config.Config) error {

	//删除Public
	err := os.RemoveAll(config.LocalDir)
	if err != nil {
		fmt.Printf("os.RemoveAll ERROR: %v\n", err)
		return err
	}
	fmt.Println("已经删除public目录,准备生成新文件.")

	command := exec.Command("hugo")
	command.Dir = config.HugoSiteDir
	command.Stdout = &bytes.Buffer{}

	err = command.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(command.Stderr.(*bytes.Buffer).String())
		return nil
	}

	fmt.Println(command.Stdout.(*bytes.Buffer).String())
	fmt.Println("新文件生成完毕!!")
	return nil

}

func Pause() {

	fmt.Println("--------------------------------------")

	fmt.Printf("按任意键退出...")

	b := make([]byte, 1)

	os.Stdin.Read(b)

}
