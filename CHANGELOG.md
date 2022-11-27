# Changelog

All notable changes to this project will be documented in this file.

## 1.2.2

### Fixed

- upgrade dependencies (this is a maintenance release)

## 1.2.1

### Changed

- Using withFields from logrus to pass the set of parameters

### Removed

- Colorful option (because Logrus can't really handle it, especially when logs are redirected to a file)

### Fixed

- Duplicate information in trace logs (duration and lines affected)

## 1.2.0

### Added

- Gorm options support

### Fixed

- upgrade dependencies for maintenance

## 1.1.2

### Fixed

- upgrade dependencies (this is a maintenance release)

## 1.1.1

### Fixed

- upgrade dependencies (this is a maintenance release)

## 1.1.0

### Added

- Missing log level management
- "SLOW SQL" Prefix on slow sql logs

## 1.0.0

### Added

- Gormlog library
- Example in readme
