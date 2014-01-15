package server

type BufferManager struct {
	buffers chan []byte
}

func NewBufferManager() *BufferManager {
	manager := &BufferManager{
		buffers: make(chan []byte, KMaxOnlineClients*2),
	}
	return manager
}

func (this *BufferManager) Get() []byte {
	select {
	case buffer := <-this.buffers:
		return buffer
	default:
		return make([]byte, KBufferSize)
	}
}

func (this *BufferManager) Put(in []byte) {
	select {
	case this.buffers <- in:
	default:
	}
}
