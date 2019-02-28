// 生成一个PNG格式的Mandelbrot分形图
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/height*(xmax-xmin) + xmin
			z := complex(x, y)
			// 点(px, py)表示复数值z
			img.Set(px, py, mandelbrot(z))
		}
	}
	// png.Encode(os.Stdout, img) // 注意：忽略错误
	f, err := os.OpenFile("p1.png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	defer f.Close()
	png.Encode(f, img)  // 注意：忽略错误
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
