# Changelog

All changes to the GraphQL Test Tool (gtt) are documented here. Releases follow semantic versioning.

## [Unreleased]

## [1.7.3] - 2021-08-18
### Fixed
- Fixed incorrect error check.

## [1.7.2] - 2021-08-17
### Added
- Add support for JSONPath for remembered variables.

## [1.7.1] - 2021-07-09
### Added
- Add memory access for remembered variables.

## [1.7.0] - 2021-07-09
### Added
- More display options for pretty format and SEN (from OjG) formatting.

## [1.6.1] - 2021-03-24
### Changed
- Updated go.mod to use the latest OjG.

## [1.6.0] - 2021-02-02
### Added
- The OjG SEN format is now allowed for script files. This make
  reading the scripts very tolerant of JSON errors.

## [1.5.1] - 2020-10-08
### Fixed
- Timeout option fixed.

## [1.5.0] - 2020-09-11
### Added
- Allow variable replacement in headers and content with `$(var)` pattern.

## [1.4.0] - 2020-09-10
### Added
- Options to set the expected response status code.

## [1.3.0] - 2020-09-02
### Added
- Request timeout added.
- Headers can now be set on requests.

## [1.2.1] - 2020-01-29
### Fixed
- Fixed display of use case filename when longer than 80 characters.

## [1.2.0] - 2020-01-16
### Added
- Exact matches on element in a map can be achieved with use of `"*": null` indicating
  any element not explicitly called out must be either `null` or not present.

## [1.1.0] - 2020-01-09
### Added
- Include file option added.
- Step option `always` added which indicates the step should always be run even if the test has failed.

## [1.0.11] - 2020-01-08
### Added
- Regular expressions can now be used to match non-string responses.

## [1.0.10] - 2020-01-02
### Added
- Enable line-wise regex match when expecting strings

## [1.0.9] - 2019-12-31

Sort nested arrays
### Fixed
- Nested arrays are now sorted when specified with a `sortBy` step option.

## [1.0.8] - 2019-12-18
Add dockerfile to repository
### Added- Dockerfile

## [1.0.7] - 2019-12-18
Array limits bug fix (v1.0.6 wasn't picked up by go mod, bug in sum.golang.org)
### Fixed
- Check for array comparisons now check for length.

## [1.0.5] - 2019-12-17
Test Reporting Improvements
### Added
- Displays mismatched values as part of error message when actual is not the same as expected.

## [1.0.2] - 2019-12-16
Sorting
### Fixed
- Sort before displaying results instead of after.

## [1.0.1] - 2019-12-16
Reorganize
### Fixed
- Folder and file layout so go.mod works as expected.
- Fixed support for comparing text responses.

## [1.0.0] - 2019-12-15

Iniital release
