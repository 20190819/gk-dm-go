package demo

import "strings"

// router 路由
type router struct {
	tree map[string]*node
}

func (r *router) addRoute(method, path string, handle HandleFunc) {

	if r.tree == nil {
		r.tree = make(map[string]*node)
	}

	root, ok := r.tree[method]
	if !ok {
		root = &node{path: "/"}
		r.tree[method] = root
	}
	if path == "/" {
		root.handler = handle
		return
	}

	// 去除前、后的 /
	path = strings.Trim(path, "/")

	rootCopy := root
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg != "" {
			rootCopy = rootCopy.childOrCreate(seg)
		}
	}
	rootCopy.handler = handle
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.tree[method]
	// 没找到
	if !ok {
		return nil, false
	}

	// 如果是根节点
	if path == "/" {
		return &matchInfo{n: root}, true
	}

	// 去除两端的 /
	path = strings.Trim(path, "/")
	cur := root

	// 将路由分段
	segs := strings.Split(path, "/")
	for _, seg := range segs {

		// 没有子节点，
		// 可能是参数 /:xxx/ 路由
		// 可能是通配符 * 路由
		if cur.Children == nil {
			if cur.ParamsChild != nil {
				return &matchInfo{
					n: cur.ParamsChild,
					pathParams: map[string]string{
						cur.ParamsChild.path[1:]: seg,
					},
				}, true
			}
			return &matchInfo{n: cur.StarNode}, cur.StarNode != nil
		}

		// /order/cancel
		child, ok := cur.Children[seg]
		// 有子节点，但没找到 key=seg  这个，
		// 可能是参数 /:xxx/ 路由
		// 可能是通配符 * 路由
		if !ok {
			if cur.ParamsChild != nil {
				return &matchInfo{
					n: cur.ParamsChild,
					pathParams: map[string]string{
						cur.ParamsChild.path[1:]: seg,
					},
				}, true
			}
			return &matchInfo{n: cur.StarNode}, cur.StarNode != nil
		}
		// 找到了子节点
		cur = child
	}

	// 最终返回
	return &matchInfo{n: cur}, true
}

type node struct {
	path        string
	handler     HandleFunc
	StarNode    *node
	Children    map[string]*node
	ParamsChild *node
}

func (n *node) childOrCreate(path string) *node {

	// a/*/c
	if path == "*" {
		if n.StarNode == nil {
			n.StarNode = &node{path: path}
		}
		return n.StarNode
	}

	// 含参数且 ParamsChild==nil
	// 初始化 ParamsChild
	if path[0] == ':' && n.ParamsChild == nil {
		n.ParamsChild = &node{path: path}
		return n.ParamsChild
	}

	if n.Children == nil {
		n.Children = make(map[string]*node)
	}

	// 如果没有则先构建节点
	if _, ok := n.Children[path]; !ok {
		n.Children[path] = &node{path: path}
	}
	return n.Children[path]
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
