package mistakes

import "errors"

var (
	ErrBucketNotFound  = errors.New("bucket not found")
	ErrNetNotFound     = errors.New("net not found")
	ErrNetExist        = errors.New("net is already exist")
	ErrNetAnotherExist = errors.New("net is already exist in another list")
)
