# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-02-03)


### Features

* expand toolcompose set/skill with examples and tests ([fb7ce0d](https://github.com/jonwraymond/toolcompose/commit/fb7ce0de923b2bcfd0ddf6dc674612373b583e15))
* initial repository structure ([7ae48ca](https://github.com/jonwraymond/toolcompose/commit/7ae48ca7a7984bbdf808a276f1f3f60824edebe1))
* **set:** migrate toolset to toolcompose/set package (PRD-150) ([4997ef4](https://github.com/jonwraymond/toolcompose/commit/4997ef448f61f92f1eab34a9c1cab25a041d3a10))
* **skill:** migrate toolskill to toolcompose/skill package (PRD-151) ([d643092](https://github.com/jonwraymond/toolcompose/commit/d64309212377aeaedd873bea87906aee2d143077))


### Documentation

* add mkdocs config ([dd9b6d2](https://github.com/jonwraymond/toolcompose/commit/dd9b6d2d9dbd0852c6f8d8bf9e214e8f1d502062))
* **toolcompose:** align set/skill documentation ([b96dd40](https://github.com/jonwraymond/toolcompose/commit/b96dd40bd6616a3cee4b03c4271a42d6241b0e74))

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
