package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type procStatCollector struct {
	prev *procfs.Stat
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
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Collector: &procStatCollector{
			prev: stat,
		},
	}
	return nil
}

func interval(before []uint64, after []uint64) uint64 {
	var prev uint64 = 0
	for _, value := range before {
		prev += value
	}
	var curr uint64 = 0
	for _, value := range after {
		curr += value
	}
	return curr - prev
}

func (reader *procStatCollector) Collect() []uint64 {
	stat, err := procfs.ReadStat(procStatPath)
	if err != nil {
		panic(err)
	}
	values := parseCPUStats(stat, reader.prev)
	reader.prev = stat
	return values
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
	for _, cpuStatType := range cpuStatTypes {
		name := fmt.Sprintf("%s.%s", "cpu", cpuStatType)
		names = append(names, name)
	}
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

func parseCPUStats(curr *procfs.Stat, prev *procfs.Stat) []uint64 {
	var values []uint64
	interval := runtime(curr.CPUStatAll) - runtime(prev.CPUStatAll)
	values = append(values, difference(curr.CPUStatAll.User, prev.CPUStatAll.User, interval))
	values = append(values, difference(curr.CPUStatAll.Nice, prev.CPUStatAll.Nice, interval))
	values = append(values, difference(curr.CPUStatAll.System, prev.CPUStatAll.System, interval))
	values = append(values, difference(curr.CPUStatAll.Idle, prev.CPUStatAll.Idle, interval))
	values = append(values, difference(curr.CPUStatAll.IOWait, prev.CPUStatAll.IOWait, interval))
	values = append(values, difference(curr.CPUStatAll.IRQ, prev.CPUStatAll.IRQ, interval))
	values = append(values, difference(curr.CPUStatAll.SoftIRQ, prev.CPUStatAll.SoftIRQ, interval))
	values = append(values, difference(curr.CPUStatAll.Steal, prev.CPUStatAll.Steal, interval))
	values = append(values, difference(curr.CPUStatAll.Guest, prev.CPUStatAll.Guest, interval))
	values = append(values, difference(curr.CPUStatAll.GuestNice, prev.CPUStatAll.GuestNice, interval))
	for idx, currCpuStat := range curr.CPUStats {
		prevCpuStat := prev.CPUStats[idx]
		interval := runtime(currCpuStat) - runtime(prevCpuStat)
		values = append(values, difference(currCpuStat.User, prevCpuStat.User, interval))
		values = append(values, difference(currCpuStat.Nice, prevCpuStat.Nice, interval))
		values = append(values, difference(currCpuStat.System, prevCpuStat.System, interval))
		values = append(values, difference(currCpuStat.Idle, prevCpuStat.Idle, interval))
		values = append(values, difference(currCpuStat.IOWait, prevCpuStat.IOWait, interval))
		values = append(values, difference(currCpuStat.IRQ, prevCpuStat.IRQ, interval))
		values = append(values, difference(currCpuStat.SoftIRQ, prevCpuStat.SoftIRQ, interval))
		values = append(values, difference(currCpuStat.Steal, prevCpuStat.Steal, interval))
		values = append(values, difference(currCpuStat.Guest, prevCpuStat.Guest, interval))
		values = append(values, difference(currCpuStat.GuestNice, prevCpuStat.GuestNice, interval))
	}
	values = append(values, curr.ContextSwitches-prev.ContextSwitches)
	return values
}

func runtime(cpuStat procfs.CPUStat) uint64 {
	return cpuStat.User + cpuStat.Nice + cpuStat.System + cpuStat.Idle + cpuStat.IOWait + cpuStat.IRQ + cpuStat.SoftIRQ
}

func difference(curr uint64, prev uint64, interval uint64) uint64 {
	return uint64(float64(curr-prev) / float64(interval) * 100)
}
