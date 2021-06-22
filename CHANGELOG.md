# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [1.0.857] - 2021/04/28### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added### Added
- `upload` command processes file parts in parallel. Defaults to 25, but can be changed via flag `--max-concurrent-connections 50`

### Fix
- Reduce memory usage when not in debug mode.

## [1.0.856] - 2021/04/28### Added
### Added
- New flag `--max-concurrent-connections` to `download` and `upload` commands. (Default is 10)
- `download` only shows one progress bar when downloading a single file.

### Fix
- Uploading nested folders could sometimes skip folders.
- `files-cli download documents/report.pdf local-files` would result in `local-files/documents/report.pdf`. This is now fixed resulting in `local-files/report.pdf`
- `files-cli download documents/report.pdf local-files/report-2020.pdf` would result in `local-files/documents/report-2020.pdf/report.pdf`. This is now fixed resulting in `local-files/report-2020.pdf`
- Allows for downloaded of folders that contain more than 1000 files/folders.
- Fix file size attribute from failing on 32-bit releases.
- Downloading large files could hang once showing 100% due to inconsistencies in reported size.

## [1.0.669] - 2021/04/12
### Added
- New commands `login` and `logout`
- `config reset` now takes flags to reset a specific key.

### Fix
- In some cases API errors were not being returned correctly.
- `session delete` no longer returns an error.

## [1.0.215] - 2021/02/22
### Fix
- Windows command prompt for session login now formats input correctly.

## [1.0.210] - 2021/02/19
### Fix
- version command now displays the current version correctly.

## [1.0.196] - 2021/02/17
### Fix
- `download` and `upload` command now support session login

## [1.0.194] - 2021/02/16
### Added
- Support basic login and 2FA methods sms, u2f, yubi, and otp.
- Command `config` holds configuration in `~/.config/files-cli`
  - subcommand `set` with flags `subdomain`, `username`, and `api-key`
  - subcommand `reset` deletes the config to start from fresh state.
