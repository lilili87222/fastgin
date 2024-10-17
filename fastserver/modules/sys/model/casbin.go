package model

// 角色权限规则
type RoleCasbin struct {
	Keyword string // 角色关键字
	Path    string // 访问路径
	Method  string // 请求方式
}
