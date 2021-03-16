package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type CopyInfo struct {
	src1 string
	dst1 string
	src2 string
	dst2 string
}

func main() {
	flag.Parse()
	var path = flag.Arg(0)

	Walk(path)
}

func Walk(root string) {
	copyInfos := make([]*CopyInfo, 0)

	{
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			err, copyInfo := process(path, info, err)
			if err != nil {
				return err
			}
			if copyInfo != nil {
				copyInfos = append(copyInfos, copyInfo)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error on filepath.Walk : ", err)
		}
	}

	{
		for _, ci := range copyInfos {
			// 親ディレクトリの名前を変更
			fmt.Printf("親: %s -> %s\n", ci.src1, ci.dst1)
			if err := os.Rename(ci.src1, ci.dst1); err != nil {
				fmt.Println("Error on os.Rename1 : ", err)
			}
			// 子ディレクトリの名前を変更
			fmt.Printf("子: %s -> %s\n", ci.src2, ci.dst2)
			if err := os.Rename(ci.src2, ci.dst2); err != nil {
				fmt.Println("Error on os.Rename2 : ", err)
			}
		}
	}

	fmt.Println("END")
}

func process(path string, info os.FileInfo, err error) (error, *CopyInfo) {
	if err != nil {
		return err, nil
	}
	if info.IsDir() {
		// 自分がディレクトリで親の名前と同じかチェック
		dirName := filepath.Base(path)
		parent := filepath.Dir(path)
		parentName := filepath.Base(parent)
		pparent := filepath.Dir(parent)
		backName := filepath.Join(pparent, parentName+".bak")
		if dirName == parentName {
			// 親ディレクトリの名前を変更
			src1 := filepath.Join(pparent, parentName)
			dst1 := backName
			// 子ディレクトリの名前を変更
			src2 := filepath.Join(backName, dirName)
			dst2 := filepath.Join(pparent, dirName)
			copyInfo := new(CopyInfo)
			copyInfo.src1 = src1
			copyInfo.dst1 = dst1
			copyInfo.src2 = src2
			copyInfo.dst2 = dst2
			return nil, copyInfo
		}
	}
	return nil, nil
}
