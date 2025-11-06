package errno

import (
	"net/http"

	"github.com/moweilong/milady/pkg/errorsx"
)

var (
	// ErrUsernameInvalid indicates that the username is invalid.
	ErrUsernameInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.UsernameInvalid",
		Message: "Invalid username: Username must consist of letters, digits, and underscores only, and its length must be between 3 and 20 characters.",
	}

	// ErrNicknameInvalid indicates that the nickname is invalid.
	ErrNicknameInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.NicknameInvalid",
		Message: "Invalid nickname: Nickname must consist of letters, digits, and underscores only, and its length must be between 3 and 30 characters.",
	}

	// ErrPasswordInvalid indicates that the password is invalid.
	ErrPasswordInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.PasswordInvalid",
		Message: "Password is incorrect.",
	}

	// ErrUserAlreadyExists indicates that the user already exists.
	ErrUserAlreadyExists = &errorsx.ErrorX{Code: http.StatusBadRequest, Reason: "AlreadyExist.UserAlreadyExists", Message: "User already exists."}

	// ErrUserNotFound indicates that the specified user was not found.
	ErrUserNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.UserNotFound", Message: "User not found."}
)
