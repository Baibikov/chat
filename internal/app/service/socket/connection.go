package socket

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"chat/pkg/utils"
)

type group map[GroupName]map[ConnectName]*connect

func (s *_socket) CreateGroup(group GroupName) error {
	if s._group.exists(group) {
		return ErrExistsGroup
	}

	s._group[group] = make(map[ConnectName]*connect)
	return nil
}

func (s *_socket) GetGroups() []string {
	groups := make([]string, 0, len(s._group))
	for name := range s._group {
		groups = append(groups, name.String())
	}

	return groups
}

func (s *_socket) CreateConnect(group GroupName, ws *websocket.Conn) (Connector, error) {
	if !s._group.exists(group) {
		return nil, ErrExistsGroup
	}

	connName := ConnectName(utils.RandBytesN(10))

	conn := &connect{
		ws:       ws,
		group:    group,
		conn:     connName,
		_message: s._message,
	}

	s._group[group][connName] = conn

	return conn, nil
}

func (g group) exists(group GroupName) bool {
	_, ok := g[group]
	return ok
}

func (g group) connectExists(group GroupName, conn ConnectName) bool {
	_, ok := g[group][conn]
	return ok
}

func (g group) sendMessages(group GroupName, message interface{}) (err error) {
	if !g.exists(group) {
		return ErrUnknownGroup
	}

	for _, conn := range g[group] {
		err = multierr.Append(
			err,
			conn.ws.WriteJSON(message),
		)
	}

	return err
}

func (g group) connectClose(group GroupName, conn ConnectName) error {
	if !g.exists(group) {
		return ErrUnknownGroup
	}

	if !g.connectExists(group, conn) {
		return ErrUnknownConnect
	}

	err := g[group][conn].ws.Close()
	delete(g[group], conn)

	return err
}

type Payload struct {
	Body []byte
	resp interface{}
}

func (p *Payload) Write(v interface{}) {
	p.resp = v
}

type HandleFunc func(p *Payload) error

// Connector - implement socket connect usage
type Connector interface {
	Read(ctx context.Context, handler HandleFunc) error
}

type connect struct {
	ws    *websocket.Conn
	group GroupName
	conn  ConnectName

	_message chan<- messageChan
}

// Read - method takes into account all possible
// errors of sending and receiving messages via web socket
// takes into account the errors of the user's handler
// and sends a message on arrival
func (c *connect) Read(ctx context.Context, handleFunc HandleFunc) (err error) {
	for {
		select {
		case <-ctx.Done():
			c.sendError(err)
			return errors.New("ctx closed")
		default:
			_, p, err := c.ws.ReadMessage()
			if err != nil {
				c.sendError(err)
				return err
			}
			if len(p) == 0 {
				continue
			}

			r := &Payload{
				Body: p,
			}

			if err := handleFunc(r); err != nil {
				c.sendError(err)
				return err
			}

			c.sendSuccess(r.resp)
		}
	}
}

// sendError - send err message
func (c *connect) sendError(err error) {
	c._message <- messageChan{
		group: c.group,
		conn:  c.conn,
		err:   err,
	}
}

// sendSuccess - send success message
func (c *connect) sendSuccess(payload interface{}) {
	c._message <- messageChan{
		group:   c.group,
		conn:    c.conn,
		payload: payload,
	}
}
