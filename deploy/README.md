# eDairy Docker Deployment

Each dairy cooperative runs as an identical Docker Compose stack under `/opt/edairy/{slug}/`. HAProxy, MySQL, Certbot, and monitoring remain on the host.

## Architecture

```
Internet → HAProxy (TLS) → host ports → Docker containers
                                              ├── frontend (nginx + React SPA)
                                              └── backend (Go API)
                                                    └── MySQL on host (host.docker.internal)
```

## Images

| Image | Source |
|-------|--------|
| `ghcr.io/rubewafula/edairy-backend:{version}` | [edairy-go/Dockerfile](../edairy-go/Dockerfile) |
| `ghcr.io/rubewafula/edairy-frontend:{version}` | [react-admin-app/Dockerfile](../edairy-react-admin/react-admin-app/Dockerfile) |

Build and publish via GitHub Actions on tag push (`v1.0.0`) or manual workflow dispatch in each repo.

## Directory layout (server)

```text
/opt/edairy/
  templates/                 # shared templates (installed by create-dairy.sh)
  arithi/
    docker-compose.yml
    .env
    runtime-config.js
    data/uploads/
    data/storage/
  tigania-west/
    ...
```

## Port map

| Dairy | API port | UI port | Database |
|-------|----------|---------|----------|
| arithi | 8500 | 8501 | arithi_dairy_erp |
| tigania-west | 8503 | 8502 | tigania_west_dairy_erp |
| mukululu | 8504 | 8505 | mukululu_dairy_erp |
| mwimbi | 8506 | 8507 | mwimbi_dairy_erp |
| mutuati | 8508 | 8509 | mutuati_dairy_erp |
| nkuene | 8510 | 8511 | nkuene_dairy_erp |
| dev | 8512 | 8513 | dev_dairy_erp |
| s-butsotso | 8514 | 8515 | s_butsotso_dairy_erp |

## Provision a new dairy

```bash
# Copy deploy/ to server, then:
sudo MYSQL_ROOT='...' EDAIRY_VERSION=1.0.0 \
  ./deploy/create-dairy.sh mukululu 8504 8505
```

This will:

1. Create the MySQL database (schema-only copy from `arithi_dairy_erp`)
2. Grant `edairy@172.%` for Docker bridge access
3. Render `/opt/edairy/mukululu/` from templates
4. Patch HAProxy (if not already configured)
5. `docker compose pull && docker compose up -d`
6. Run health checks on `/health` and `/runtime-config.js`

## Migrate an existing systemd dairy

```bash
sudo EDAIRY_VERSION=1.0.0 \
  ./deploy/create-dairy.sh --migrate --skip-db arithi 8500 8501
```

The `--migrate` flag:

- Preserves JWT secret and DB password from `/apps/go/{slug}/.env`
- Copies `uploads/` and `storage/` into the stack data volumes
- Stops and disables the legacy systemd unit
- Removes the host nginx site config (backed up with timestamp)

## Rename a dairy

```bash
sudo MYSQL_ROOT='...' ./deploy/rename-dairy.sh nguene nkuene 8510 8511
```

## Development dairy

The `dev` stack (API 8512, UI 8513) is for integration testing and local frontend dev:

```bash
# Local frontend dev server targets dev dairy API
npm run dev   # uses env/.env.dev -> api.dev.edairy.africa

# Per-dairy build config (optional)
npm run build --dev-dairy
```


```bash
sudo EDAIRY_VERSION=1.0.0 ./deploy/migrate-all-to-docker.sh
```

Runs nguene first (rehearsal), then the remaining five. Assumes databases and HAProxy are already configured.

## Upgrade a dairy

```bash
cd /opt/edairy/arithi
export EDAIRY_VERSION=1.1.0
docker compose pull
docker compose up -d
```

## Rollback

```bash
cd /opt/edairy/arithi
docker compose down

# Restore legacy systemd (binary must still exist under /apps/go/arithi/)
systemctl start arithi-erp.service

# Restore nginx config from backup if removed
mv /etc/nginx/conf.d/arithi.conf.bak.TIMESTAMP /etc/nginx/conf.d/arithi.conf
nginx -t && systemctl reload nginx
```

To rollback an image version without stopping:

```bash
export EDAIRY_VERSION=1.0.0   # previous tag
docker compose pull && docker compose up -d
```

## Build images locally (fallback)

When images are built on the server (not pulled from GHCR), use `SKIP_PULL=true`:

```bash
sudo SKIP_PULL=true EDAIRY_VERSION=1.0.0 ./deploy/create-dairy.sh --migrate --skip-db arithi 8500 8501
```

```bash
# Backend
cd edairy-go
docker build -t ghcr.io/rubewafula/edairy-backend:1.0.0 .

# Frontend
cd edairy-react-admin/react-admin-app
docker build -t ghcr.io/rubewafula/edairy-frontend:1.0.0 .
```

## Health checks

```bash
curl http://127.0.0.1:8500/health          # {"status":"ok"}
curl http://127.0.0.1:8501/runtime-config.js
docker compose -f /opt/edairy/arithi/docker-compose.yml logs -f
```

## MySQL grants

Docker containers connect via `host.docker.internal`. The provisioning script grants:

```sql
GRANT ALL ON `{db}`.* TO 'edairy'@'172.%';
```

Verify from a throwaway container:

```bash
docker run --rm --add-host=host.docker.internal:host-gateway mysql:8 \
  mysql -h host.docker.internal -u edairy -p'...' -e 'SELECT 1'
```

## Templates

| File | Purpose |
|------|---------|
| [templates/docker-compose.yml](templates/docker-compose.yml) | Per-dairy stack |
| [templates/.env.template](templates/.env.template) | Backend environment |
| [templates/runtime-config.js.template](templates/runtime-config.js.template) | Frontend API URLs |
| [templates/nginx/default.conf](templates/nginx/default.conf) | SPA nginx config (baked into frontend image) |

Placeholders: `{{SLUG}}`, `{{APP_NAME}}`, `{{API_PORT}}`, `{{UI_PORT}}`, `{{DB_NAME}}`, `{{DB_PASSWORD}}`, `{{JWT_SECRET}}`.
