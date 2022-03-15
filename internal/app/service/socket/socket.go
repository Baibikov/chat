package socket

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type (
	// GroupName - name of socket group
	GroupName string

	// ConnectName - name of socket connection
	ConnectName string
)

func (g GroupName) String() string {
	return string(g)
}

// IOSocket - merge use cases for current
// method uses
type IOSocket interface {
	Listener
	Grouper
}

// Listener - implement listen case
type Listener interface {
	ListenAndRoute(ctx context.Context) error
}

type _socket struct {
	_message chan messageChan
	_group   group
}

func NewSocket() IOSocket {
	return &_socket{
		_message: make(chan messageChan),
		_group:   map[GroupName]map[ConnectName]*connect{},
	}
}

// Grouper - distributing interface
// 1. create group
// 2. create connect
// 3. get created groups
type Grouper interface {
	CreateGroup(group GroupName) error
	CreateConnect(group GroupName, ws *websocket.Conn) (Connector, error)
	GetGroups() []string
}

// ListenAndRoute - listen messages for groups
// and send for sockets
func (s *_socket) ListenAndRoute(ctx context.Context) error {
	return s.sender(ctx)
}

// sender - the main task is to distribute messages into groups
// and check if the application context has been corrupted
// also the sender to monitor the errors of sent messages
func (s *_socket) sender(ctx context.Context) (err error) {
	defer func(err error) {
		multierr.AppendInto(&err, s.shutdown())
	}(err)
	for {
		select {
		case <-ctx.Done():
			return errors.New("ctx closed")
		case m, ok := <-s._message:
			if !ok {
				return nil
			}
			if m.err != nil {
				logrus.Infof("connect message err: %+v", m.err)
				err = s._group.connectClose(m.group, m.conn)
				if err != nil {
					logrus.Info(err)
				}

				continue
			}

			err = s._group.sendMessages(m.group, m.payload)
			if err != nil {
				logrus.Info(err)
			}
		}
	}
}

// shutdown - close all group connections
func (s *_socket) shutdown() (err error) {
	for name, g := range s._group {
		logrus.Infof("close [%s] group", name)
		for conn := range g {
			err = multierr.Append(err, s._group.connectClose(name, conn))
		}
	}

	close(s._message)
	return err
}
