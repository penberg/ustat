package ustat

import (
	"fmt"
	procfs "github.com/c9s/goprocinfo/linux"
)

type procNetDevCollector struct {
	values []uint64
}

const procNetDevPath = "/proc/net/dev"

// NewNetStat returns a new Stat, which collects networking stats from /proc/net/dev.
func NewNetStat() *Stat {
	stats, err := procfs.ReadNetworkStat(procNetDevPath)
	if err != nil {
		panic(err)
	}
	names := parseNetStatNames(stats)
	descriptions := parseNetStatDescriptions(stats)
	values := parseNetStats(stats)
	return &Stat{
		Names:        names,
		Descriptions: descriptions,
		Collector:    &procNetDevCollector{values: values},
	}
}

func (reader *procNetDevCollector) Collect() []uint64 {
	stats, err := procfs.ReadNetworkStat(procNetDevPath)
	if err != nil {
		panic(err)
	}
	values := parseNetStats(stats)
	diff := Difference(reader.values, values)
	reader.values = values
	return diff
}

var netStatTypes = []string{
	"rx.bytes",
	"rx.packets",
	"rx.errors",
	"rx.drop",
	"tx.bytes",
	"tx.packets",
	"tx.errors",
	"tx.drop",
}

var netStatDescriptions = map[string]string{
	"rx.bytes":   "Number of bytes received",
	"rx.packets": "Number of packets received",
	"rx.errors":  "Number of receive errors",
	"rx.drop":    "Number of receive packets dropped",
	"tx.bytes":   "Number of bytes transmitd",
	"tx.packets": "Number of packets transmitd",
	"tx.errors":  "Number of transmit errors",
	"tx.drop":    "Number of transmit packets dropped",
}

func parseNetStatNames(stats []procfs.NetworkStat) []string {
	var names []string
	for _, stat := range stats {
		for _, netStatType := range netStatTypes {
			name := fmt.Sprintf("net.%s.%s", stat.Iface, netStatType)
			names = append(names, name)
		}
	}
	return names
}

func parseNetStatDescriptions(stats []procfs.NetworkStat) []string {
	var descriptions []string
	for _, stat := range stats {
		for _, netStatType := range netStatTypes {
			netStatDescription := netStatDescriptions[netStatType]
			description := fmt.Sprintf("net.%s.%s = %s %s", stat.Iface, netStatType, stat.Iface, netStatDescription)
			descriptions = append(descriptions, description)
		}
	}
	return descriptions
}

func parseNetStats(stats []procfs.NetworkStat) []uint64 {
	var values []uint64
	for _, stat := range stats {
		values = append(values, stat.RxBytes)
		values = append(values, stat.RxPackets)
		values = append(values, stat.RxErrs)
		values = append(values, stat.RxDrop)
		values = append(values, stat.TxBytes)
		values = append(values, stat.TxPackets)
		values = append(values, stat.TxErrs)
		values = append(values, stat.TxDrop)
	}
	return values
}
