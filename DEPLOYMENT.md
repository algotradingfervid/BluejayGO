# Bluejay CMS — Production Deployment Guide

## Production Architecture

```
                    INTERNET
                       |
                       v
                 [Caddy Server]
              (TLS, Reverse Proxy)
                       |
              +--------+--------+
              |                 |
        Static Files      Reverse Proxy
      (/public, /uploads)   (localhost:28090)
              |                 |
              |                 v
              |         [Go Binary: bluejay-cms]
              |           (Echo Framework)
              |                 |
              |                 v
              |           [SQLite Database]
              |              bluejay.db
              |          (WAL mode, Foreign Keys)
              |                 |
              +--------+--------+
                       |
                       v
                 [Litestream]
            (Continuous Replication)
                       |
                       v
                   [AWS S3]
              (Remote Backup Storage)
```

## Prerequisites

### Hardware Requirements
- Linux server (Ubuntu 20.04+ or Debian 11+ recommended)
- 2GB RAM minimum (4GB recommended)
- 10GB disk space minimum
- Stable internet connection

### Software Requirements
- Go 1.25+ (for building the binary)
- Caddy 2.x (web server and reverse proxy)
- Litestream (SQLite backup tool)
- systemd (for process management)

### External Services
- Domain name with DNS configured
- AWS S3 bucket for database backups (optional but highly recommended)
- AWS Access Key ID and Secret Access Key (for Litestream)

## Build for Production

### Cross-Compilation for Linux

If building from macOS or Windows, compile for Linux x64:

```bash
GOOS=linux GOARCH=amd64 go build -o bluejay-cms cmd/server/main.go
```

Or use the Makefile target:

```bash
make deploy-build
```

This produces a single static binary `bluejay-cms` with no external dependencies.

### Verify the Build

Check the binary:

```bash
file bluejay-cms
# Expected output: bluejay-cms: ELF 64-bit LSB executable, x86-64...
```

## Server Setup (Step-by-Step)

### 1. Create Directory Structure

SSH into your production server and create the application directory:

```bash
sudo mkdir -p /var/www/bluejay-cms
sudo mkdir -p /var/www/bluejay-cms/public/uploads
sudo mkdir -p /var/www/bluejay-cms/templates
sudo mkdir -p /var/www/bluejay-cms/db/migrations
```

### 2. Create www-data User

The application runs as the `www-data` user for security isolation:

```bash
# www-data usually exists on Debian/Ubuntu, but verify:
id www-data

# If it doesn't exist, create it:
sudo useradd -r -s /bin/false www-data
```

Set ownership:

```bash
sudo chown -R www-data:www-data /var/www/bluejay-cms
```

### 3. Upload Application Files

From your local machine, upload the necessary files:

```bash
# Upload the Go binary
scp bluejay-cms user@yourserver:/tmp/
ssh user@yourserver 'sudo mv /tmp/bluejay-cms /var/www/bluejay-cms/ && sudo chmod +x /var/www/bluejay-cms/bluejay-cms'

# Upload templates directory (recursive)
scp -r templates/ user@yourserver:/tmp/templates
ssh user@yourserver 'sudo mv /tmp/templates /var/www/bluejay-cms/'

# Upload public assets
scp -r public/ user@yourserver:/tmp/public
ssh user@yourserver 'sudo mv /tmp/public /var/www/bluejay-cms/'

# Upload database migrations
scp -r db/migrations/ user@yourserver:/tmp/migrations
ssh user@yourserver 'sudo mv /tmp/migrations /var/www/bluejay-cms/db/'

# Fix ownership after upload
ssh user@yourserver 'sudo chown -R www-data:www-data /var/www/bluejay-cms'
```

### 4. Install and Configure Caddy

#### Install Caddy

```bash
# Install Caddy (Debian/Ubuntu)
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

#### Configure Caddyfile

Edit `/etc/caddy/Caddyfile`:

```bash
sudo nano /etc/caddy/Caddyfile
```

Replace contents with:

```caddy
yourdomain.com {
    # Reverse proxy all non-static requests to the Go application
    reverse_proxy localhost:28090

    # Enable gzip compression for faster transfers
    encode gzip

    # Security headers
    header {
        # Force HTTPS for 1 year, apply to all subdomains (set by Caddy)
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"

        # Prevent MIME type sniffing
        X-Content-Type-Options "nosniff"

        # Prevent clickjacking attacks
        X-Frame-Options "DENY"

        # Enable XSS protection (legacy browsers)
        X-XSS-Protection "1; mode=block"

        # Control referrer information sent to external sites
        Referrer-Policy "strict-origin-when-cross-origin"
    }

    # Serve static files directly (bypass Go app for performance)
    file_server /public/* {
        root /var/www/bluejay-cms
    }

    # Serve uploaded files directly
    file_server /uploads/* {
        root /var/www/bluejay-cms/public
    }

    # Cache static assets for 1 year (they're immutable)
    header /public/* Cache-Control "public, max-age=31536000, immutable"

    # Cache uploaded files for 1 year
    header /uploads/* Cache-Control "public, max-age=31536000, immutable"
}
```

**Caddyfile Explanation:**

- `yourdomain.com`: Replace with your actual domain. Caddy automatically obtains and renews Let's Encrypt TLS certificates.
- `reverse_proxy localhost:28090`: Routes requests to the Go application running on port 28090.
- `encode gzip`: Compresses text responses (HTML, CSS, JS) before sending to clients.
- `Strict-Transport-Security`: Tells browsers to always use HTTPS for your domain.
- `file_server` directives: Serve static files directly from filesystem without hitting the Go app (much faster).
- `Cache-Control` headers: Allow browsers and CDNs to cache static assets aggressively.

#### Test and Start Caddy

```bash
# Validate Caddyfile syntax
sudo caddy validate --config /etc/caddy/Caddyfile

# Enable Caddy to start on boot
sudo systemctl enable caddy

# Start Caddy
sudo systemctl start caddy

# Check status
sudo systemctl status caddy
```

### 5. Install and Configure systemd Service

#### Create systemd Service File

```bash
sudo nano /etc/systemd/system/bluejay-cms.service
```

Paste the following:

```ini
[Unit]
Description=BlueJay CMS Website
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/bluejay-cms
ExecStart=/var/www/bluejay-cms/bluejay-cms
Restart=always
RestartSec=5s
Environment="ENVIRONMENT=production"
Environment="DATABASE_PATH=/var/www/bluejay-cms/bluejay.db"

[Install]
WantedBy=multi-user.target
```

**Service Configuration Explanation:**

- `Description`: Human-readable description shown in systemctl output.
- `After=network.target`: Ensures network is available before starting.
- `Type=simple`: The process runs in the foreground (doesn't daemonize).
- `User=www-data`: Runs as unprivileged user for security.
- `WorkingDirectory`: Sets the current directory for the binary (affects relative paths).
- `ExecStart`: Full path to the binary to execute.
- `Restart=always`: Automatically restart on crashes or clean exits.
- `RestartSec=5s`: Wait 5 seconds before restarting (prevents rapid restart loops).
- `Environment`: Sets environment variables for the Go application.
- `WantedBy=multi-user.target`: Enables the service to start automatically on boot.

#### Enable and Start Service

```bash
# Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable bluejay-cms

# Start the service now
sudo systemctl start bluejay-cms

# Check service status
sudo systemctl status bluejay-cms

# View logs in real-time
sudo journalctl -u bluejay-cms -f
```

### 6. Install and Configure Litestream

Litestream provides continuous, real-time replication of your SQLite database to S3.

#### Install Litestream

```bash
# Download and install Litestream
wget https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.tar.gz
tar -xzf litestream-v0.3.13-linux-amd64.tar.gz
sudo mv litestream /usr/local/bin/
sudo chmod +x /usr/local/bin/litestream

# Verify installation
litestream version
```

#### Configure Litestream

Create configuration file:

```bash
sudo mkdir -p /etc/litestream
sudo nano /etc/litestream/litestream.yml
```

Paste the following configuration:

```yaml
dbs:
  - path: /var/www/bluejay-cms/bluejay.db
    replicas:
      - type: s3
        bucket: your-backup-bucket
        path: bluejay-cms
        region: us-east-1
        # access-key-id: $AWS_ACCESS_KEY_ID
        # secret-access-key: $AWS_SECRET_ACCESS_KEY
```

**Litestream Configuration Explanation:**

- `path`: Location of your SQLite database file.
- `type: s3`: Uses AWS S3 for backup storage.
- `bucket`: Your S3 bucket name (must already exist).
- `path`: Prefix path within the bucket (organizes backups).
- `region`: AWS region where your bucket is located.
- `access-key-id` / `secret-access-key`: AWS credentials (set via environment variables for security).

#### Set AWS Credentials

Create environment file for Litestream:

```bash
sudo nano /etc/litestream/litestream.env
```

Add your AWS credentials:

```bash
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
```

Secure the file:

```bash
sudo chmod 600 /etc/litestream/litestream.env
sudo chown root:root /etc/litestream/litestream.env
```

#### Create Litestream systemd Service

```bash
sudo nano /etc/systemd/system/litestream.service
```

Paste:

```ini
[Unit]
Description=Litestream SQLite Replication
After=network.target

[Service]
Type=simple
User=www-data
EnvironmentFile=/etc/litestream/litestream.env
ExecStart=/usr/local/bin/litestream replicate -config /etc/litestream/litestream.yml
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable litestream
sudo systemctl start litestream
sudo systemctl status litestream
```

**Litestream Backup Strategy:**

- Continuous replication: Changes are streamed to S3 in near real-time (typically within 1 second).
- Point-in-time recovery: Restore your database to any point in time (down to the second).
- Low overhead: Litestream adds minimal CPU/memory overhead to your application.
- Disaster recovery: If your server dies, restore your database from S3 in seconds.

## First Deployment Checklist

Before going live, verify all components:

```bash
# 1. Check Go binary is executable and owned by www-data
ls -lh /var/www/bluejay-cms/bluejay-cms

# 2. Verify directory permissions
ls -ld /var/www/bluejay-cms
ls -ld /var/www/bluejay-cms/public/uploads

# 3. Check service is running
sudo systemctl status bluejay-cms

# 4. Verify port 28090 is listening
sudo netstat -tlnp | grep 28090

# 5. Check Caddy is running and TLS is active
sudo systemctl status caddy
curl -I https://yourdomain.com

# 6. Check Litestream is replicating
sudo systemctl status litestream
litestream databases -config /etc/litestream/litestream.yml

# 7. Test health check endpoint
curl http://localhost:28090/health

# 8. Review logs for errors
sudo journalctl -u bluejay-cms -n 50
sudo journalctl -u caddy -n 50
sudo journalctl -u litestream -n 50

# 9. Test the admin login page
curl https://yourdomain.com/admin/login

# 10. Verify file uploads work (check write permissions)
ls -ld /var/www/bluejay-cms/public/uploads
```

## Production Hardening

### 1. Change Session Secret

**CRITICAL:** Change the default session secret in production.

Generate a secure random secret:

```bash
openssl rand -base64 32
```

Update the systemd service file:

```bash
sudo nano /etc/systemd/system/bluejay-cms.service
```

Add the environment variable:

```ini
Environment="SESSION_SECRET=your-generated-secret-here"
```

Update your code to read from environment:

```go
// In cmd/server/main.go
secret := os.Getenv("SESSION_SECRET")
if secret == "" {
    secret = "change-this-secret-in-production-minimum-32-chars"
    logger.Warn("using default session secret - set SESSION_SECRET environment variable")
}
customMiddleware.InitSessionStore(secret)
```

Restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart bluejay-cms
```

### 2. Enable Secure Cookie Flag

In `internal/middleware/session.go`, update the session options:

```go
SessionStore.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   86400 * 7,
    HttpOnly: true,
    Secure:   true,  // Changed from false to true
    SameSite: http.SameSiteLaxMode,
}
```

**Important:** Only set `Secure: true` after TLS is working. This flag prevents cookies from being sent over HTTP.

### 3. Tighten Content Security Policy

Edit `internal/middleware/security.go` to restrict CSP:

```go
// Remove 'unsafe-inline' and 'unsafe-eval' from script-src
c.Response().Header().Set("Content-Security-Policy",
    "default-src 'self'; "+
    "script-src 'self' cdn.tailwindcss.com cdn.jsdelivr.net; "+  // Removed unsafe-inline/eval
    "style-src 'self' 'unsafe-inline' fonts.googleapis.com cdn.tailwindcss.com; "+
    "font-src 'self' fonts.gstatic.com; "+
    "img-src 'self' data: https:;")
```

**Note:** This may require refactoring inline scripts to external files.

### 4. Enable CSRF Protection on All POST Routes

Update the middleware stack in `cmd/server/main.go`:

```go
// Add CSRF middleware to admin routes
adminGroup := e.Group("/admin",
    customMiddleware.RequireAuth(),
    customMiddleware.CSRF(),  // Add CSRF protection
)
```

Implement CSRF middleware in `internal/middleware/csrf.go` (if not already present).

### 5. Use Environment Variables for Sensitive Config

Create a production environment file:

```bash
sudo nano /etc/bluejay-cms/production.env
```

Add sensitive configuration:

```bash
DATABASE_PATH=/var/www/bluejay-cms/bluejay.db
SESSION_SECRET=your-32-char-secret
ENVIRONMENT=production
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
ADMIN_EMAIL=admin@yourdomain.com
```

Update systemd service to use this file:

```ini
[Service]
EnvironmentFile=/etc/bluejay-cms/production.env
```

Secure the file:

```bash
sudo chmod 600 /etc/bluejay-cms/production.env
sudo chown root:root /etc/bluejay-cms/production.env
```

### 6. Set Restrictive File Permissions

```bash
# Binary should be executable only
sudo chmod 755 /var/www/bluejay-cms/bluejay-cms

# Database should be readable/writable only by www-data
sudo chmod 600 /var/www/bluejay-cms/bluejay.db

# Templates and static files can be world-readable
sudo chmod -R 755 /var/www/bluejay-cms/templates
sudo chmod -R 755 /var/www/bluejay-cms/public

# Upload directory needs write permissions
sudo chmod 755 /var/www/bluejay-cms/public/uploads
sudo chown -R www-data:www-data /var/www/bluejay-cms/public/uploads
```

### 7. Enable Firewall

```bash
# Install ufw (if not already installed)
sudo apt install ufw

# Allow SSH (important - don't lock yourself out!)
sudo ufw allow OpenSSH

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Verify status
sudo ufw status verbose
```

### 8. Set Up Automatic Security Updates

```bash
# Install unattended-upgrades
sudo apt install unattended-upgrades

# Enable automatic security updates
sudo dpkg-reconfigure -plow unattended-upgrades
```

## Updating / Redeploying

### Quick Update (Binary Only)

Use the Makefile deploy workflow:

```bash
# From your local machine:
make deploy
```

This runs three targets sequentially:
1. `make deploy-build`: Compiles for Linux x64
2. `make deploy-upload`: Uploads binary and config files
3. `make deploy-restart`: Restarts services

### Manual Update Process

```bash
# 1. Build new binary locally
GOOS=linux GOARCH=amd64 go build -o bluejay-cms cmd/server/main.go

# 2. Upload to server
scp bluejay-cms user@yourserver:/tmp/

# 3. SSH to server and replace binary
ssh user@yourserver
sudo systemctl stop bluejay-cms
sudo mv /tmp/bluejay-cms /var/www/bluejay-cms/bluejay-cms
sudo chmod +x /var/www/bluejay-cms/bluejay-cms
sudo chown www-data:www-data /var/www/bluejay-cms/bluejay-cms
sudo systemctl start bluejay-cms

# 4. Verify deployment
sudo systemctl status bluejay-cms
curl http://localhost:28090/health
```

### Zero-Downtime Deployment (Advanced)

For critical production systems, use systemd socket activation or run two instances behind a load balancer.

### Database Migrations

If your update includes schema changes:

```bash
# Upload new migration files
scp db/migrations/*.sql user@yourserver:/var/www/bluejay-cms/db/migrations/

# SSH to server and restart (migrations run automatically on startup)
ssh user@yourserver
sudo systemctl restart bluejay-cms

# Check logs to verify migrations succeeded
sudo journalctl -u bluejay-cms -n 50
```

### Rollback Procedure

If deployment fails:

```bash
# 1. Keep old binary as backup before deploying
sudo cp /var/www/bluejay-cms/bluejay-cms /var/www/bluejay-cms/bluejay-cms.backup

# 2. If new version fails, restore backup
sudo systemctl stop bluejay-cms
sudo mv /var/www/bluejay-cms/bluejay-cms.backup /var/www/bluejay-cms/bluejay-cms
sudo systemctl start bluejay-cms

# 3. Restore database from Litestream if needed (see Backup and Recovery section)
```

## Monitoring

### Health Check Endpoint

The application provides a health check at `/health`:

```bash
curl http://localhost:28090/health
```

Expected response:

```json
{
  "status": "ok",
  "time": "2026-02-10T15:30:45Z"
}
```

Configure external monitoring services (UptimeRobot, Pingdom, etc.) to poll this endpoint every 60 seconds.

### systemd Service Logs

View real-time logs:

```bash
# Follow logs (live tail)
sudo journalctl -u bluejay-cms -f

# View last 100 lines
sudo journalctl -u bluejay-cms -n 100

# View logs since 1 hour ago
sudo journalctl -u bluejay-cms --since "1 hour ago"

# View logs with specific priority (error and above)
sudo journalctl -u bluejay-cms -p err

# Export logs to file
sudo journalctl -u bluejay-cms --since "2026-02-01" > /tmp/bluejay-logs.txt
```

### Caddy Access Logs

Caddy logs to syslog by default:

```bash
# View Caddy logs
sudo journalctl -u caddy -f

# View access logs (if configured separately)
sudo tail -f /var/log/caddy/access.log
```

Enable JSON access logging in Caddyfile:

```caddy
yourdomain.com {
    log {
        output file /var/log/caddy/access.log
        format json
    }
    # ... rest of config
}
```

### Application Performance Monitoring (APM)

Consider integrating:

- **Prometheus + Grafana**: Metrics collection and visualization
- **Sentry**: Error tracking and alerting
- **New Relic / DataDog**: Full-stack APM

### Custom Alerts

Set up systemd email alerts on service failures:

```bash
sudo nano /etc/systemd/system/bluejay-cms-alert@.service
```

```ini
[Unit]
Description=Service Alert for %i

[Service]
Type=oneshot
ExecStart=/usr/local/bin/alert.sh %i
```

Create alert script:

```bash
sudo nano /usr/local/bin/alert.sh
sudo chmod +x /usr/local/bin/alert.sh
```

```bash
#!/bin/bash
SERVICE=$1
echo "Service $SERVICE has failed!" | mail -s "Alert: $SERVICE Failed" admin@yourdomain.com
```

## Backup and Recovery

### Litestream Continuous Replication

Litestream automatically backs up your database to S3 in real-time.

**Verify Replication Status:**

```bash
# Check replication is active
litestream databases -config /etc/litestream/litestream.yml

# View replication lag
litestream snapshots -config /etc/litestream/litestream.yml /var/www/bluejay-cms/bluejay.db
```

Expected output shows recent snapshots and WAL files.

### Restore from S3 (Disaster Recovery)

If your server fails and you need to restore:

```bash
# Stop the application
sudo systemctl stop bluejay-cms

# Restore database from S3 (restores to latest point in time)
litestream restore -config /etc/litestream/litestream.yml /var/www/bluejay-cms/bluejay.db

# Restore to specific point in time (if needed)
litestream restore -config /etc/litestream/litestream.yml -timestamp 2026-02-10T12:00:00Z /var/www/bluejay-cms/bluejay.db

# Fix permissions
sudo chown www-data:www-data /var/www/bluejay-cms/bluejay.db
sudo chmod 600 /var/www/bluejay-cms/bluejay.db

# Start the application
sudo systemctl start bluejay-cms
```

### Manual Database Backup (sqlite3)

For ad-hoc backups before major changes:

```bash
# Create backup
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db ".backup /var/www/bluejay-cms/backup-$(date +%Y%m%d-%H%M%S).db"

# Verify backup
ls -lh /var/www/bluejay-cms/backup-*.db
```

### Restore from Manual Backup

```bash
# Stop application
sudo systemctl stop bluejay-cms

# Copy backup to production database
sudo cp /var/www/bluejay-cms/backup-20260210-120000.db /var/www/bluejay-cms/bluejay.db

# Fix permissions
sudo chown www-data:www-data /var/www/bluejay-cms/bluejay.db
sudo chmod 600 /var/www/bluejay-cms/bluejay.db

# Start application
sudo systemctl start bluejay-cms
```

### Backup Uploaded Files

User uploads are not backed up by Litestream (they're not in the database). Use rsync:

```bash
# Set up daily cron job to backup uploads to S3
sudo crontab -e
```

Add:

```bash
0 2 * * * aws s3 sync /var/www/bluejay-cms/public/uploads s3://your-backup-bucket/uploads/ --delete
```

Or use AWS CLI to sync uploads:

```bash
aws s3 sync /var/www/bluejay-cms/public/uploads s3://your-backup-bucket/uploads/ --delete
```

## Troubleshooting Common Issues

### Issue: Port 28090 Already in Use

**Symptoms:**
```
Error: bind: address already in use
```

**Solution:**

```bash
# Find process using port 28090
sudo lsof -i :28090

# Kill the process
sudo kill -9 <PID>

# Or stop the service properly
sudo systemctl stop bluejay-cms
sudo systemctl start bluejay-cms
```

### Issue: Database Locked

**Symptoms:**
```
Error: database is locked
```

**Causes:**
- Multiple processes trying to write simultaneously
- Connection pool misconfiguration
- Long-running transaction not committed

**Solution:**

```bash
# Check for multiple instances
ps aux | grep bluejay-cms

# Kill duplicate processes
sudo systemctl restart bluejay-cms

# Check WAL mode is enabled
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db "PRAGMA journal_mode;"
# Should output: wal

# Force WAL mode if not enabled
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db "PRAGMA journal_mode=WAL;"
```

### Issue: Permission Denied Errors

**Symptoms:**
```
Error: open /var/www/bluejay-cms/bluejay.db: permission denied
```

**Solution:**

```bash
# Fix ownership
sudo chown -R www-data:www-data /var/www/bluejay-cms

# Fix database permissions
sudo chmod 600 /var/www/bluejay-cms/bluejay.db

# Fix directory permissions
sudo chmod 755 /var/www/bluejay-cms

# Fix upload directory permissions
sudo chmod 755 /var/www/bluejay-cms/public/uploads
```

### Issue: Caddy TLS Certificate Failures

**Symptoms:**
```
Error: obtaining certificate: acme: error presenting token
```

**Causes:**
- DNS not pointing to server
- Firewall blocking port 80
- Rate limiting from Let's Encrypt

**Solution:**

```bash
# Verify DNS is correct
dig yourdomain.com +short
# Should return your server's IP

# Check firewall allows HTTP
sudo ufw status | grep 80

# View Caddy logs for detailed error
sudo journalctl -u caddy -n 50

# Test certificate issuance manually
sudo caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
```

If rate-limited, wait 1 hour or use Let's Encrypt staging server:

```caddy
yourdomain.com {
    tls {
        ca https://acme-staging-v02.api.letsencrypt.org/directory
    }
    # ... rest of config
}
```

### Issue: 502 Bad Gateway

**Symptoms:** Caddy returns 502 error when accessing the site.

**Causes:**
- Go application not running
- Go application crashed
- Port mismatch

**Solution:**

```bash
# Check if application is running
sudo systemctl status bluejay-cms

# Check if port 28090 is listening
sudo netstat -tlnp | grep 28090

# Check application logs for crashes
sudo journalctl -u bluejay-cms -n 100

# Restart application
sudo systemctl restart bluejay-cms
```

### Issue: Session Not Persisting (Logged Out Immediately)

**Symptoms:** Can't stay logged in, redirected to login page after authentication.

**Causes:**
- Session secret changed (invalidated existing sessions)
- Secure flag enabled but accessing via HTTP
- Cookie domain/path mismatch

**Solution:**

```bash
# Check session configuration in logs
sudo journalctl -u bluejay-cms | grep -i session

# Verify HTTPS is working
curl -I https://yourdomain.com

# Test cookies are being set
curl -v -X POST https://yourdomain.com/admin/login -d "email=admin@example.com&password=test" -c cookies.txt
cat cookies.txt
```

If using HTTPS, ensure `Secure: true` is set in session options. If testing locally with HTTP, set `Secure: false`.

### Issue: File Uploads Failing

**Symptoms:**
```
Error: failed to save file: permission denied
```

**Solution:**

```bash
# Check upload directory exists and is writable
ls -ld /var/www/bluejay-cms/public/uploads

# Fix permissions
sudo mkdir -p /var/www/bluejay-cms/public/uploads
sudo chown www-data:www-data /var/www/bluejay-cms/public/uploads
sudo chmod 755 /var/www/bluejay-cms/public/uploads

# Check disk space
df -h /var/www/bluejay-cms
```

### Issue: Litestream Not Replicating

**Symptoms:** Database not being backed up to S3.

**Solution:**

```bash
# Check Litestream service status
sudo systemctl status litestream

# View Litestream logs
sudo journalctl -u litestream -n 50

# Test AWS credentials
aws s3 ls s3://your-backup-bucket/ --region us-east-1

# Verify configuration
litestream databases -config /etc/litestream/litestream.yml

# Manually trigger replication
sudo systemctl restart litestream
```

Common causes:
- Invalid AWS credentials
- S3 bucket doesn't exist
- IAM permissions insufficient (needs s3:PutObject, s3:GetObject, s3:ListBucket)

### Issue: High Memory Usage

**Symptoms:** Server running out of memory, OOM killer terminating processes.

**Solution:**

```bash
# Check memory usage
free -h
top -o %MEM

# Check if SQLite cache is too large
# Default cache_size=2000 pages (~8MB) should be fine

# Add memory limits to systemd service
sudo nano /etc/systemd/system/bluejay-cms.service
```

Add under `[Service]`:

```ini
MemoryMax=512M
MemoryHigh=400M
```

```bash
sudo systemctl daemon-reload
sudo systemctl restart bluejay-cms
```

## Scaling Considerations

### SQLite Limits

SQLite is suitable for most small-to-medium websites:

**What SQLite Can Handle:**
- 100,000+ requests per day
- Databases up to 281 TB (theoretical limit)
- Thousands of concurrent readers
- ~1000 writes per second (with WAL mode)

**When SQLite is Enough:**
- Single-server deployments
- Read-heavy workloads
- < 100 concurrent users
- < 10 writes per second

### When to Consider PostgreSQL

Migrate to PostgreSQL if you encounter:

- **High write concurrency**: Multiple servers need to write simultaneously
- **Geographic distribution**: Need multi-region database replication
- **Advanced features**: Full-text search, JSON indexing, custom extensions
- **Horizontal scaling**: Sharding or read replicas across multiple servers
- **Team preferences**: Developers more comfortable with PostgreSQL

### Horizontal Scaling (Multiple Servers)

If outgrowing a single server:

1. **Load Balancer**: Add nginx or AWS ALB in front of multiple app servers
2. **Shared Database**: Move SQLite to PostgreSQL or use managed database (AWS RDS, DigitalOcean Managed Databases)
3. **Shared Storage**: Move uploads to S3 or object storage
4. **Session Store**: Move sessions to Redis or PostgreSQL
5. **CDN**: Serve static assets and uploads via CloudFront or Cloudflare

### Vertical Scaling (Bigger Server)

Often cheaper and simpler than horizontal scaling:

- Upgrade to 4GB/8GB RAM
- Add SSD storage
- Increase CPU cores

SQLite can easily handle 10x-100x more traffic on better hardware before needing architectural changes.

### Caching Strategies

Reduce database load with caching:

- **Application-level cache**: Already implemented in `internal/services/cache.go`
- **HTTP caching**: Already implemented with `Cache-Control` headers
- **CDN caching**: Put CloudFlare or Fastly in front of Caddy
- **Redis cache**: Add Redis for frequently accessed data

### Database Optimization

Improve SQLite performance:

```bash
# Analyze database and update statistics
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db "ANALYZE;"

# Check for missing indexes
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db ".schema" | grep -i index

# Vacuum database to reclaim space and optimize
sudo systemctl stop bluejay-cms
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db "VACUUM;"
sudo systemctl start bluejay-cms
```

---

## Summary

You now have a production-ready deployment of Bluejay CMS with:

- Automatic HTTPS via Caddy
- Process supervision via systemd
- Continuous database backups via Litestream
- Security hardening (headers, CSRF, secure sessions)
- Monitoring via health checks and journald logs
- Disaster recovery procedures

**Key Commands Reference:**

```bash
# View application logs
sudo journalctl -u bluejay-cms -f

# Restart application
sudo systemctl restart bluejay-cms

# Reload Caddy config
sudo systemctl reload caddy

# Check all services
sudo systemctl status bluejay-cms caddy litestream

# Backup database manually
sudo -u www-data sqlite3 /var/www/bluejay-cms/bluejay.db ".backup /tmp/backup.db"

# Restore from Litestream
litestream restore -config /etc/litestream/litestream.yml /var/www/bluejay-cms/bluejay.db
```

**Next Steps:**

1. Set up automated monitoring (UptimeRobot, Pingdom, etc.)
2. Configure email notifications for service failures
3. Document your specific deployment configuration
4. Create runbooks for common maintenance tasks
5. Schedule regular backup tests to verify recovery procedures

For questions or issues, check the logs first, then consult this guide's troubleshooting section.
