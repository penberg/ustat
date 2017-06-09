package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type interruptsStatReader struct {
	counts []uint64
}

const procInterruptsPath = "/proc/interrupts"

// NewInterruptsStat returns a new Stat, which collects interrupt stats from /proc/interrupts.
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
		Reader:       &interruptsStatReader{counts: counts},
	}
}

func (reader *interruptsStatReader) Read() []uint64 {
	interrupts, err := procfs.ReadInterrupts(procInterruptsPath)
	if err != nil {
		panic(err)
	}
	counts := parseInterruptCounts(interrupts)
	diff := Difference(reader.counts, counts)
	reader.counts = counts
	return diff
}

func parseInterruptNames(interrupts *procfs.Interrupts) []string {
	var names []string
	for _, interrupt := range interrupts.Interrupts {
		for cpu := range interrupt.Counts {
			name := fmt.Sprintf("int%s.cpu%d", interrupt.Name, cpu)
			names = append(names, name)
		}
	}
	return names
}

func parseInterruptDescriptions(interrupts *procfs.Interrupts) []string {
	var descriptions []string
	for _, interrupt := range interrupts.Interrupts {
		description := fmt.Sprintf("intr.%s = %s", interrupt.Name, interrupt.Description)
		descriptions = append(descriptions, description)
	}
	return descriptions
}

func parseInterruptCounts(interrupts *procfs.Interrupts) []uint64 {
	var values []uint64
	for _, interrupt := range interrupts.Interrupts {
		for _, count := range interrupt.Counts {
			values = append(values, count)
		}
	}
	return values
}
