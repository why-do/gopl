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
	zscale        = height * 0.04       // z轴上每个单位长度的像素
	angle         = math.Pi / 6         // x、y轴的角度，30度
	color         = "grey"              // 默认的线条颜色
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

var f1 func(x, y float64) float64 // 生成形状的函数
var webColor string               // 生成图形的线条颜色

func svg(w io.Writer) {
	if webColor == "" {
		webColor = color
	}
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: %s; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", webColor, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, ok := corner(i+1, j)
			if !ok {
				continue
			}
			bx, by, ok := corner(i, j)
			if !ok {
				continue
			}
			cx, cy, ok := corner(i, j+1)
			if !ok {
				continue
			}
			dx, dy, ok := corner(i+1, j+1)
			if !ok {
				continue
			}
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func corner(i, j int) (float64, float64, bool) {
	// 求出网格单元(i,j)的顶点坐标(x,y)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	// 计算曲面高度z
	if f1 == nil {
		f1 = f
	}
	z := f1(x, y)
	// 判断z是否是无穷大
	if math.IsInf(z, 0) {
		return 0, 0, false
	}
	// 将(x,y,z)等角投射到二维SVG绘图平面上，坐标是(sx,sy)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // 到(0,0)的距离
	return math.Sin(r)
}

// 鸡蛋盒
func eggbox(x, y float64) float64 {
	return 0.2 * (math.Cos(x) + math.Cos(y))
}

// 马鞍
func saddle(x, y float64) float64 {
	r := y*y/44 - x*x/44
	return r
}

func main() {
	fmt.Println("http://localhost:8000/?color=gray&f=f")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		for k, v := range r.Form {
			switch k {
			case "color":
				webColor = v[0]
			case "f":
				switch v[0] {
				case "eggbox":
					f1 = eggbox
				case "saddle":
					f1 = saddle
				default:
					f1 = f
				}
			}
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		svg(w)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}
