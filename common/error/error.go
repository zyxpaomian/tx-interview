package error

type AgentError struct {
	msg string
}

func (a *AgentError) Error() string {
	return a.msg
}

func New(msg string) *AgentError {
	return &AgentError{msg: msg}
}

func DBError() *AgentError {
	return New("[DB状态异常] DB状态异常")
}

func AuthError() *AgentError {
	return New("[认证异常] 认证异常")
}
