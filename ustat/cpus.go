package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type CpusStatReader struct {
	values []uint64
}

const procStatPath = "/proc/stat"

func NewCpusStat() *Stat {
	stat, err := procfs.ReadStat(procStatPath)
	if err != nil {
		panic(err)
	}
	names := parseCpuStatNames(stat)
	descriptions := parseCpuStatDescriptions(stat)
	values := parseCpuStats(stat)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Reader:       &CpusStatReader{values: values},
	}
	return nil
}

func (reader *CpusStatReader) Read() []uint64 {
	stat, err := procfs.ReadStat(procStatPath)
	if err != nil {
		panic(err)
	}
	values := parseCpuStats(stat)
	diff := make([]uint64, 0)
	for idx := range reader.values {
		diff = append(diff, values[idx]-reader.values[idx])
	}
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

func parseCpuStatNames(stat *procfs.Stat) []string {
	names := make([]string, 0)
	for _, cpuStat := range stat.CPUStats {
		for _, cpuStatType := range cpuStatTypes {
			name := fmt.Sprintf("%s.%s", cpuStat.Id, cpuStatType)
			names = append(names, name)
		}
	}
	return names
}

func parseCpuStatDescriptions(stat *procfs.Stat) []string {
	descriptions := make([]string, 0)
	for _, cpuStat := range stat.CPUStats {
		for _, cpuStatType := range cpuStatTypes {
			cpuStatDescription := cpuStatDescriptions[cpuStatType]
			description := fmt.Sprintf("%s.%s = %s %s", cpuStat.Id, cpuStatType, cpuStat.Id, cpuStatDescription)
			descriptions = append(descriptions, description)
		}
	}
	return descriptions
}

func parseCpuStats(stat *procfs.Stat) []uint64 {
	values := make([]uint64, 0)
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
	return values
}
