package level

// Generator is an interface for MCPE map generator.
type Generator interface {
	Gen(int32, int32, Chunk) error
}
