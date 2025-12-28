Project Title: URL Shortener

A fast, lightweight URL shortener built with Golang, Fiber, Redis, and Docker, featuring rate limiting, custom short URLs, and expiry support.


🚀 Features

Shorten long URLs

Custom short URL support

URL expiry (TTL)

Redis-backed storage

IP-based rate limiting

REST API

Docker & Docker Compose support


🛠 Tech Stack

Backend: Golang (Fiber)

Database: Redis

Containerization: Docker, Docker Compose

Validation: govalidator

Rate Limit: Redis-based


📁 Project Structure
.
├── api/
├── db/
├── database/
├── helpers/
├── routes/
├── main.go
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md

⚙️ Environment Variables

Create a .env file in the root:

APP_PORT=3000
DB_ADDR=db:6379
DB_PASS=
DOMAIN=localhost:3000
API_QUOTA=10

🐳 Run with Docker
docker-compose up --build


Server will run at:

http://localhost:3000

📡 API Endpoints
POST: /urlshortner/api/v1

Request Body

{
  "url": "https://toufiq.dev",
  "short": "",
  "expiry": 24
}


Response

{
  "url": "https://toufiq.dev",
  "short": "localhost:3000/abc123",
  "expiry": 24,
  "rate_limit": 9,
  "rate_limit_reset": 29
}


Redirects to the original URL.

⏱ Rate Limiting

Rate limit is applied per IP

Default quota: API_QUOTA

Auto reset after 30 minutes

✅ Status

✔ Fully working
✔ Dockerized
✔ Production-ready base

📌 Future Improvements

Authentication

Analytics (click count per URL)

Custom domain support

Admin dashboard

CI/CD pipeline

👤 Author

Md Toufiquzzaman
Software Engineer
