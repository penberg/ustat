package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type cpusStatReader struct {
	values []uint64
}

const procStatPath = "/proc/stat"

// NewCPUsStat returns a new Stat, which collects CPU stats from /proc/stat.
func NewCPUsStat() *Stat {
	stat, err := procfs.ReadStat(procStatPath)
	if err != nil {
		panic(err)
	}
	names := parseCPUStatNames(stat)
	descriptions := parseCPUStatDescriptions(stat)
	values := parseCPUStats(stat)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Reader:       &cpusStatReader{values: values},
	}
	return nil
}

func (reader *cpusStatReader) Read() []uint64 {
	stat, err := procfs.ReadStat(procStatPath)
	if err != nil {
		panic(err)
	}
	values := parseCPUStats(stat)
	diff := Difference(reader.values, values)
	reader.values = values
	return diff
}

var cpuStatTypes = []string{
	"usr",
	"nice",
	"system",
	"idle",
	"iowait",
	"irq",
	"softirq",
	"steal",
	"guest",
	"guestnice",
}

var cpuStatDescriptions = map[string]string{
	"usr":       "User",
	"nice":      "Nice",
	"system":    "System",
	"idle":      "Idle",
	"iowait":    "IOWait",
	"irq":       "IRQ",
	"softirq":   "SoftIRQ",
	"steal":     "Steal",
	"guest":     "Guest",
	"guestnice": "GuestNice",
}

func parseCPUStatNames(stat *procfs.Stat) []string {
	var names []string
	for _, cpuStat := range stat.CPUStats {
		for _, cpuStatType := range cpuStatTypes {
			name := fmt.Sprintf("%s.%s", cpuStat.Id, cpuStatType)
			names = append(names, name)
		}
	}
	names = append(names, "ctxt.switch")
	return names
}

func parseCPUStatDescriptions(stat *procfs.Stat) []string {
	var descriptions []string
	for _, cpuStat := range stat.CPUStats {
		for _, cpuStatType := range cpuStatTypes {
			cpuStatDescription := cpuStatDescriptions[cpuStatType]
			description := fmt.Sprintf("%s.%s = %s %s", cpuStat.Id, cpuStatType, cpuStat.Id, cpuStatDescription)
			descriptions = append(descriptions, description)
		}
	}
	descriptions = append(descriptions, "ctx.switch = Number of context switches")
	return descriptions
}

func parseCPUStats(stat *procfs.Stat) []uint64 {
	var values []uint64
	for _, cpuStat := range stat.CPUStats {
		values = append(values, cpuStat.User)
		values = append(values, cpuStat.Nice)
		values = append(values, cpuStat.System)
		values = append(values, cpuStat.Idle)
		values = append(values, cpuStat.IOWait)
		values = append(values, cpuStat.IRQ)
		values = append(values, cpuStat.SoftIRQ)
		values = append(values, cpuStat.Steal)
		values = append(values, cpuStat.Guest)
		values = append(values, cpuStat.GuestNice)
	}
	values = append(values, stat.ContextSwitches)
	return values
}
