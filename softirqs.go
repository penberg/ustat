package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type procSoftIRQsCollector struct {
	counts []uint64
}

const procSoftIRQsPath = "/proc/softirqs"

// NewSoftIRQsStat returns a new Stat, which collects interrupt stats from /proc/interrupts.
func NewSoftIRQsStat() *Stat {
	interrupts, err := procfs.ReadInterrupts(procSoftIRQsPath)
	if err != nil {
		panic(err)
	}
	names := parseSoftIRQNames(interrupts)
	descriptions := parseSoftIRQDescriptions(interrupts)
	counts := parseSoftIRQCounts(interrupts)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Collector:    &procSoftIRQsCollector{counts: counts},
	}
}

func (reader *procSoftIRQsCollector) Collect() []uint64 {
	interrupts, err := procfs.ReadInterrupts(procSoftIRQsPath)
	if err != nil {
		panic(err)
	}
	counts := parseSoftIRQCounts(interrupts)
	diff := Difference(reader.counts, counts)
	reader.counts = counts
	return diff
}

func parseSoftIRQNames(interrupts *procfs.Interrupts) []string {
	var names []string
	for _, interrupt := range interrupts.Interrupts {
		for cpu := range interrupt.Counts {
			name := fmt.Sprintf("softirq.%s.cpu%d", interrupt.Name, cpu)
			names = append(names, name)
		}
	}
	return names
}

func parseSoftIRQDescriptions(interrupts *procfs.Interrupts) []string {
	var descriptions []string
	for _, interrupt := range interrupts.Interrupts {
		description := fmt.Sprintf("softirq.%s = %s", interrupt.Name, interrupt.Description)
		descriptions = append(descriptions, description)
	}
	return descriptions
}

func parseSoftIRQCounts(interrupts *procfs.Interrupts) []uint64 {
	var values []uint64
	for _, interrupt := range interrupts.Interrupts {
		for _, count := range interrupt.Counts {
			values = append(values, count)
		}
	}
	return values
}
