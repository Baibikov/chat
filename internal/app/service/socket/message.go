package socket

type messageChan struct {
	group   GroupName
	conn    ConnectName
	payload interface{}
	err     error
}
