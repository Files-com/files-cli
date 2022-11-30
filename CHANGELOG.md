# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [2.1.36] - 2022/11/30
### Fix
- `files` `move`/`copy` not returning initial errors.
- Default table format now suppress empty columns.

## [2.1.35] - 2022/11/30
### Fix
- Windows table format would overflow it now clips extra columns to fit screen width.

## [2.1.34] - 2022/11/29
### Fix
- field `name` was not settable on resource `remote-mounts`

## [2.1.33] - 2022/11/29
### Fix
- On Windows when listing a resource with pager resulted in garbled text. Changed to use non unicode character for table.

## [2.1.32] - 2022/11/29
### Add
- New resource `remote-mounts` for managing Remote Mounts.

## [2.1.31] - 2022/11/29
### Add
- Protected resources could require reauthenication when using session based authentication. This adds the `reauthentication` flag, which prompts for the users password for protected resources. See doc https://developers.files.com/#reauthentication

## [2.1.25] - 2022/11/23
### Add
- flag enum/field values now convert dashed value to underscore values.
- flag enum can return a did you mean response if value is unknown. 

## [2.1.17] - 2022/11/18
### Fix
- On create of .config directory use permission `0755` instead of `0600` 

## [2.1.5] - 2022/11/10
### Fix
- In cases of listing resources, where the server returned unexpected results, it would not show beyond the first page of results.

## [2.1.2] - 2022/11/09
### Fix
- Better parsing of API errors.

## [2.1.0] - 2022/11/07
### Add
- Single resource results now removes empty columns from `table` format.
- Single resource results now display vertically.
- `format` flag now take these options: '{format} {style} {direction}' - formats: {json, csv, table}
    - table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
    - json-styles: {raw, pretty}

### Fix
- `debug` flag sends more detailed logs to file including generated curl commands.

## [2.0.8] - 2022/11/07
### Fix
- Listing a resource with format JSON was missing an ending bracket. [issue#2](https://github.com/Files-com/files-cli/issues/2)

## [2.0.6] - 2022/11/04
### Fix
- Uploading a zero byte file, with `upload`/`sync`, resulted in error `Upload Not Found`

## [2.0.5] - 2022/11/03
### Fix
- Remote paths for `upload`/`download`/`sync` that start with a slash are now normalized, fixing a possible panic.

## [2.0.4] - 2022/11/02
## Fix
- Validate enum flag values and return error if invalid.

## [1.6.5] - 2022/10/26
## Add
- `profile` flag - Can setup different profiles to many sites using an api key or session login.
  ```bash
  files-cli config set --api-key {API_KEY_GOES_HERE} --profile site1
  files-cli folders ls --profile site1
  
  files-cli login --profile site2
  files-cli folders ls --profile site2
  ```
  
- `agent` currently hidden and under beta. This is an on-perm remote server. Contact customer support before use. 
  
## Remove
- u2f login, removing support due to OS deprecation warnings and [Chrome Browser](https://developer.chrome.com/blog/deps-rems-95/#deprecate-u2f-api-cryptotoken) dropping support. Will migrate to [WebAuthn](https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API) in the future.

## [1.6.1] - 2022/10/03
### Fix
- Don't require valid api-key/session for `version` command.
- version checking did not properly recheck once version was upgraded. If having issue use reset using `config reset --version-check`

## [1.6.0] - 2022/10/03
### Add
- Check if running the latest version and if out of date provide upgrade instruction. Can be ignored with a global flag `ignore-version-check`.

## [1.5.10] - 2022/09/07
### Fix
- `sync` logs now remove `size 0 B` when status is skipped.

## [1.5.9] - 2022/09/07
### Fix
- `login` remove error messaging when login requires 2FA.

## [1.5.8] - 2022/09/07
### Fix
- Boolean flags now support setting to false `--flag=false`.

## [1.5.2] - 2022/07/20
### Fix
- `--format` `table-markdown` for non paginated requests.

## [1.5.1] - 2022/07/19
### Fix
- Resources with parse errors no longer require format `json` to recover and display results.

## [1.5.0] - 2022/07/19
### Add
- Use system $PAGER (ie. less or more) when returning any more than a single result. Can be disabled with `--use-pager=false`
- Add loading indicator for listing slow queries.
- `--format` `table-markdown`

## [1.4.3] - 2022/07/18
### Improvement
- Removes flag `--times` for `upload`/`sync push` this is now done by default.

## [1.4.2] - 2022/07/18
### Add
- Commands `list` and `list-*` with the default format of `table`
  - Prefetches the next page making it faster to display when next page is requested.
  - Shows loading indicator when pages on slow page loads.
- `folders list-for` now as alias `folders ls`
- Commands now show a descriptions when listing sub commands.

### Fix
- Resources with parse errors can now recover and display results when using format `json`.
- Better parsing of server response errors.
- Address resources where id is inserted into path.
- Listing resources with nested values now defaults to JSON instead of internal Map value.
- Resources that returned non paginated collection results no longer return parse error.
- `upload`/`download`/`sync` flag `--times` preserve modification times now falls back to mtime if no provided_time.

## [1.4.1] - 2022/07/18
### Fix
- Address an issue formatting non paginated result sets.

## [1.4.0] - 2022/07/17
### Add
- `upload`/`download`/`sync` new flag `--times` preserve modification times. Defaults to false.

## [1.3.70] - 2022/07/08
### Fix
- `upload`/`sync` when on an unstable connection was failing with `Your socket connection to the server was not read from or written to within the timeout period. Idle connections will be closed.`. This was fixed by property rewinding the file part before trying again.

## [1.3.67] - 2022/06/29
### Fix
- `users update --id XXX` would return `Datetime Parse - 'authenticate_until must contain valid date and time'`

## [1.3.65] - 2022/06/23
### Fix
- Upload to remote mounts could fail with an etag error.

## [1.3.61] - 2022/06/15
### Fix
- `* delete` would return blank attributes for an entity. It now should return nothing unless there is an error.

## [1.3.60] - 2022/06/14
### Fix
- `upload`/`download`/`sync`
  - when using progress bars now shows full path of file if there is enough width in the terminal.
  - Could delay files from transferring until scanning phase was complete.

## [1.3.58] - 2022/06/07
### Fix
- `sync` command could panic at the end of an operation.

## [1.3.10] - 2021/12/16
### Add
- `files` async commands `copy` and `move` 
  - flag `block`. This returns a progress bar of known status. To skip progress bars use flag `no-progress`. If task fails it will return a non-zero exit code.
  - flag `event-log`. After the operation is complete this returns a log line for each file. This can be formatted with the standard `format` flag.
  
## [1.3.8] - 2021/12/03
### Add
- `folders list-for` new flag `recursive` (list folders/files recursively)
- `folders`

## [1.3.0] - 2021/10/27
### Add
- `folders list-for` new flag `only-folders`

## [1.2.3] - 2021/10/25
### Fix
- `upload` or `sync push` could cause a panic error when uploading to a remote server.

## [1.2.2] - 2021/10/22
### Improvement
- 3x performance when uploading to remote mounts.
- Less jittery transfer rate and ETA indicators when uploading to slower remote mounts.

## [1.2.1] - 2021/10/19
### Fix
- Improved error handling for unexpected HTML errors.
- When uploading to a remote server files would incorrectly report 0 bytes transferred after completion.

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
