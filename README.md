# Gatus Status

Clean and simple **public-facing status page** built on top of a [Gatus](https://github.com/TwiN/gatus) instance.

Gatus' own UI (whilst great) shows everything, this app uses the API to a Gatus instance and allows us to cherrypick specific endpoints from this and show a slimmed down public facing page which can be exposed to all, whilst Gatus itself remains protected.

- Cherrypick Groups/Endpoints from a Master Gatus Instance.
- Strip 'sensitive' data such as errors and condition results before anything reaches the browser.

## Configuration

Configure with a YAML file (default `config/config.yml`, override with `GATUS_STATUS_CONFIG_PATH`):

```yaml
title: Status

gatus:
  url: http://gatus:8080
  # Optional Basic Auth, if the Gatus API is protected
  # username: ${GATUS_USERNAME}
  # password: ${GATUS_PASSWORD}

# Endpoints to show, in display order.
endpoints:
  - frontend        # endpoint name
  - core/api        # group/name
  - core_database   # Gatus key

# Show the announcements configured in Gatus (default true)
announcements: true
```

Each entry is matched against the Gatus endpoint **name**, **`group/name`**, or **key**. Use `group/name` or the key if names are duplicated across groups. Endpoints not found in Gatus are shown as "Unknown".

Or configure entirely with environment variables (these override the file):

| Variable | Description |
| --- | --- |
| `GATUS_URL` | Base URL of the Gatus instance |
| `ENDPOINTS` | Comma-separated list of endpoints, e.g. `frontend,core/api` |
| `TITLE` | Page title (default `Status`) |
| `GATUS_USERNAME` / `GATUS_PASSWORD` | Optional Basic Auth for the Gatus API |
| `ANNOUNCEMENTS` | Show the announcements configured in Gatus, `true` or `false` (default `true`) |

Statuses are fetched from Gatus' `/api/v1/endpoints/statuses` and `/api/v1/endpoints/{key}/uptimes/{duration}` APIs, cached for 30s, and the page refreshes every minute.

### Announcements

[Gatus announcements](https://github.com/TwiN/gatus#announcements) are read from `/api/v1/config` and shown as-is: active ones above the overall status, archived ones under "Past Announcements". They're the place to tell visitors about planned maintenance, as Gatus' maintenance windows only silence alerting - checks that fail during a window still show as down.

## Running

```sh
docker compose up -d
```

### Development

```sh
go run .                       # API on :8080
cd frontend && npm run dev     # Vite dev server, proxies /api to :8080
```
