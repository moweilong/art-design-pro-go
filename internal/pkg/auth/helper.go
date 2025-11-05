package auth

import "time"

// keyExpired checks if a key has expired, if the value of user.Secret.Expires is 0, it will be ignored.
func keyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}
