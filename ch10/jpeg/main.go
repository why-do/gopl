// 读取 PNG 图像，并把它作为 JPEG 图像保存
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // 注册 PNG 解码器
	"io"
	"os"
	"path/filepath"
)

func main() {
	fileName := "test"   // 不要扩展名
	dir, _ := os.Getwd() // 返回当前文件路径的字符串和一个err信息，忽略err
	pngPath := filepath.Join(dir, fileName+".png")
	jpgPath := filepath.Join(dir, fileName+".jpg")

	// 打开 png 文件
	pngFile, err := os.Open(pngPath)
	if err != nil {
		// 文件可能不存在
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer pngFile.Close()

	// 创建 jpg 文件
	jpgFile, err := os.Create(jpgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer jpgFile.Close()

	// 调用文件转换
	if err := toJPEG(pngFile, jpgFile); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
