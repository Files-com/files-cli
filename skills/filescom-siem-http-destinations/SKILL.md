---
name: filescom-siem-http-destinations
description: |
  Files.com Siem Http Destinations via files-cli.
---

# filescom-siem-http-destinations

Files.com Siem Http Destinations management via files-cli.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli siem-http-destinations list`

List SIEM HTTP Destinations.

No flags beyond the global ones.

### `files-cli siem-http-destinations find`

Show SIEM HTTP Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Siem Http Destination ID. **Required.** |

### `files-cli siem-http-destinations create`

Create SIEM HTTP Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name for this Destination |
| `--additional-headers` | object | Additional HTTP Headers included in calls to the destination URL |
| `--sending-active` | bool | Whether this SIEM HTTP Destination is currently being sent to or not |
| `--generic-payload-type` | enum | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. One of: `json_newline`, `json_array`. |
| `--file-destination-path` | string | Applicable only for destination type: file. Destination folder path on Files.com. |
| `--file-format` | enum | Applicable only for destination type: file. Generated file format. One of: `json`, `csv`. |
| `--file-interval-minutes` | int64 | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `--splunk-token` | string | Applicable only for destination type: splunk. Authentication token provided by Splunk. |
| `--azure-dcr-immutable-id` | string | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `--azure-stream-name` | string | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `--azure-oauth-client-credentials-tenant-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `--azure-oauth-client-credentials-client-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `--azure-oauth-client-credentials-client-secret` | string | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `--qradar-username` | string | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `--qradar-password` | string | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `--solar-winds-token` | string | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `--new-relic-api-key` | string | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `--datadog-api-key` | string | Applicable only for destination type: datadog. API key provided by Datadog. |
| `--action-send-enabled` | bool | Whether or not sending is enabled for action logs. |
| `--sftp-action-send-enabled` | bool | Whether or not sending is enabled for sftp_action logs. |
| `--ftp-action-send-enabled` | bool | Whether or not sending is enabled for ftp_action logs. |
| `--web-dav-action-send-enabled` | bool | Whether or not sending is enabled for web_dav_action logs. |
| `--sync-send-enabled` | bool | Whether or not sending is enabled for sync logs. |
| `--outbound-connection-send-enabled` | bool | Whether or not sending is enabled for outbound_connection logs. |
| `--automation-send-enabled` | bool | Whether or not sending is enabled for automation logs. |
| `--api-request-send-enabled` | bool | Whether or not sending is enabled for api_request logs. |
| `--public-hosting-request-send-enabled` | bool | Whether or not sending is enabled for public_hosting_request logs. |
| `--email-send-enabled` | bool | Whether or not sending is enabled for email logs. |
| `--exavault-api-request-send-enabled` | bool | Whether or not sending is enabled for exavault_api_request logs. |
| `--settings-change-send-enabled` | bool | Whether or not sending is enabled for settings_change logs. |
| `--destination-type` | enum | Destination Type. One of: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`. **Required.** |
| `--destination-url` | string | Destination Url |

### `files-cli siem-http-destinations send-test-entry`

send_test_entry SIEM HTTP Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--siem-http-destination-id` | int64 | SIEM HTTP Destination ID |
| `--destination-type` | enum | Destination Type. One of: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`. |
| `--destination-url` | string | Destination Url |
| `--name` | string | Name for this Destination |
| `--additional-headers` | object | Additional HTTP Headers included in calls to the destination URL |
| `--sending-active` | bool | Whether this SIEM HTTP Destination is currently being sent to or not |
| `--generic-payload-type` | enum | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. One of: `json_newline`, `json_array`. |
| `--file-destination-path` | string | Applicable only for destination type: file. Destination folder path on Files.com. |
| `--file-format` | enum | Applicable only for destination type: file. Generated file format. One of: `json`, `csv`. |
| `--file-interval-minutes` | int64 | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `--splunk-token` | string | Applicable only for destination type: splunk. Authentication token provided by Splunk. |
| `--azure-dcr-immutable-id` | string | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `--azure-stream-name` | string | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `--azure-oauth-client-credentials-tenant-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `--azure-oauth-client-credentials-client-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `--azure-oauth-client-credentials-client-secret` | string | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `--qradar-username` | string | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `--qradar-password` | string | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `--solar-winds-token` | string | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `--new-relic-api-key` | string | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `--datadog-api-key` | string | Applicable only for destination type: datadog. API key provided by Datadog. |
| `--action-send-enabled` | bool | Whether or not sending is enabled for action logs. |
| `--sftp-action-send-enabled` | bool | Whether or not sending is enabled for sftp_action logs. |
| `--ftp-action-send-enabled` | bool | Whether or not sending is enabled for ftp_action logs. |
| `--web-dav-action-send-enabled` | bool | Whether or not sending is enabled for web_dav_action logs. |
| `--sync-send-enabled` | bool | Whether or not sending is enabled for sync logs. |
| `--outbound-connection-send-enabled` | bool | Whether or not sending is enabled for outbound_connection logs. |
| `--automation-send-enabled` | bool | Whether or not sending is enabled for automation logs. |
| `--api-request-send-enabled` | bool | Whether or not sending is enabled for api_request logs. |
| `--public-hosting-request-send-enabled` | bool | Whether or not sending is enabled for public_hosting_request logs. |
| `--email-send-enabled` | bool | Whether or not sending is enabled for email logs. |
| `--exavault-api-request-send-enabled` | bool | Whether or not sending is enabled for exavault_api_request logs. |
| `--settings-change-send-enabled` | bool | Whether or not sending is enabled for settings_change logs. |

### `files-cli siem-http-destinations update`

Update SIEM HTTP Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Siem Http Destination ID. **Required.** |
| `--name` | string | Name for this Destination |
| `--additional-headers` | object | Additional HTTP Headers included in calls to the destination URL |
| `--sending-active` | bool | Whether this SIEM HTTP Destination is currently being sent to or not |
| `--generic-payload-type` | enum | Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. One of: `json_newline`, `json_array`. |
| `--file-destination-path` | string | Applicable only for destination type: file. Destination folder path on Files.com. |
| `--file-format` | enum | Applicable only for destination type: file. Generated file format. One of: `json`, `csv`. |
| `--file-interval-minutes` | int64 | Applicable only for destination type: file. Interval, in minutes, between file deliveries. Valid values are 5, 10, 15, 20, 30, 60, 90, 180, 240, 360. |
| `--splunk-token` | string | Applicable only for destination type: splunk. Authentication token provided by Splunk. |
| `--azure-dcr-immutable-id` | string | Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule. |
| `--azure-stream-name` | string | Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table. |
| `--azure-oauth-client-credentials-tenant-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID. |
| `--azure-oauth-client-credentials-client-id` | string | Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID. |
| `--azure-oauth-client-credentials-client-secret` | string | Applicable only for destination type: azure. Client Credentials OAuth Client Secret. |
| `--qradar-username` | string | Applicable only for destination type: qradar. Basic auth username provided by QRadar. |
| `--qradar-password` | string | Applicable only for destination type: qradar. Basic auth password provided by QRadar. |
| `--solar-winds-token` | string | Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds. |
| `--new-relic-api-key` | string | Applicable only for destination type: new_relic. API key provided by New Relic. |
| `--datadog-api-key` | string | Applicable only for destination type: datadog. API key provided by Datadog. |
| `--action-send-enabled` | bool | Whether or not sending is enabled for action logs. |
| `--sftp-action-send-enabled` | bool | Whether or not sending is enabled for sftp_action logs. |
| `--ftp-action-send-enabled` | bool | Whether or not sending is enabled for ftp_action logs. |
| `--web-dav-action-send-enabled` | bool | Whether or not sending is enabled for web_dav_action logs. |
| `--sync-send-enabled` | bool | Whether or not sending is enabled for sync logs. |
| `--outbound-connection-send-enabled` | bool | Whether or not sending is enabled for outbound_connection logs. |
| `--automation-send-enabled` | bool | Whether or not sending is enabled for automation logs. |
| `--api-request-send-enabled` | bool | Whether or not sending is enabled for api_request logs. |
| `--public-hosting-request-send-enabled` | bool | Whether or not sending is enabled for public_hosting_request logs. |
| `--email-send-enabled` | bool | Whether or not sending is enabled for email logs. |
| `--exavault-api-request-send-enabled` | bool | Whether or not sending is enabled for exavault_api_request logs. |
| `--settings-change-send-enabled` | bool | Whether or not sending is enabled for settings_change logs. |
| `--destination-type` | enum | Destination Type. One of: `generic`, `splunk`, `azure_legacy`, `qradar`, `sumo`, `rapid7`, `solar_winds`, `new_relic`, `datadog`, `azure`, `file`. |
| `--destination-url` | string | Destination Url |

### `files-cli siem-http-destinations delete`

Delete SIEM HTTP Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Siem Http Destination ID. **Required.** |

