package errors

type MsgErr struct {
	msg string
}

// NewMsg 创建一个携带面向用户提示信息的 error
func NewMsg(msg string) *MsgErr {
	return &MsgErr{
		msg: msg,
	}
}

func (e *MsgErr) Error() string {
	return e.msg
}
