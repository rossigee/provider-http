# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.1] - 2026-09-24

### Changed
- **Dependencies**: Refreshed Go, Crossplane, and provider dependencies for the v2.5.0 release baseline.
- **Release publishing**: Restricted publishing to exact SemVer tags at the current `origin/master`, built and published `linux_amd64` and `linux_arm64` xpkg files, aliased `latest`, verified equal digests and both architectures, and used GitHub token authentication without OIDC claims.
