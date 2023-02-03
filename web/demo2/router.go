package demo

import "strings"

// router 路由
type router struct {
	tree map[string]*node
}

func (r *router) addRoute(method, path string, handle HandleFunc) {

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
}

type node struct {
	path        string
	handler     HandleFunc
	starNode    *node
	Children    map[string]*node
	ParamsChild *node
}

func (n *node) childOrCreate(path string) *node {

	// a/*/c
	if path == "*" {
		if n.starNode == nil {
			n.starNode = &node{path: path}
		}
		return n.starNode
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
