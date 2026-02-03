# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4](https://github.com/jonwraymond/toolcompose/compare/v0.1.3...v0.1.4) (2026-02-03)


### Documentation

* update version matrix ([#9](https://github.com/jonwraymond/toolcompose/issues/9)) ([ea6de6e](https://github.com/jonwraymond/toolcompose/commit/ea6de6eca36fb1dd96ea68cc83f2555ee3ac2355))

## [Unreleased]

### Added

#### set package
- Sentinel errors (`ErrNoSource`, `ErrNilAdapter`, `ErrNilToolset`) for programmatic error handling with `errors.Is()`
- Comprehensive package documentation in `doc.go` covering core concepts, usage patterns, thread safety, and error handling
- Contract documentation for `FilterFunc` type specifying thread-safety and mutation requirements
- Thread-safety documentation for `Builder` type clarifying single-goroutine usage
- Example tests for pkg.go.dev: `ExampleNewBuilder`, `ExampleNewBuilder_withNamespace`, `ExampleNewBuilder_withPolicy`, `ExampleToolset_Filter`, `ExampleAllowScopes`
- Benchmark tests: `BenchmarkToolset_Add`, `BenchmarkToolset_Get`, `BenchmarkToolset_Tools`, `BenchmarkToolset_Filter`, `BenchmarkBuilder_Build`, `BenchmarkBuilder_Build_WithFilters`, `BenchmarkBuilder_Build_WithPolicy`, `BenchmarkNamespaceFilter`, `BenchmarkTagsAll`

#### skill package
- Consolidated sentinel errors in `errors.go`: `ErrInvalidRunner`, `ErrInvalidSkillName`, `ErrInvalidStepID`, `ErrInvalidToolID`, `ErrNoSteps`, `ErrMaxStepsExceeded`, `ErrToolNotAllowed`
- Comprehensive package documentation in `doc.go` covering execution model, guards, determinism, and thread safety
- Determinism contract for `Guard` interface specifying idempotency requirements
- Example tests for pkg.go.dev: `ExampleNewPlanner`, `ExampleExecute`, `ExampleMaxStepsGuard`, `ExampleAllowedToolIDsGuard`, `ExampleSkill_Validate`
- Benchmark tests: `BenchmarkPlanner_Plan`, `BenchmarkExecute`, `BenchmarkSkill_Validate`, `BenchmarkMaxStepsGuard`, `BenchmarkAllowedToolIDsGuard`
- Additional test coverage: `TestExecute_NilRunner`, `TestExecute_EmptyPlan`, `TestExecute_PartialResultsOnError`, edge cases for guards and validation

### Changed

#### set package
- `Exposure.Export()` now returns `ErrNilAdapter` sentinel error instead of inline error string
- Builder execution flow documented in type comment

#### skill package
- Error definitions consolidated from scattered locations into single `errors.go` file
- Improved test coverage from 84.8% to 97.8%

### Fixed

- `exposure_test.go` updated to use sentinel error comparison instead of string matching

## [0.1.0] - Initial Release

### Added
- `set` package: Toolset composition, filtering, and protocol exposure
- `skill` package: Declarative skill planning and execution
- Builder pattern for constructing filtered toolsets
- Built-in filters: namespace, tags, scopes
- Policy system with DenyByDefault and AllowByDefault
- Guard system for skill execution constraints
- Integration with toolfoundation adapter system
