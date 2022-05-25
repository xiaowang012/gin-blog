package index

type SendMessageBoard struct {
	Username       string `form:"username" binding:"required"`
	MessageContent string `form:"content" binding:"required"`
	IfAnonymous    string `form:"anonymous"`
}
