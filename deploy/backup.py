#!/usr/bin/env python3
"""Create a verified SQLite snapshot and upload archive; never archive secrets."""
import datetime
import fcntl
import grp
import hashlib
import io
import json
import os
from pathlib import Path
import sqlite3
import stat
import tarfile
import tempfile
import urllib.parse

BACKUPS = Path(os.environ.get("BACKUP_DIR", "/var/backups/bluejay-cms/daily"))
DATABASE = Path(os.environ.get("DB_PATH", "/var/lib/bluejay-cms/bluejay.db"))
UPLOADS = Path(os.environ.get("UPLOAD_DIR", "/var/www/bluejay-cms/public/uploads"))
GROUP = os.environ.get("BACKUP_GROUP", "bluejay")
RETENTION = int(os.environ.get("BACKUP_RETENTION", "14"))

def utc():
    return datetime.datetime.now(datetime.timezone.utc).replace(microsecond=0)

def publish_metadata(path, value):
    temp = path.with_suffix(".tmp")
    temp.write_text(json.dumps(value, indent=2) + "\n")
    os.chmod(temp, 0o640)
    if GROUP:
        os.chown(temp, -1, grp.getgrnam(GROUP).gr_gid)
    os.replace(temp, path)

def add_regular(archive, fd, name, manifest):
    before = os.fstat(fd)
    if not stat.S_ISREG(before.st_mode):
        raise ValueError("Backup source must contain only regular files")
    # Hash the same open descriptor that is archived, then detect concurrent edits.
    digest = hashlib.sha256()
    with os.fdopen(os.dup(fd), "rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
        stream.seek(0)
        info = tarfile.TarInfo(name)
        info.size, info.mode, info.mtime = before.st_size, 0o600, int(before.st_mtime)
        archive.addfile(info, stream)
    after = os.fstat(fd)
    if (before.st_size, before.st_mtime_ns) != (after.st_size, after.st_mtime_ns):
        raise RuntimeError("An upload changed during backup; retry required")
    manifest[name] = {"sha256": digest.hexdigest(), "size": before.st_size}

def add_uploads(archive, directory_fd, prefix, manifest):
    # Use openat/O_NOFOLLOW for every component: a writable upload tree must
    # never make the root-run backup job follow a link into private host files.
    for name in sorted(os.listdir(directory_fd)):
        info = os.stat(name, dir_fd=directory_fd, follow_symlinks=False)
        if stat.S_ISDIR(info.st_mode):
            child = os.open(name, os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW, dir_fd=directory_fd)
            try:
                add_uploads(archive, child, prefix + name + "/", manifest)
            finally:
                os.close(child)
        elif stat.S_ISREG(info.st_mode):
            fd = os.open(name, os.O_RDONLY | os.O_NOFOLLOW, dir_fd=directory_fd)
            try:
                add_regular(archive, fd, prefix + name, manifest)
            finally:
                os.close(fd)
        else:
            raise ValueError("Symlinks and special upload files are not backed up")

def main():
    os.umask(0o077)
    if RETENTION < 2:
        raise ValueError("Keep at least two daily backups")
    BACKUPS.mkdir(parents=True, exist_ok=True)
    with (BACKUPS / ".backup.lock").open("a") as lock:
        fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
        started = utc()
        filename = "bluejay-backup-" + started.strftime("%Y%m%dT%H%M%SZ") + ".tar.gz"
        try:
            with tempfile.TemporaryDirectory(prefix=".snapshot-", dir=BACKUPS) as work:
                db_copy = Path(work) / "bluejay.db"
                uri = "file:" + urllib.parse.quote(str(DATABASE)) + "?mode=ro"
                with sqlite3.connect(uri, uri=True) as source, sqlite3.connect(db_copy) as target:
                    source.backup(target)
                    if target.execute("PRAGMA quick_check").fetchone()[0] != "ok":
                        raise RuntimeError("SQLite backup integrity check failed")
                candidate = Path(work) / filename
                manifest = {}
                with tarfile.open(candidate, "w:gz") as archive:
                    fd = os.open(db_copy, os.O_RDONLY | os.O_NOFOLLOW)
                    try:
                        add_regular(archive, fd, "database/bluejay.db", manifest)
                    finally:
                        os.close(fd)
                    directory_fd = os.open(UPLOADS, os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW)
                    try:
                        add_uploads(archive, directory_fd, "uploads/", manifest)
                    finally:
                        os.close(directory_fd)
                    b = json.dumps({"created": started.isoformat(), "files": manifest}, indent=2).encode()
                    member = tarfile.TarInfo("manifest.json")
                    member.size, member.mode = len(b), 0o600
                    archive.addfile(member, io.BytesIO(b))
                # Check archive contents before publishing a download link.
                with tarfile.open(candidate, "r:gz") as archive:
                    for name, expected in manifest.items():
                        stream = archive.extractfile(name)
                        digest = hashlib.sha256()
                        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
                            digest.update(chunk)
                        if digest.hexdigest() != expected["sha256"]:
                            raise RuntimeError("Archive checksum verification failed")
                os.chmod(candidate, 0o640)
                if GROUP:
                    os.chown(candidate, -1, grp.getgrnam(GROUP).gr_gid)
                os.replace(candidate, BACKUPS / filename)
            publish_metadata(BACKUPS / "status.json", {"success": True, "finished": utc().isoformat(), "filename": filename})
            # Only prune this job's own archives, after a new verified success.
            for old in sorted(BACKUPS.glob("bluejay-backup-????????T??????Z.tar.gz"), reverse=True)[RETENTION:]:
                if old.is_file() and not old.is_symlink():
                    old.unlink()
            print(json.dumps({"success": True, "filename": filename, "files": len(manifest)}))
        except Exception:
            publish_metadata(BACKUPS / "status.json", {"success": False, "finished": utc().isoformat()})
            raise

if __name__ == "__main__":
    main()
