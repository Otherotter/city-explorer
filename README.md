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
