package main

import (
	"fmt"
	"github.com/penberg/ustat"
	"gopkg.in/urfave/cli.v1"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "ustat"
	app.Usage = "Unified system statistics collector"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "c,cpu",
			Usage: "Enable CPU stats collection",
		},
		cli.BoolFlag{
			Name:  "i,int",
			Usage: "Enable interrupt stats collection",
		},
		cli.BoolFlag{
			Name:  "n,net",
			Usage: "Enable network stats collection",
		},
		cli.BoolFlag{
			Name:  "d,disk",
			Usage: "Enable disk stats collection",
		},
		cli.StringFlag{
			Name:  "o,output",
			Usage: "Write output to a file",
		},
		cli.StringFlag{
			Name:  "grep",
			Usage: "Regular expression patter for filtering stats",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		var stats []*ustat.Stat
		if ctx.Bool("cpu") {
			stats = append(stats, ustat.NewCPUsStat())
		}
		if ctx.Bool("int") {
			stats = append(stats, ustat.NewInterruptsStat())
		}
		if ctx.Bool("net") {
			stats = append(stats, ustat.NewNetStat())
		}
		if ctx.Bool("disk") {
			stats = append(stats, ustat.NewDiskStat())
		}
		output := os.Stdout
		outputPath := ctx.String("output")
		if outputPath != "" {
			file, err := os.Create(outputPath)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Unable to open file: %v", err), 2)
			}
			defer file.Close()
			output = file
		}
		pattern := ctx.String("grep")
		if len(stats) == 0 {
			return cli.NewExitError("No stats enabled", 1)
		}
		fmt.Fprintf(output, "# This file has been generated by ustat.\n")
		fmt.Fprintf(output, "#\n")
		fmt.Fprintf(output, "# Column descriptions:\n")
		for _, stat := range stats {
			for _, description := range stat.Descriptions {
				if pattern != "" {
					matched, err := regexp.MatchString(pattern, description)
					if err != nil || !matched {
						continue
					}
				}
				fmt.Fprintf(output, "# %s\n", description)
			}
		}
		for _, stat := range stats {
			for _, name := range stat.Names {
				if pattern != "" {
					matched, err := regexp.MatchString(pattern, name)
					if err != nil || !matched {
						continue
					}
				}
				fmt.Fprintf(output, "%s\t", name)
			}
		}
		fmt.Fprintln(output, "")
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			done <- true
		}()
		ticker := time.NewTicker(time.Second * 1)
		go func() {
			for _ = range ticker.C {
				for _, stat := range stats {
					values := stat.Collector.Collect()
					for idx, value := range values {
						if pattern != "" {
							matched, err := regexp.MatchString(pattern, stat.Names[idx])
							if err != nil || !matched {
								continue
							}
						}
						fmt.Fprintf(output, "%d\t", value)
					}
				}
				fmt.Fprintln(output, "")
			}
		}()
		<-done

		return nil
	}
	app.Run(os.Args)
}
