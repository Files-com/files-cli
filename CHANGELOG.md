# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [1.1.1588] - 2021/08/03
### Fix
- `upload`/`download` flag `send-logs-to-cloud` has now been fixed.

## [1.1.1585] - 2021/08/03
### Add
- Flag `format` default: `table`, options: "json, csv, table, table-dark, table-light"

### Change
- Default output format was `json` now it's `table`

## [1.1.1584] - 2021/07/28### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add### Add
- `upload` and `download`
    - Flag `sync` to only transfer files based on the modified date.
    - Flag `send-logs-to-cloud` to sends output as external log.
    - Flag `disable-progress-output` to disable progress bars and only show status when file is complete.

### Change
- `upload` and `download`
    - Flag `max-concurrent-connections` removed and replaced with `concurrent-file-uploads` and `concurrent-file-downloads`

### Fix
- Incorrectly typed subcommands return error rather than no output.
- `config reset` now resets all persistent config rather that having to reset by flags.

## [1.0.1318] - 2021/06/29
### Fix
- `files-cli download` Fix Windows issue `The process cannot access the file because it is being used by another process.`
- `files-cli download` in some cases the CLI hangs after all files are download.

- String flags that only accept a list of values are validated before sending request. This list of values is also listed in the usage description.
### Added
- String flags that only accept a list of values are validated before sending request. This list of values is also listed in the usage description.

## [1.0.857] - 2021/04/28
### Added
- `upload` command processes file parts in parallel. Defaults to 25, but can be changed via flag `--max-concurrent-connections 50`

### Fix
- Reduce memory usage when not in debug mode.

## [1.0.856] - 2021/04/28
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
