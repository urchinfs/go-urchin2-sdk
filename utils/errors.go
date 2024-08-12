package utils

import "errors"

var (
	ErrNotDir        = errors.New("this node is not a directory")
	ErrNotFile       = errors.New("this node is not a regular file")
	ErrOffline       = errors.New("this action must be run in online mode, try running 'ipfs daemon' first")
	ErrNotSupported  = errors.New("operation not supported")
	ErrNotReceiveRet = errors.New("no results received from ipfs peer")
	ErrBadResponse   = errors.New("bad response from server")
)
