package message

// message type
const (
	KNil           uint32 = iota // begin
	KLoginRequest                // login request
	KLoginResopnse               // login response
	KEnd                         //end
)
