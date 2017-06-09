# ustat

`ustat` is an unified system stats collector for Linux, which combines capabilities of tools like `vmstat`, `mpstat`, `iostat`, and `ifstat`.
The tool is designed for low collection overhead to make it suitable for stats collection when evaluating system performance under load.
The main objective of `ustat` is to collect detailed stats rather than aggregate stats so that it is possible to drill down to details during analysis.
The `ustat` tool reports collected stats in a self-describing, [delimiter-separated values](https://en.wikipedia.org/wiki/Delimiter-separated_values) (DSV) format file that is easy to post process using tools like [ggplot2](http://ggplot2.org/) for R and [gnuplot](http://www.gnuplot.info/).

## Install

```sh
go get -u github.com/penberg/ustat/cmd/ustat
```

## Usage

To run `ustat`:

```sh
ustat 1
```

In the above example, `ustat` collects all stats it supports and samples them every one second.

Please use the `ustat --help` command for more information on supported stats collectors and other command line options.

## Related Tools

* [dstat](http://dag.wiee.rs/home-made/dstat/) - Versatile resource statistics tool. The tool provides similar capabilities as `ustat` but is written in [Python](https://www.python.org/), which has higher collection overhead, and does not provide detailed stats for everything (e.g. interrupts).

## Authors

* [Pekka Enberg](https://penberg.github.io/)

See also the list of [contributors](https://github.com/penberg/ustat/contributors) who participated in this project.

## License

`ustat` is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
