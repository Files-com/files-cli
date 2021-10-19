# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [1.2.0] - 2021/10/19
### Add
- `sync { push | pull } [flags]`
  - Post sync actions:
    - `--local-path`, `--remote-path`
    - `--move-source {path}`
    - `--delete-source`

### Fix
- `sync` compares file size instead of modified time to match the server sync.
- Fixes uploading errors to some remote servers for files over 1GB.

## [1.1.1691] - 2021/09/30
### Add
- An upload that fails in the middle will be retried at the point it failed. If there are multiple files the failed upload will be retried after all other files have finished.

## [1.1.1679] - 2021/09/23
### Fix
- Address `upload` & `download` commands not displaying progress.

## [1.1.1677] - 2021/09/22
### Change
- Remove all single letter flags for generated commands

### Fix
- Missing flag for `bundles create` `--paths`

## [1.1.1666] - 2021/09/14
### Change
- `upload`/`download` better handles transfers with many files.
  - Improved performance with a new scanning state while the total number of files is not known.
  - No longer displays a progress bar for each file, but shows a single overall progress bar with a status line below.

### Add
- `upload`/`download`
  - `ignore` flag ignore individual files or match by patterns. See https://git-scm.com/docs/gitignore#_pattern_format
- When listing results and using default `--format table` results paginated, similar to `less`. To advance hit enter/return.

### Fix
- `fields` flag
  - fields are now case-insensitive.
  - Invalid fields names return an error
- Errors now are sent to stderr rather than stdout.
- `format` flag with `table` returns file size in a human readable format.

## [1.1.1645] - 2021/08/11
### Fix
- Rendering output in the default table format was printed twice.
- Fix documentation for `format` flag `table-light` to `table-bright`

## [1.1.1644] - 2021/08/11
### Fix
- `upload`/`download` flag `send-logs-to-cloud` has now been fixed.

## [1.1.1588] - 2021/08/03
### Add
- Flag `format` default: `table`, options: "json, csv, table, table-dark, table-bright"

### Change
- Default output format was `json` now it's `table`

## [1.1.1585] - 2021/08/03### Add
### Add
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
