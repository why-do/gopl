// 根据一个三维曲面函数计算并生成SVG，并输出到Web页面
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	defaultWidth, defaultHeight = 600, 320    // 以像素表示的画布大小
	cells                       = 100         // 网格单元的个数
	xyrange                     = 30.0        // 坐标轴的范围，-xyrange ~ xyrange
	angle                       = math.Pi / 6 // x、y轴的角度，30度
	defaultColor                = "grey"      // 默认的线条颜色
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

var width, height int
var xyscale, zscale float64
var f1 func(x, y float64) float64 // 生成形状的函数
var color string                  // 生成图形的线条颜色

func svg(w io.Writer) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: %s; stroke-width: 0.7' "+
		"width='%d' height='%d'>", color, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, fill, ok := corner(i+1, j)
			if !ok {
				continue
			}
			bx, by, fill, ok := corner(i, j)
			if !ok {
				continue
			}
			cx, cy, fill, ok := corner(i, j+1)
			if !ok {
				continue
			}
			dx, dy, fill, ok := corner(i+1, j+1)
			if !ok {
				continue
			}
			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill: %s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, fill)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func corner(i, j int) (float64, float64, string, bool) {
	// 求出网格单元(i,j)的顶点坐标(x,y)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	// 计算曲面高度z
	z := f1(x, y)
	// 判断z是否是无穷大
	if math.IsInf(z, 0) {
		return 0, 0, "", false
	}
	var fill string
	switch {
	case z > 0.02:
		fill = "#ff0000"
	case z < -0.02:
		fill = "#0000ff"
	default:
		fill = "white"
	}
	// 将(x,y,z)等角投射到二维SVG绘图平面上，坐标是(sx,sy)
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, fill, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // 到(0,0)的距离
	return math.Sin(r) / r
}

// 鸡蛋盒
func eggbox(x, y float64) float64 {
	return (math.Cos(x) + math.Cos(y)) / 10
}

// 马鞍
func saddle(x, y float64) float64 {
	r := y*y/600 - x*x/300
	return r
}

func main() {
	fmt.Println("http://localhost:8000/?color=%2300ff00&f=saddle&w=800&h=600")
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}
		for k, v := range r.Form {
			switch k {
			case "w":
				width, _ = strconv.Atoi(v[0]) // 忽略错误，那么就是0，之后会把0设置成默认值
			case "h":
				height, _ = strconv.Atoi(v[0])
			case "color":
				color = v[0]
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
		if width == 0 {
			width = defaultWidth
		}
		if height == 0 {
			height = defaultHeight
		}
		xyscale = float64(width) / 2 / xyrange // x 或 y 轴上每个单位长度的像素
		zscale = float64(height) * 0.4         // z轴上每个单位长度的像素
		if color == "" {
			color = defaultColor
		}
		if f1 == nil {
			f1 = f
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		svg(w)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}
