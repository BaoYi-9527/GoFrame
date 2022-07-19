package main

import "strings"

type node struct {
	pattern  string  // 待匹配的路由 例如：/p/:lang
	part     string  // 路由中的一部分 例如：:lang
	children []*node // 子节点 例如：[doc, tutorial, intro]
	isWild   bool    // 是否精准匹配, part 含有 : 或 * 时为 true
}

// matchChild
// @Description: 第一个匹配成功的节点，用于插入
// @receiver n
// @param part
// @return *node
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 若 part 等于子节点的 part 或者是模糊匹配就返回子节点
		if (child.part == part) || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren
// @Description: 返回所有匹配成功的节点 用于查找
// @receiver n
// @param part
// @return []*node
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert
// @Description:
// @receiver n
// @param pattern	待匹配的模式
// @param parts
// @param height
func (n *node) insert(pattern string, parts []string, height int) {
	// 若深度遍历到了路由的深度 则表示遍历结束
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 获取当前深度的路由的 part 然后去当前节点的子节点中去寻找这个 part
	// 找得到就说明 已经在子节点中 -> 遍历下一层
	// 找不到就新增该节点
	part := parts[height]
	child := n.matchChild(part) // 找不到这个路由
	if child == nil {
		// 若 part 以 : 或 * 开头 则是模糊匹配
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 将当前路由添加到 node 的子节点中
		n.children = append(n.children, child)
	}
	// 递归 下一层
	child.insert(pattern, parts, height+1)
}

// search
// @Description:查找符合路由规则的节点
// @receiver n
// @param parts
// @param height
// @return *node
func (n *node) search(parts []string, height int) *node {
	// 如果遍历到了最底层 或者 当前节点的 part 是模糊匹配则进入到了最后一次
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果其没有后续路由 则表明匹配失败
		if n.pattern == "" {
			return nil
		}
		// 返回模糊匹配节点
		return n
	}
	// 遍历 part
	part := parts[height]
	// 查找所有符合条件的子节点
	children := n.matchChildren(part)
	// 遍历子节点
	for _, child := range children {
		// 查找下一个 节点
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	// 查找失败
	return nil
}
