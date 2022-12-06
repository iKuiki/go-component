package utils

import "container/list"

// Stack 栈
type Stack struct {
	list *list.List
}

// NewStack 创建栈
func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

// Push 将元素推入栈
func (stack *Stack) Push(value interface{}) {
	stack.list.PushBack(value)
}

// Pop 将元素从栈顶取出
func (stack *Stack) Pop() interface{} {
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		return e.Value
	}
	return nil
}

// Peak 检查栈顶元素
func (stack *Stack) Peak() interface{} {
	e := stack.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

// Len 获取栈长度
func (stack *Stack) Len() int {
	return stack.list.Len()
}

// IsEmpty 检查栈是否为空
func (stack *Stack) IsEmpty() bool {
	return stack.list.Len() == 0
}
