# Bluejay hardened deployment

The service runs as `bluejay`, binds `127.0.0.1:28090`, and requires a protected random 32-byte signing key at `/etc/bluejay/session.key`. Code, templates and configuration are root-owned. Only `/var/lib/bluejay-cms` and `/var/www/bluejay-cms/public/uploads` are writable by the app. The SQLite database and its WAL/SHM files belong in the private state directory. Do not reintroduce a shared fallback key or recursively assign application code to the service account.

## Runtime settings

- `SESSION_KEY_FILE`: file containing exactly 32 random bytes, root:bluejay, mode 0640, private containing directory. Alternatively, `SESSION_SECRET` accepts 64 hex characters. Do not set both.
- `ENVIRONMENT=production` requires Secure session cookies. Local HTTP development may use a generated `SESSION_SECRET` without production mode.
- `DB_PATH=/var/lib/bluejay-cms/bluejay.db`, `LISTEN_HOST=127.0.0.1`, `PORT=28090`.
- `BACKUP_DIR=/var/backups/bluejay-cms/daily` grants the service read-only access to completed backups.
- Rotate the signing key and restart to invalidate all existing signed cookies. Rotation forces users to sign in again. The separate session-revocation findings in the audit are not fully solved by key rotation.

## Daily backups and manual off-server copies

`bluejay-backup.timer` runs daily at 02:30 UTC (08:00 IST), with up to 10 minutes jitter. It keeps the newest 14 successful archives. Existing manual rollback snapshots are not pruned. A job failure exits nonzero into systemd/journald, records failure for the Backups page, and keeps previous archives. A successful snapshot older than 36 hours triggers a warning on the page. No email/webhook destination is configured.

Admin users can visit `/admin/backups` and download an archive. Editor accounts and anonymous visitors cannot. The archives contain `database/bluejay.db`, `uploads/`, and a checksum manifest. They exclude signing secrets, SSH keys and server credentials. They still contain private site content and administrator password hashes: store downloaded copies securely. Per the owner's instruction, there is no automated off-host replication; the team must download regularly. Server-local backups alone do not protect against server/disk loss.

The SQLite online backup API makes a consistent database snapshot; uploads are copied separately. The job detects files changing while read and refuses symlinks or special files. It does not promise a single atomic transaction spanning the database and filesystem. For a fully quiescent disaster-recovery snapshot, pause writes while backing up.

## Restore checks

Extract a trusted backup into a private empty directory, never directly over the running site. Verify every file against `manifest.json`; run SQLite `PRAGMA integrity_check` on the restored database. Test the application with that database and restored uploads on an isolated loopback port with a new test-only signing key. Do not expose a restored database to public routes. For an actual rollback, stop the service, preserve current data, restore the database and uploads with bluejay ownership/private modes, then restart and verify the site. The live session signing key remains separate.

## Release and server operations

The live hardening was deployed on 29 September 2026 from `codex/bluejay-live-security`, based on the previously deployed commit `931160c`. The security changes and both completed UX review rounds are now integrated into `main`. Use `main` for future releases.

The `bluejay-ops` SSH account accepts the existing authorized operator keys. Its only passwordless sudo entry is the root-owned `/usr/local/sbin/bluejay-operations` wrapper, accepting exactly one of `status`, `restart`, `backup`, or `logs`. For example: `sudo /usr/local/sbin/bluejay-operations backup`. Password and keyboard-interactive SSH authentication are disabled. Root key access remains available for recovery and deployment; both fresh root and operator key logins were checked after reload. Give each team member an individually managed public key when onboarding them.

Build with Go 1.27.1 or newer supported patched release; run tests, `go vet`, module verification and `govulncheck` for host/Linux. `deploy.sh` preserves the hardened ownership and refuses deployment without private key/state files. It does not install or change systemd/SSH configuration automatically. Changes to units/Caddy require validation before reload. Preserve the known-working key-based SSH recovery connection until a second login has passed.

The Caddy configuration has explicit header/body/idle/write time limits and a 64 MB body cap. The app further caps ordinary bodies at 4 MiB, multipart bodies at 64 MiB, form values at 512 and file uploads at 32 per request; existing per-file caps still apply. Backup downloads have a 30-minute response deadline. Large files or slow connections may require adjusted, tested limits.

This deployment addresses the live-configuration recommendations and adds protected backup downloads. It does not resolve every separate application finding in the security review, such as CSRF enforcement, login throttling, draft visibility or all upload validation issues.
