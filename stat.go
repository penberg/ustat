package ustat

// A Stat is a collection of named stats.
type Stat struct {
	Names        []string
	Descriptions []string
	Collector    StatCollector
}

// A StatCollector is an interface for collecting stats.
type StatCollector interface {
	Collect() []uint64
}

// Difference calculates the change in values for two arrays.
func Difference(before []uint64, after []uint64) []uint64 {
	var diff []uint64
	for idx := range before {
		diff = append(diff, after[idx]-before[idx])
	}
	return diff
}
