# marketflow
Real-Time Cryptocurrency Market Data Processing System using Go, Redis, and PostgreSQL with support for concurrency patterns and hexagonal architecture.

---

## 🧭 Table of Contents

* [🚀 Launch the Project](#launch-the-project)
* [📁 Project Structure](#project-structure)
* [🌐 API Endpoints](#api-endpoints)

---

## 🚀 Launch the Project

Make sure you have Docker installed, then run:

```bash
docker load -i exchange1_amd64.tar

docker load -i exchange2_amd64.tar

docker load -i exchange3_amd64.tar
```

```bash
docker compose up --build
```

🔧 Don’t forget to configure `.env` properly before launch.

---

## 📁 Project Structure

```bash
.
├── cmd
│   └── main.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── driven
│   │   │   ├── exchange
│   │   │   │   ├── infra.go
│   │   │   │   ├── ping.go
│   │   │   │   └── subscribe.go
│   │   │   ├── generator
│   │   │   │   ├── generator.go
│   │   │   │   └── infra.go
│   │   │   ├── postgres
│   │   │   │   ├── averageDuration.go
│   │   │   │   ├── average.go
│   │   │   │   ├── eskerim
│   │   │   │   ├── fallback.go
│   │   │   │   ├── highestDuration.go
│   │   │   │   ├── highest.go
│   │   │   │   ├── infra.go
│   │   │   │   ├── latest.go
│   │   │   │   ├── lowestDuration.go
│   │   │   │   ├── lowest.go
│   │   │   │   ├── migration
│   │   │   │   │   └── init.sql
│   │   │   │   ├── saver.go
│   │   │   │   └── testcleaner.go
│   │   │   └── redis
│   │   │       ├── add.go
│   │   │       ├── aggregator.go
│   │   │       ├── allAvg.go
│   │   │       ├── average.go
│   │   │       ├── eskerim
│   │   │       ├── getExchanges.go
│   │   │       ├── highest.go
│   │   │       ├── infra.go
│   │   │       ├── latest.go
│   │   │       ├── lowest.go
│   │   │       ├── testCleaner.go
│   │   │       └── tsParser.go
│   │   └── driver
│   │       └── http
│   │           ├── api
│   │           │   ├── handler
│   │           │   │   ├── average.go
│   │           │   │   ├── eskerim
│   │           │   │   ├── health.go
│   │           │   │   ├── highest.go
│   │           │   │   ├── infra.go
│   │           │   │   ├── latest.go
│   │           │   │   ├── lowest.go
│   │           │   │   ├── mode.go
│   │           │   │   └── utils.go
│   │           │   └── router.go
│   │           └── server.go
│   ├── app
│   │   ├── app.go
│   │   └── stream.go
│   ├── config
│   │   ├── config.go
│   │   ├── psql.go
│   │   ├── redis.go
│   │   ├── server.go
│   │   ├── sources.go
│   │   ├── utils.go
│   │   └── worker.go
│   ├── domain
│   │   ├── errors.go
│   │   └── exchange.go
│   ├── ports
│   │   ├── eskerim
│   │   ├── inbound
│   │   │   ├── app.go
│   │   │   └── usecase.go
│   │   └── outbound
│   │       ├── generator.go
│   │       ├── postgres.go
│   │       ├── redis.go
│   │       ├── server.go
│   │       └── stream.go
│   └── services
│       ├── one
│       │   ├── average.go
│       │   ├── barcherGo.go
│       │   ├── infra.go
│       │   ├── interface.go
│       │   ├── outbound.go
│       │   └── timerGo.go
│       ├── streams
│       │   ├── infra.go
│       │   ├── interface.go
│       │   └── test.go
│       ├── syncPool
│       │   ├── infra.go
│       │   └── interface.go
│       ├── usecase
│       │   ├── average.go
│       │   ├── health.go
│       │   ├── highest.go
│       │   ├── infra.go
│       │   ├── latest.go
│       │   ├── lowest.go
│       │   └── mode.go
│       └── workers
│           ├── interface.go
│           ├── pool.go
│           └── worker.go
├── Makefile
├── readme.md
└── todo

26 directories, 89 files
```

🧱 **Architecture:** Clean Hexagonal (Ports & Adapters)<br>
⚙️ **Patterns Used:** Fan-in, Fan-out, Worker Pool, Generator

---

## 🌐 API Endpoints

### 📊 Price Endpoints

| Method | Endpoint                                                |
| ------ | ------------------------------------------------------- |
| GET    | `/prices/latest/{symbol}`                               |
| GET    | `/prices/latest/{exchange}/{symbol}`                    |
| GET    | `/prices/highest/{symbol}`                              |
| GET    | `/prices/highest/{exchange}/{symbol}`                   |
| GET    | `/prices/highest/{symbol}?period={duration}`            |
| GET    | `/prices/highest/{exchange}/{symbol}?period={duration}` |
| GET    | `/prices/lowest/{symbol}`                               |
| GET    | `/prices/lowest/{exchange}/{symbol}`                    |
| GET    | `/prices/lowest/{symbol}?period={duration}`             |
| GET    | `/prices/lowest/{exchange}/{symbol}?period={duration}`  |
| GET    | `/prices/average/{symbol}`                              |
| GET    | `/prices/average/{exchange}/{symbol}`                   |
| GET    | `/prices/average/{symbol}?period={duration}`            |
| GET    | `/prices/average/{exchange}/{symbol}?period={duration}` |

### 🔁 Mode Switching

| Method | Endpoint                               |
| ------ | -------------------------------------- |
| POST   | `/mode/live` – Switch to **Live Mode** |
| POST   | `/mode/test` – Switch to **Test Mode** |

### 🩺 Health Check

| Method | Endpoint                          |
| ------ | --------------------------------- |
| GET    | `/health` – Returns system status |

---