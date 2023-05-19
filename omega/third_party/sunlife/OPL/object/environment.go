package object

// 创建一个新的环境
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// 创建一个新的环境
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// 环境 绑定变量名字的
type Environment struct {
	store map[string]Object
	outer *Environment
}

// 获取根据变量名字 获取是否有绑定值
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// 设置变量名字与绑定值
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
