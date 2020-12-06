package profile

// JSON messages
var (
	MsgStatusOK          = []byte("{\"status\": \"OK\"}")
	MsgInvalidUserID     = []byte("{\"error\": \"invalid id\"}")
	MsgUserNotFound      = []byte("{\"error\": \"user not found\"}")
	MsgInvalidDataFormat = []byte("{\"error\": \"data format error\"}")
	MsgInternalError     = []byte("{\"error\": \"internal error\"}")
	MsgGetUserError      = []byte("{\"error\": \"get user error\"}")
	MsgCreateUserError   = []byte("{\"error\": \"create user error\"}")
	MsgDeleteUserError   = []byte("{\"error\": \"delete user error\"}")
	MsgEditUserError     = []byte("{\"error\": \"edit user error\"}")
)
