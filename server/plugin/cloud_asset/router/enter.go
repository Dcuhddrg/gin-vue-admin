package router

type RouterGroup struct {
	CloudProviderRouter
	CloudInstanceRouter
}

var RouterGroupApp = new(RouterGroup)
