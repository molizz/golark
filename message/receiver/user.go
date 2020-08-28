package receiver

// 这里暂时实现user、未来可能有针对部门的信息发送等等...
//
type UserReceiver struct {
	receivers []string
}

func NewUserReceiver(receivers ...string) *UserReceiver {
	return &UserReceiver{
		receivers: receivers,
	}
}

func (u *UserReceiver) Type() string {
	return "user"
}

func (u *UserReceiver) ParamName() string {
	return "open_ids"
}

func (u *UserReceiver) List() (receivers []string) {
	return u.receivers
}
