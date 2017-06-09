package ustat

type Stat struct {
	Names        []string
	Descriptions []string
	Reader       StatReader
}

type StatReader interface {
	Read() []uint64
}

// Difference calculates the change in values for two arrays.
func Difference(before []uint64, after []uint64) []uint64 {
	var diff []uint64
	for idx := range before {
		diff = append(diff, after[idx]-before[idx])
	}
	return diff
}