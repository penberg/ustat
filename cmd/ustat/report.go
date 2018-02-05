package main

import (
	"encoding/csv"
	"fmt"
	"github.com/montanaflynn/stats"
	"gopkg.in/urfave/cli.v1"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type cpuStat struct {
	values map[string][]float64
}

type interruptStat struct {
	values map[string][]float64
}

var reportCommand = cli.Command{
	Name:      "report",
	Usage:     "summarise stats that are recored to a file",
	ArgsUsage: "[file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "delimiter",
			Usage: "delimiter used in the output file",
			Value: "\t",
		},
	},
	Action: reportAction,
}

func reportAction(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) == 0 {
		return cli.NewExitError(fmt.Sprintf("No stats file specified"), 3)
	}
	filename := args[0]
	delimiter := ctx.String("delimiter")
	file, err := os.Open(filename)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to open file: %v", err), 2)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	comma, _ := utf8.DecodeRuneInString(delimiter)
	reader.Comma = comma
	header := map[string]int{}
	cpuStats := map[string]cpuStat{}
	interruptStats := map[string]interruptStat{}
	softIrqStats := map[string]interruptStat{}
	sampleCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(record) == 0 {
			continue
		}
		if strings.HasPrefix(record[0], "#") {
			continue
		}
		if len(header) == 0 {
			for idx, column := range record {
				header[column] = idx
			}
			continue
		}
		for column, idx := range header {
			result := strings.Split(column, ".")
			resource := result[0]
			if strings.HasPrefix(resource, "cpu") {
				rawValue := record[idx]
				value, err := strconv.Atoi(rawValue)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Unable to parse value '%s': %v", rawValue, err), 2)
				}
				class := result[1]
				stat, ok := cpuStats[resource]
				if !ok {
					stat = cpuStat{values: map[string][]float64{}}
				}
				values := stat.values[class]
				values = append(values, float64(value))
				stat.values[class] = values
				cpuStats[resource] = stat
			}
			if strings.HasPrefix(resource, "int") {
				rawValue := record[idx]
				value, err := strconv.Atoi(rawValue)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Unable to parse value '%s': %v", rawValue, err), 2)
				}
				class := result[1]
				stat, ok := interruptStats[resource]
				if !ok {
					stat = interruptStat{values: map[string][]float64{}}
				}
				values := stat.values[class]
				values = append(values, float64(value))
				stat.values[class] = values
				interruptStats[resource] = stat
			}
			if strings.HasPrefix(resource, "softirq") {
				resource := result[1]
				rawValue := record[idx]
				value, err := strconv.Atoi(rawValue)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Unable to parse value '%s': %v", rawValue, err), 2)
				}
				class := result[2]
				stat, ok := softIrqStats[resource]
				if !ok {
					stat = interruptStat{values: map[string][]float64{}}
				}
				values := stat.values[class]
				values = append(values, float64(value))
				stat.values[class] = values
				softIrqStats[resource] = stat
			}
		}
		sampleCount++
	}
	fmt.Printf("Processing %s ...\n", filename)
	fmt.Printf("\n")
	fmt.Printf("N = %d\n", sampleCount)
	cpus := []string{}
	for cpu, _ := range cpuStats {
		cpus = append(cpus, cpu)
	}
	sort.Strings(cpus)
	classes := []string{"system", "usr", "nice", "irq", "softirq", "iowait", "guest", "guestnice", "steal", "idle"}
	fmt.Printf("\n")
	fmt.Printf("CPU utilization, mean (SD):\n")
	fmt.Printf("\n")
	fmt.Printf("      ")
	for _, class := range classes {
		fmt.Printf(" %-12s", class)
	}
	fmt.Printf("\n")
	for _, cpu := range cpus {
		fmt.Printf("  %-4s", cpu)
		cpuStats := cpuStats[cpu]
		for _, class := range classes {
			values := cpuStats.values[class]
			mean, err := stats.Mean(values)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("%v", err), 2)
			}
			stddev, err := stats.StandardDeviation(values)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("%v", err), 2)
			}
			format := fmt.Sprintf("%.2f (%.2f)", mean, stddev)
			fmt.Printf(" %-12s", format)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
	if err := printInterrupts("Interrupts", interruptStats, cpus); err != nil {
		return err
	}
	fmt.Printf("\n")
	if err := printInterrupts("SoftIRQs", softIrqStats, cpus); err != nil {
		return err
	}
	fmt.Printf("\n")
	return nil
}

func printInterrupts(title string, interruptStats map[string]interruptStat, cpus []string) error {
	interrupts := []string{}
	for interrupt, _ := range interruptStats {
		interrupts = append(interrupts, interrupt)
	}
	sort.Strings(interrupts)
	fmt.Printf("%s, mean (SD):\n", title)
	fmt.Printf("\n")
	fmt.Printf("  %-10s", "interrupt")
	for _, cpu := range cpus {
		if cpu == "cpu" {
			continue
		}
		fmt.Printf("%20s", cpu)
	}
	fmt.Printf("\n")
	for _, intr := range interrupts {
		interruptStat := interruptStats[intr]
		fmt.Printf("  %-10s", intr)
		for _, cpu := range cpus {
			if cpu == "cpu" {
				continue
			}
			values, ok := interruptStat.values[cpu]
			if !ok {
				fmt.Printf("%20s", "")
				continue
			}
			mean, err := stats.Mean(values)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("%v", err), 2)
			}
			stddev, err := stats.StandardDeviation(values)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("%v", err), 2)
			}
			format := fmt.Sprintf("%.2f (%.2f)", mean, stddev)
			fmt.Printf("%20s", format)
		}
		fmt.Printf("\n")
	}
	return nil
}
