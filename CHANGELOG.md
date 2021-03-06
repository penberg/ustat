# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## [0.2.0] - 2017-07-13
### Added
- Collect aggregate CPU stats.
- Collect SoftIRQ stats.
- Add `ustat report` command for summarizing results.

### Changed
- Move stats collection under `ustat record` command.

### Fixed
- Fix CPU utilization percentage calculation.

## 0.1.0 - 2017-06-08
### Added
- A `ustat` command line tool written in Go, similar to `dstat`, for example.
- CPU utilization stats per core and number of context switches, which are collected from `/proc/stat`.
- Interrupt count stats per core, which are collected from `/proc/interrupts`.
- Network stats per interface, which are collected from `/proc/net/dev`.
- Disk stats per block device, which are collected from `/proc/diskstats`.

[Unreleased]: https://github.com/penberg/ustat/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/penberg/ustat/compare/v0.1.0...v0.2.0
