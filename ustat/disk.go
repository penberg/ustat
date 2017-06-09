package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type DiskStatReader struct {
	values []uint64
}

const procDiskStatPath = "/proc/diskstats"

func NewDiskStat() *Stat {
	stats, err := procfs.ReadDiskStats(procDiskStatPath)
	if err != nil {
		panic(err)
	}
	names := parseDiskStatNames(stats)
	descriptions := parseDiskStatDescriptions(stats)
	values := parseDiskStats(stats)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Reader:       &DiskStatReader{values: values},
	}
}

func (reader *DiskStatReader) Read() []uint64 {
	stats, err := procfs.ReadDiskStats(procDiskStatPath)
	if err != nil {
		panic(err)
	}
	values := parseDiskStats(stats)
	diff := make([]uint64, 0)
	for idx := range reader.values {
		diff = append(diff, values[idx]-reader.values[idx])
	}
	reader.values = values
	return diff
}

var diskStatTypes = []string{
	"read.sectors",
	"write.sectors",
}

var diskStatDescriptions = map[string]string{
	"read.sectors":  "Number of 512 byte sectors read",
	"write.sectors": "Number of 512 byte sectors written",
}

func parseDiskStatNames(stats []procfs.DiskStat) []string {
	names := make([]string, 0)
	for _, stat := range stats {
		for _, diskStatType := range diskStatTypes {
			name := fmt.Sprintf("disk.%s.%s", stat.Name, diskStatType)
			names = append(names, name)
		}
	}
	return names
}

func parseDiskStatDescriptions(stats []procfs.DiskStat) []string {
	descriptions := make([]string, 0)
	for _, stat := range stats {
		for _, diskStatType := range diskStatTypes {
			diskStatDescription := diskStatDescriptions[diskStatType]
			description := fmt.Sprintf("disk.%s.%s = %s %s", stat.Name, diskStatType, stat.Name, diskStatDescription)
			descriptions = append(descriptions, description)
		}
	}
	return descriptions
}

func parseDiskStats(stats []procfs.DiskStat) []uint64 {
	values := make([]uint64, 0)
	for _, stat := range stats {
		values = append(values, stat.ReadSectors)
		values = append(values, stat.WriteSectors)
	}
	return values
}
