package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type InterruptsStatReader struct {
	counts []uint64
}

const procInterruptsPath = "/proc/interrupts"

func NewInterruptsStat() *Stat {
	interrupts, err := procfs.ReadInterrupts(procInterruptsPath)
	if err != nil {
		panic(err)
	}
	names := parseInterruptNames(interrupts)
	descriptions := parseInterruptDescriptions(interrupts)
	counts := parseInterruptCounts(interrupts)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Reader:       &InterruptsStatReader{counts: counts},
	}
}

func (reader *InterruptsStatReader) Read() []uint64 {
	interrupts, err := procfs.ReadInterrupts(procInterruptsPath)
	if err != nil {
		panic(err)
	}
	counts := parseInterruptCounts(interrupts)
	diff := make([]uint64, 0)
	for idx := range reader.counts {
		diff = append(diff, counts[idx]-reader.counts[idx])
	}
	reader.counts = counts
	return diff
}

func parseInterruptNames(interrupts *procfs.Interrupts) []string {
	names := make([]string, 0)
	for _, interrupt := range interrupts.Interrupts {
		for cpu := range interrupt.Counts {
			name := fmt.Sprintf("int%s.cpu%d", interrupt.Name, cpu)
			names = append(names, name)
		}
	}
	return names
}

func parseInterruptDescriptions(interrupts *procfs.Interrupts) []string {
	descriptions := make([]string, 0)
	for _, interrupt := range interrupts.Interrupts {
		description := fmt.Sprintf("intr.%s = %s", interrupt.Name, interrupt.Description)
		descriptions = append(descriptions, description)
	}
	return descriptions
}

func parseInterruptCounts(interrupts *procfs.Interrupts) []uint64 {
	values := make([]uint64, 0)
	for _, interrupt := range interrupts.Interrupts {
		for _, count := range interrupt.Counts {
			values = append(values, count)
		}
	}
	return values
}
