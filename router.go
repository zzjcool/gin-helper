package ginHelper

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinRouter interface {
	gin.IRoutes
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
	BasePath() string
}

type GroupRouter struct {
	Path        string   // 路由组的根路径，与Gin的Group一样，定义一组接口的公共路径
	Name        string   // 路由组的名称
	Description string   // 路由组的说明
	Routes      []*Route // 路由组中的具体路由
}

type Route struct {
	Param    Parameter         // 接口的参数
	Path     string            // 接口的路径
	Method   string            // 接口的方法
	Summary  string            // 接口说明
	Handlers []gin.HandlerFunc // 接口的处理函数
}

func (rt *Route) genHandlerFunc() gin.HandlerFunc {
	// 取变量a的反射类型对象
	paramsType := reflect.TypeOf(rt.Param).Elem()
	// 根据反射类型对象创建类型实例
	handler := func(c *gin.Context) {
		newParam := reflect.New(paramsType).Interface().(Parameter)
		err := newParam.Bind(c, newParam)
		if err != nil {
			newParam.Result(c, nil, err)
			return
		}
		data, err := newParam.Handler(c)
		newParam.Result(c, data, err)
	}
	return handler
}

func (rt *Route) AddHandler(r GinRouter) {
	if rt.Param != nil {
		replace := false
		for i, handler := range rt.Handlers {
			if handler == nil {
				//如果内部存在GenHandlerFunc表示占位，那么就替换
				rt.Handlers[i] = rt.genHandlerFunc()
				replace = true
			}
		}
		if !replace {
			rt.Handlers = append(rt.Handlers, rt.genHandlerFunc())
		}
	}
	switch strings.ToUpper(rt.Method) {
	case "GET":
		r.GET(rt.Path, rt.Handlers...)
	case "POST":
		r.POST(rt.Path, rt.Handlers...)
	case "PUT":
		r.PUT(rt.Path, rt.Handlers...)
	case "PATCH":
		r.PATCH(rt.Path, rt.Handlers...)
	case "HEAD":
		r.HEAD(rt.Path, rt.Handlers...)
	case "OPTIONS":
		r.OPTIONS(rt.Path, rt.Handlers...)
	case "DELETE":
		r.DELETE(rt.Path, rt.Handlers...)
	case "ANY":
		r.Any(rt.Path, rt.Handlers...)
	default:
		panic("Method: " + rt.Method + " is error")
	}
}
