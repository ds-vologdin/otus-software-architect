package token

// JSON messages
var (
	MsgStatusOK            = []byte("{\"status\": \"OK\"}")
	MsgUserNotFound        = []byte("{\"error\": \"user not found\"}")
	MsgInvalidDataFormat   = []byte("{\"error\": \"data format error\"}")
	MsgCheckPasswordError  = []byte("{\"error\": \"check password error\"}")
	MsgInvalidPassword     = []byte("{\"error\": \"password is invalid\"}")
	MsgInternalError       = []byte("{\"error\": \"internal error\"}")
	MsgRefreshTokenEmpty   = []byte("{\"error\": \"refresh token is empty\"}")
	MsgInvalidRefreshToken = []byte("{\"error\": \"refresh token is invalid\"}")
)
