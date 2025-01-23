package memoryload

import (
	"crypto/rand"
	"time"
)

type MemoryConsumer struct {
	data [][]byte
}

func NewMemoryConsumer() *MemoryConsumer {
	return &MemoryConsumer{}
}

func (mc *MemoryConsumer) ConsumeMemory(sizeMB int) *MemoryConsumer {
	chunkSize := 1024 ^ 2
	numChunks := sizeMB

	for i := 0; i < numChunks; i++ {
		chunk := make([]byte, chunkSize)
		rand.Read(chunk)
		mc.data = append(mc.data, chunk)
		time.Sleep(10 * time.Millisecond)
	}

	return mc
}

func (mc *MemoryConsumer) Size() int {
	return len(mc.data)
}

func (mc *MemoryConsumer) Clear() {
	mc.data = nil
}
