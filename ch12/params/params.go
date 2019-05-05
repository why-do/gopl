package params

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Unpack 从 HTTP 请求 req 的参数中提取数据填充到 ptr 指向的结构体的各个字段
func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	// 创建字段映射表，key 为有效名称
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem() // reflect.ValueOf(&x).Elem() 获得任意变量 x 可寻址的值，用于设置值。
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.Structtag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// 对请求中的每个参数更新结构体中对应的字段
	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue // 忽略不能识别的 HTTP 参数
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
