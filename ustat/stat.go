package ustat

type Stat struct {
	Names        []string
	Descriptions []string
	Reader       StatReader
}

type StatReader interface {
	Read() []uint64
}
