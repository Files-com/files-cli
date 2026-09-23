# Path comparison data

`path_comparison.json` is the version 1 comparison map for MySQL
`utf8mb4_0900_ai_ci`, shared with Files.com server path comparison.
`comparison_examples.json` contains exact input/output examples.

Source: `files-rails` commit `41cec53c4cc711a193a16fff5d1dbadaffb32de4`,
`gems/files-path-native/data`, published in `files-path-native` 0.1.1914082.
SHA-256 of `path_comparison.json`: `a64d5de7fd4450f9889b57a5ac3088eb41451545985017a2beada38aaad2e371`.

Normalize path syntax first, then replace each Unicode scalar once using the
hexadecimal keys in `mapping`. Preserve scalars without entries. An empty value
removes a scalar; a value can also contain several scalars. Do not normalize,
lowercase, trim, or apply path syntax normalization after replacement. Keep the
original path for requests and display; comparison keys are for matching only.

Update the map and server examples together from an approved server revision,
record the new digest here, and run every SDK's path comparison tests. Runtime
resources under each generated SDK are copied from this shared source by the
generator; do not edit those copies independently.
