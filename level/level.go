package level

// Level is a interface for storing block chunks, entities, tile objects, etc.
type Level interface {
	GetChunk(int64, int64) *Chunk
	GetName() string
}
