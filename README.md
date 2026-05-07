# City Explorer

> If I moved to a new city tomorrow, what would
> my ideal routine look like?

A personal city exploration tool that aggregates
data from public APIs, serves it through a custom
Go API, and monitors the entire system with
Grafana Cloud observability.

---

## What It Does

Search any city and explore it by category:

- 🍜 **Food & Drink** — restaurants, cafes, markets
- 🌿 **Nature** — parks, trails, wildlife
- 📚 **Study Spots** — libraries, cafes, universities
- 🎉 **Events** — local happenings
- 🏛️ **Attractions** — museums, landmarks
- 🛍️ **Thrift** — secondhand shops, vintage stores
- 🤝 **Social** — meetups, community groups
- 🏗️ **Architecture** — notable buildings

Data is fetched on demand from OpenStreetMap
via the Overpass API and stored in PostgreSQL.
First request per city/category triggers
collection. Subsequent requests serve from DB.

---

## Architecture
Browser (Next.js)
│
▼
API Service (Go + Chi)    :8080
│
├── checks PostgreSQL for fresh data
│
└── triggers Collector if stale
│
▼
Collector Service (Go)    :8081
│
└── fetches from Overpass API
│
▼
PostgreSQL + PostGIS



---
## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend API | Go + Chi |
| Data Collection | Go (Collector service) |
| Database | PostgreSQL 16 + PostGIS |
| Cache | Redis (planned) |
| Frontend | Next.js + Tailwind + Leaflet |
| Observability | Grafana Cloud (Loki + Prometheus + Tempo) |
| Agent | Grafana Alloy |
| Container | Docker + Docker Compose |
| Deploy | Fly.io (planned) |

---

## Observability

This project treats observability as a
first-class concern, not an afterthought.

### Three Pillars

**Metrics → Prometheus**
> cityexplorer_collector_overpass_requests_total
> cityexplorer_collector_overpass_duration_seconds
> cityexplorer_collector_places_inserted_total
> cityexplorer_collector_errors_total
> pg_up, pg_database_size_bytes, pg_stat_*



----
**Logs → Loki**

{
  "level": "INFO",
  "msg": "collection complete",
  "service": "collector",
  "project": "city-explorer",
  "city": "new york city",
  "category": "food",
  "fetched": 18226,
  "inserted": 18226,
  "duration_ms": 9802
}



-------
**Traces → Tempo**

GET /cities/new york city/food [38ms]
├── db.cities.lookup [2ms]
├── db.places.check_status [1ms]
├── collector.trigger [828ms]
│     └── overpass.http_request [15s]
│           status_code = 200
│           elements.returned = 18521
└── db.places.fetch_by_category [14ms]
Shared Observability Library
All services use a shared observability package
that enforces consistent field names, log format,
and trace setup across the entire system.
Copyshared/observability/
├── logger.go      NewLogger(serviceName)
├── middleware.go  RequestLogger (structured HTTP logs)
├── metrics.go     NewCollectorMetrics()
└── tracer.go      InitTracer(ctx, serviceName)
Every log line across every service includes:
jsonCopy{
  "service": "api",
  "project": "city-explorer",
  ...
}

##Project Structure
```
city-explorer/
├── services/
│   ├── api/                    Go API server
│   │   ├── internal/
│   │   │   ├── collector/      Collector HTTP client
│   │   │   ├── db/             Database connection + queries
│   │   │   └── handlers/       HTTP handlers
│   │   ├── Dockerfile
│   │   └── main.go
│   └── collector/              Data collection service
│       ├── internal/
│       │   ├── config/         Category → amenity mapping
│       │   ├── db/             Database writes
│       │   ├── fetchers/
│       │   │   └── overpass/   Overpass API client
│       │   └── models/         Internal data models
│       ├── Dockerfile
│       └── main.go
├── shared/
│   └── observability/          Shared o11y library
│       ├── logger.go
│       ├── middleware.go
│       ├── metrics.go
│       └── tracer.go
├── frontend/                   Next.js app
│   ├── app/
│   │   ├── page.tsx            Home / city search
│   │   └── cities/[name]/      City view
│   ├── components/
│   │   ├── CitySearch.tsx
│   │   ├── CityView.tsx
│   │   ├── CategoryTabs.tsx
│   │   ├── PlaceMap.tsx
│   │   └── PlaceCard.tsx
│   └── lib/
│       └── api.ts              API client
├── db/
│   ├── migrations/             SQL schema files
│   ├── postgres/               PostgreSQL config
│   └── setupDB.sh              Database setup script
├── alloy/
│   └── config.alloy            Grafana Alloy config
├── k8s/                        Kubernetes manifests (planned)
├── docker-compose.yml
├── go.mod
└── .env
```

##Getting Started
###Prerequisites

Docker + Docker Compose
Go 1.26+
Node.js 20+
1. Clone and configure
bashCopygit clone https://github.com/Otherotter/city-explorer.git
cd city-explorer

Create .env from the example:
cp .env.example .env

```
*Fill in your values:*
bashCopy# Database
POSTGRES_USER=cityexplorer
POSTGRES_PASSWORD=yourpassword
POSTGRES_DB=city_explorer
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
POSTGRES_DSN=postgresql://cityexplorer:yourpassword@db:5432/city_explorer?sslmode=disable

# Monitoring user (read-only)
POSTGRES_MONITOR_DSN=postgresql://grafana_monitor:monitorpassword@db:5432/city_explorer?sslmode=disable

# Grafana Cloud
GRAFANA_PROMETHEUS_URL=https://prometheus-prod-xx.grafana.net/api/prom/push
GRAFANA_PROMETHEUS_USERNAME=your_username
GRAFANA_LOKI_URL=https://logs-prod-xx.grafana.net/loki/api/v1/push
GRAFANA_LOKI_USERNAME=your_username
GRAFANA_TEMPO_URL=https://tempo-prod-xx.grafana.net/tempo
GRAFANA_TEMPO_USERNAME=your_username
GRAFANA_API_KEY=your_api_key

# Collector
COLLECTOR_URL=http://collector:8081
```

# Tracing
```
OTEL_EXPORTER_OTLP_ENDPOINT=alloy:4318
2. Start the database
bashCopydocker compose up db -d
3. Run database setup
bashCopy./db/setupDB.sh
4. Start all services
bashCopydocker compose up --build
5. Start the frontend
bashCopycd frontend
npm install
npm run dev
6. Open the app
http://localhost:3000
```

API Endpoints
```
MethodEndpointDescriptionGET/healthService health checkGET/metricsPrometheus metricsGET/cities/{name}/foodFood spotsGET/cities/{name}/natureNature spotsGET/cities/{name}/studyStudy spotsGET/cities/{name}/eventsEventsGET/cities/{name}/attractionAttractionsGET/cities/{name}/thriftThrift shopsGET/cities/{name}/socialSocial groupsGET/cities/{name}/architectureArchitecturePOST/personal-logLog a visitGET/personal-log/{city}Your visits
Example
```
```
curl http://localhost:8080/cities/new%20york%20city/food
```

```
jsonCopy{
  "city": "new york city",
  "count": 100,
  "places": [
    {
      "id": 1,
      "name": "Joe's Pizza",
      "subcategory": "pizza",
      "address": "7 Carmine St",
      "latitude": 40.7303,
      "longitude": -74.0033,
      "website": "https://joespizzanyc.com",
      "phone": "+1 212-366-1182",
      "source": "overpass"
    }
  ]
}
```
Data Sources

> SourceCategoriesLicenseOpenStreetMap via OverpassAllODbLOpenTripMapAttractionsCC BYNational Park ServiceNaturePublic DomainEventbriteEventsAPI TermsMeetupSocialAPI Terms

Database Schema
Copycities ──< places ──< place_tags >── tags
cities ──< events
cities ──< social_groups
cities ──< personal_log
places ──< place_history ──< history_tags
categories ──< places
personal_log ──< personal_log_tags >── tags

Observability Phases
This project documents a deliberate two-phase
observability journey.
Phase 1 — Inconsistent (documented, not fixed)
Each service logged with different field names.
Loki queries failed across service boundaries.
Preserved as a learning artifact.
CopyAPI:        city, category, duration_ms
Collector:  city_name, cat, duration
Phase 2 — Consistent (current)
Shared observability library enforces
consistent fields across all services.
CopyBoth:       city, category, duration_ms
One Loki query covers the entire system:
Copy{project="city-explorer"} | json | city="new york city"

Version History
TagDescriptionv0.1Phase 1 observability — intentional inconsistencyv0.2Shared observability libraryv0.3Structured request loggingv0.4Custom Prometheus metricsv0.5Overpass call instrumented with trace spanv0.6Database observabilityv0.7Frontend — Next.js with Leaflet map

Known Issues
CopyOnly food category returns data
  Other API routes not built yet.

Neighborhood field empty
  Overpass does not return neighborhood data.
  Needs reverse geocoding enrichment.

No pagination
  API hardcoded to LIMIT 100.

CORS only allows localhost:3000
  Will need updating for production.

database_observability.postgres not working
  Alloy v1.12.1 cannot parse Debian Postgres
  version string. Blocked by upstream bug.

Roadmap
⬜ Remaining API routes (nature, study, events...)
⬜ Personal log entry form
⬜ Redis caching layer
⬜ Pagination
⬜ Minikube deployment
⬜ Kubernetes manifests
⬜ Beyla no-code instrumentation (K8s phase)
⬜ Fly.io production deployment

Engineering Notes
Detailed decision logs, debugging sessions,
and architectural decisions are documented in
NOTES.md.
This project was built to demonstrate:

Production-quality Go service architecture
Observability as a first-class concern
The evolution from inconsistent to consistent
observability across a distributed system
On-demand data collection patterns
Monorepo structure for Go microservices


License
MIT
