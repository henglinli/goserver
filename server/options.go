package server

// ===========================================
// server options
const (
	kDataBasePath     = "user.db"
	KMaxOnlineClients = 2048
	KBufferSize       = 2048
	kDefaultMagic     = "magic" // defualt messsage header magic
	kDefaultMagicLen  = len(kDefaultMagic)
)
