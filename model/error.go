package model

import "errors"

var ErrNoBytes = errors.New("no bytes received")
var ErrBadResponse = errors.New("bad response code received")
