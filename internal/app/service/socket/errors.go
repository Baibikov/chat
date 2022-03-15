package socket

import "github.com/pkg/errors"

var ErrUnknownGroup = errors.New("unknown group")
var ErrUnknownConnect = errors.New("unknown connect")
var ErrExistsGroup = errors.New("group exists")
