---
name: recipe-share-and-notify
description: |
  Share a Files.com file or folder via a Share Link (bundle) and email it to recipients. Spans `bundles`, `bundle-recipients`, and (optionally) `bundle-notifications`. Use this when the user wants to share files with someone outside Files.com by email, or to notify an internal Files.com user when activity happens on a Share Link. Distinguishes the one-shot invitation email from the ongoing activity-notification subscription — they are different surfaces and different audiences.
---

# recipe-share-and-notify

Sharing a file or folder by email isn't a single call. It involves up to three domains:

- `bundles` — the Share Link itself (URL, code, permissions, expiration).
- `bundle-recipients` — the people who are invited to a Share Link, including the option to email them at creation.
- `bundle-notifications` — internal Files.com users subscribed to ongoing activity emails about the Share Link.

`bundles create` produces the link but does not email anyone. Email delivery and notifications are separate operations.

## When to use this skill

- Share a file or folder by email with one or more external addresses.
- Restrict a Share Link to specific recipients who have to receive the email to access it.
- Subscribe an internal user to ongoing activity emails (uploads, registrations) on a Share Link.
- A previous attempt to "send a share link" only created the bundle without delivering it.

## Sending the invitation email

There are two ways to send the invitation. Either one creates `BundleRecipient` rows for the addresses involved.

### Path A: create the bundle, then send

Use when you already have the bundle, or when you want to email multiple addresses at once. `bundles share` supports up to 50 addresses per call.

    files-cli bundles create --paths=/demo/q3-report.pdf --description="Q3 report" --format=json
    files-cli bundles share --id=<BUNDLE_ID> --to=bob@example.com --to=alice@example.com --note="Q3 report attached" --format=json

### Path B: create and add a recipient that triggers the email

Use when adding one recipient (or each recipient one at a time). The `--share-after-create=true` flag on `bundle-recipients create` is the "send mail at creation" path.

    files-cli bundles create --paths=/demo/q3-report.pdf --description="Q3 report" --format=json
    files-cli bundle-recipients create --bundle-id=<BUNDLE_ID> --recipient=bob@example.com --share-after-create=true --format=json

### Recipient-restricted Share Links

If the link should only be accessible to addresses that received the invitation, set `--require-share-recipient` at creation:

    files-cli bundles create --paths=/legal/contract.pdf --require-share-recipient --format=json
    files-cli bundles share --id=<BUNDLE_ID> --to=counsel@example.com --to=ceo@example.com --format=json

## Subscribing an internal user to ongoing activity

`bundle-notifications` is a separate surface from `bundles share` and `bundle-recipients`. It subscribes an existing Files.com user to receive emails when registration or upload activity happens against the Share Link.

Notifications can only target internal Files.com users — never external contacts. External addresses go through invitation emails.

    files-cli bundle-notifications create --bundle-id=<BUNDLE_ID> --notify-user-id=<USER_ID> --notify-on-upload=true --notify-on-registration=true --format=json

`--notify-user-id` must reference an existing Files.com user — look it up with `files-cli users list --search=email@example.com --format=json` first if you don't have the ID.

## Gotchas

**Per-recipient invitation emails are rate-limited to once every 24 hours.** A repeat `bundles share` to the same address within that window will not re-send.

**A single `bundles share` call can address up to 50 recipients.** Beyond that, split across multiple calls.

**The sender's email is the reply-to.** The email address of the user who issues the invitation is used as the reply-to address.

## Anti-patterns

- **Don't delete and recreate the bundle to trigger an email.** Creation never emails anyone; the rate limit on invitations makes a recreate-and-resend approach ineffective even if the agent thinks it might force delivery.
- **Don't try to use `bundle-notifications` to send the initial Share Link to an external recipient.** That's the wrong surface — notifications can only target internal Files.com users. Use `bundles share` or `bundle-recipients create --share-after-create=true` instead.
- **Don't add `bundle-recipients` rows without `--share-after-create=true` and then expect an email to go out.** Adding the recipient row alone records the access permission but does not deliver the link.

## Related skills

- `filescom-bundles` — full reference for the Share Link surface (create/update/share/delete/find/list).
- `filescom-bundle-recipients` — the recipient resource, including the `--share-after-create` flag.
- `filescom-bundle-notifications` — the activity-notification subscription resource.
- `filescom-bundle-downloads` — audit of who downloaded a bundle.
