package demo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_addRoute(t *testing.T) {

	testRoutes := []struct {
		method string
		path   string
	}{
		// 静态路由
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "//home",
		},
		{
			method: http.MethodGet,
			path:   "//home1////",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/cancel",
		},
		{
			method: http.MethodPost,
			path:   "/order/detail/:order_sn",
		},
	}

	var exceptRouter *router = &router{
		tree: make(map[string]*node),
	}
	var myHandleFunc HandleFunc = func(ctx *context.Context) {}

	wantRouter := &router{
		tree: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: myHandleFunc,
				Children: map[string]*node{
					"home":  &node{path: "home", handler: myHandleFunc},
					"home1": &node{path: "home1", handler: myHandleFunc},
					"user": &node{
						path:    "user",
						handler: myHandleFunc,
						Children: map[string]*node{
							"detail": &node{
								path: "detail",
								Children: map[string]*node{
									"profile": &node{
										path:    "profile",
										handler: myHandleFunc,
									},
								},
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				Children: map[string]*node{
					"order": &node{
						path: "order",
						StarNode: &node{
							path:    "*",
							handler: myHandleFunc,
						},
						Children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: myHandleFunc,
							},
							"detail": &node{
								path: "detail",
								ParamsChild: &node{
									path:    ":order_sn",
									handler: myHandleFunc,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testRoute := range testRoutes {
		exceptRouter.addRoute(testRoute.method, testRoute.path, myHandleFunc)
	}

	resMsg, ok := wantRouter.myEqual(exceptRouter)
	assert.True(t, ok, resMsg) // 测试-断言

	findCases := []struct {
		name   string
		method string
		path   string

		found      bool
		wantPath   string
		hasHandler bool
	}{
		{
			method:     http.MethodGet,
			name:       "/",
			path:       "/",
			wantPath:   "/",
			found:      true,
			hasHandler: true,
		},
		{
			method:     http.MethodGet,
			name:       "/user",
			path:       "/user",
			wantPath:   "user",
			found:      true,
			hasHandler: true,
		},
		{
			method:     http.MethodGet,
			name:       "/user/detail",
			path:       "/user/detail",
			wantPath:   "detail",
			found:      true,
			hasHandler: false,
		},
		{
			method:     http.MethodGet,
			name:       "/user/detail/profile",
			path:       "/user/detail/profile",
			wantPath:   "profile",
			found:      true,
			hasHandler: true,
		},
		{
			method:     http.MethodPost,
			name:       "/order/*",
			path:       "/order/abc",
			wantPath:   "*",
			found:      true,
			hasHandler: true,
		},
		{
			method:     http.MethodPost,
			name:       "/order/:xxx",
			path:       "/order/detail/:order_sn",
			wantPath:   ":order_sn",
			found:      true,
			hasHandler: true,
		},
	}
	for _, caseItem := range findCases {
		t.Run(caseItem.name, func(t *testing.T) {
			r, ok := exceptRouter.findRoute(caseItem.method, caseItem.path)
			assert.Equal(t, caseItem.found, ok)
			if ok {
				assert.Equal(t, caseItem.wantPath, r.n.path)
				assert.Equal(t, caseItem.hasHandler, r.n.handler != nil)
			}
		})
	}
}

func (r *router) myEqual(except *router) (string, bool) {

	for methodName, node := range r.tree {
		subTree, ok := except.tree[methodName]
		if !ok {
			return fmt.Sprintf("目标树里面没有方法为 %s 的路由树", methodName), false
		}
		msg, ok := node.nodeEqual(subTree)
		if !ok {
			return fmt.Sprintf("%s 测试未通过：%v", methodName, msg), false
		}
	}
	return "", true
}

func (n *node) nodeEqual(ans *node) (string, bool) {

	var msg string = ""

	if ans == nil {
		msg = "目标节点是 nil"
		return msg, false
	}

	if n.path != ans.path {
		msg = fmt.Sprintf("节点路径不相等,%v,%v", n.path, ans.path)
		return msg, false
	}

	// 反射取值
	nHandleFunc := reflect.ValueOf(n.handler)
	ansHandleFunc := reflect.ValueOf(ans.handler)
	if nHandleFunc.Type().String() != ansHandleFunc.Type().String() {
		msg = fmt.Sprintf("节点 handler 不相等，%s,%s", nHandleFunc.Type().String(), ansHandleFunc.Type().String())
		return msg, false
	}

	if len(n.Children) != len(ans.Children) {
		msg = fmt.Sprintf("节点子节点数量不相等,期望:%d,实际:%d", len(n.Children), len(ans.Children))
		return msg, false
	}

	if len(n.Children) == 0 {
		return "", true
	}

	for name, child := range n.Children {
		ansChild, ok := ans.Children[name]
		if !ok {
			msg = fmt.Sprintf("目标节点，缺少子节点:%s", name)
			return msg, false
		}

		res, ok := child.nodeEqual(ansChild)
		if !ok {
			msg = fmt.Sprintf("%s: %s", name, res)
			return msg, false
		}
	}
	return "", true
}
