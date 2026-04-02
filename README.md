# ALRU URL Shortener

**Live Demo:** [https://alru-url-shortener.vercel.app/](https://alru-url-shortener.vercel.app/)

## Main Features

* **Zero-Collision Shortening:** Uses Base62 encoding mapped to auto-incrementing database IDs for mathematically unique short codes.
* **Custom Aliases:** Supports custom vanity URLs safely isolated in a separate `/c/:alias` routing namespace to prevent conflicts.
* **High-Speed Redirects:** Integrates Redis caching to bypass database lookups for frequently accessed links.
* **Asynchronous Analytics:** Captures click telemetry (IP hashing, devices, location) via background Go routines to ensure zero-latency redirects.
* **Interactive Dashboard:** Timezone-aware UI with Recharts for visualizing daily/hourly click trends, referrers, and geographic data.
* **Security & Privacy:** Features stateless JWT authentication, and SHA-256 IP hashing.

## Tech Stack

**Backend**
* Go (Golang)
* Echo v5 (Web Framework)
* PostgreSQL & GORM (Relational Database & ORM)
* Redis (In-Memory Cache)

**Frontend**
* React & Vite
* Tailwind CSS
* Recharts (Data Visualization)
* Lucide React (Icons)