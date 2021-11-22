package rollsum

type RollingSum struct {
	a, b, count uint16
}

func New() *RollingSum {
	return &RollingSum{}
}

func (s *RollingSum) Init(in []byte) {
	s.Reset()
	s.count = uint16(len(in))
	for i, elem := range in {
		s.a += uint16(elem)
		s.b += (s.count - uint16(i)) * uint16(elem)
	}
}

func (s *RollingSum) Reset() {
	s.count = 0
	s.a = 0
	s.b = 0
}

func (s *RollingSum) Roll(out, in byte) {
	s.a += uint16(in) - uint16(out)
	s.b += s.a - s.count*uint16(out)
}

func (s *RollingSum) Sum() uint32 {
	return (uint32(s.b) << 16) | uint32(s.a)
}
