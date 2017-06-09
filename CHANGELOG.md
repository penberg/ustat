# Changelog

## [Unreleased]
### Added
- A `ustat` command line tool written in Go, similar to `dstat`, for example.
- CPU utilization stats per core, which are collected from `/proc/stat`.
- Interrupt count stats per core, which are collected from `/proc/interrupts`.
- Network stats per interface, which are collected from `/proc/net/dev`.
- Disk stats per block device, which are collected from `/proc/diskstats`.

[Unreleased]: https://github.com/penberg/ustat/compare/8322e9b...HEAD
