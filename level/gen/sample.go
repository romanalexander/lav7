package gen

type SampleGenerator struct{}

func (s SampleGenerator) Gen(x, z int32) (bs [16 * 16 * 128][2]byte) {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y < 60; y++ {
				bs[y<<8|z<<4|x] = [2]byte{3, 0}
			}
			bs[60<<8|z<<4|x] = [2]byte{2, 0}
		}
	}
	return
}
