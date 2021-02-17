# Change Log

All notable changes to this project will be documented in this file.

## [1.0.196] - 2021/02/17
### Fix
- version command now display current version correctly.

## [1.0.195] - 2021/02/17
### Fix
- `download` and `upload` command now support session login

## [1.0.194] - 2021/02/16
### Added
- Support basic login and 2FA methods sms, u2f, yubi, and otp.
- Command `config` holds configuration in `~/.config/files-cli`
  - subcommand `set` with flags `subdomain`, `username`, and `api-key`
  - subcommand `reset` deletes the config to start from fresh state.
