package custom_errors

import "errors"

var (
	UserNotFound = errors.New("user not found")
	SegmentNotFound = errors.New("segment not found")
	SegmentNameEmpty = errors.New("segments name empty")
	SegmentNameExists = errors.New("segment with this name exists")
	UserAlreadyExistsInSegment = errors.New("user already exists in segment")
	UserSegmentNotFound = errors.New("user not found in segment")
	UsersArrayEmpty = errors.New("users array empty")
	ZeroUsers = errors.New("zero users to add")
)