# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.1.2] - 2018-11-22

### Added

- add logs for when server shutdown errors

## [1.1.1] - 2018-11-16

### Added

- use go modules
- Integrated the Degraded status
- Set up a global Timeout to run the health check tests in

### Changed

- Renamed Result to Status and renamed the values to clearer verbs
- Changed the interface of the TestFunc to include the Status
- Tests will now run concurrently instead of sequentially

## [1.1.0] - 2018-01-03

### Added

- Cache Middleware

## [1.0.0] - 2017-12-05

### Added

- Renamed the package to `healthz`

### Changed

- `Server.Start()` doesn't take a `context.Context` anymore.

## [0.0.1] - 2017-11-27

### Added

- Allow registering health check tests
- Add hook to attach the `/_healthz` endpoint to an existing set of handlers
- Add health check server
