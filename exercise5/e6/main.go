// 根据一个三维曲面函数计算并生成SVG，并输出到Web页面
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

const (
	width, height = 600, 320            // 以像素表示的画布大小
	cells         = 100                 // 网格单元的个数
	xyrange       = 30.0                // 坐标轴的范围，-xyrange ~ xyrange
	xyscale       = width / 2 / xyrange // x 或 y 轴上每个单位长度的像素
	zscale        = height * 0.4        // z轴上每个单位长度的像素
	angle         = math.Pi / 6         // x、y轴的角度，30度
	color         = "grey"              // 线条的颜色
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func svg(w io.Writer) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: %s; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", color, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func corner(i, j int) (sx, sy float64) {
	// 求出网格单元(i,j)的顶点坐标(x,y)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	// 计算曲面高度z
	z := f(x, y)
	// 将(x,y,z)等角投射到二维SVG绘图平面上，坐标是(sx,sy)
	return width/2 + (x-y)*cos30*xyscale, height/2 + (x+y)*sin30*xyscale - z*zscale
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // 到(0,0)的距离
	return math.Sin(r) / r
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		svg(w)
	}
	http.HandleFunc("/", handler)
	fmt.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}