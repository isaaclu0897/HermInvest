# HermInvest

Hermes, the guardian deity of commerce in ancient Greek mythology. Merging Hermes with the power of investment, a platform that symbolizes reliability and efficient wealth creation.

A lightweight investment management platform built with Go, featuring both CLI and web interfaces for managing and visualizing trading data.

- [Key Features](#key-features)
- [Preview](#preview)
  - [Web Usage](#web-ssage)
  - [CLI Usage](#cli-usage)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Installation and Usage](#installation-and-usage)

## Key Features

- **Dual Interface**: CLI for operations, Web UI for visualization
- **Flexible Data Ingestion**: Supports multiple data sources with normalization (commission history, manual input, stock mapping)
- **Lightweight & Portable**: SQLite-based design with minimal setup
- **Built-in Visualization**: Trading insights powered by Chart.js

## Preview

### Web Usage

```bash
./hermInvestCli stock web
```

<img src="https://meee.com.tw/gPZFWDs.png" width="500">

### CLI Usage

```bash
./hermInvestCli query --all
./hermInvestCli stock add 2023-12-01 09:00:00 0050 1 1500 23.5
```
<img src="https://meee.com.tw/LuS66rQ.png" width="400"> <img src="https://meee.com.tw/B1gn8WZ.png" width="360">


## Architecture

This project follows a three-tier architecture: Presentation, Service, and Repository layers, with Model as the connecting entity across layers.

```
.
├── cmd/
│   └── hermInvestCli/  # CLI entry (Cobra) / Web UI (HTML + Chart.js)
├── pkg/
│   ├── service/        # Service:      Business logic
│   ├── repository/     # Repository:   Data access
│   └── model/          # Models:       Data structures shared across layers
└── internal/           # Private app logic
```

- **Presentation Layer**: `cmd/` (CLI) and (Web UI) handle user interactions
- **Service Layer**: `service/` implements business logic and orchestrates operations
- **Repository Layer**: `repository/` manages data persistence and queries
- **Model**: `model/` defines data structures used by all layers

The system is designed with minimal dependencies and avoids heavy frameworks to maintain simplicity, transparency, and control.

## Tech Stack

- **Backend**: Go (Gin, GORM)
- **Frontend**: HTML, Bootstrap, jQuery, Chart.js
- **Database**: SQLite
- **CLI**: Cobra
- **Build Tool**: Makefile


## Installation and Usage

1. Ensure Go 1.20+ is installed
2. Clone the repository: `git clone <repo>`
3. Install dependencies: `go mod tidy`
4. Initialize database: `go run ./internal/createDBSchema/createDBSchema.go`
5. Seed sample data (optional): `go run ./internal/seedSampleData/seedSampleData.go`
6. Build: `make build`
7. Run CLI: `./hermInvestCli --help`
8. Start Web: `./hermInvestCli stock web`
