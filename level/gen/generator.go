package gen

// Generator is an interface for MCPE map generator.
type Generator interface {
	Gen(int32, int32) [16 * 16 * 128][2]byte
}
