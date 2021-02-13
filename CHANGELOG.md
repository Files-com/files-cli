# Change Log

All notable changes to this project will be documented in this file.

## [1.0.191] - 2021/02/13
### Fix
- version command now display current version correctly.

## [1.0.190] - 2021/02/12
### Fix
- `download` and `upload` command now support session login


### Added
- Support basic login and 2FA methods sms, u2f, yubi, and otp.
- Command `config` holds configuration in `~/.config/files-cli`
  - subcommand `set` with flags `subdomain`, `username`, and `api-key`
  - subcommand `reset` deletes the config to start from fresh state.
