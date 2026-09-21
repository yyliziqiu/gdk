package xcq

// Handle 处理元素
type Handle func(item any)

// Filter 元素符合条件返回 true，否则返回 false
type Filter func(item any) bool

// Remove 元素需要删除返回 true，否则返回 false
type Remove func(item any) bool
