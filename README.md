# HermInvest

Hermes, the guardian deity of commerce in ancient Greek mythology. Merging Hermes with the power of investment, a platform that symbolizes reliability and efficient wealth creation.

A lightweight investment management platform built with Go, featuring both CLI and web interfaces for managing and visualizing trading data.

- [Key Features](#key-features)
- [Preview](#preview)
  - [Web Dashboard](#web-dashboard)
  - [CLI Usage](#cli-usage)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Installation and Usage](#installation-and-usage)

## Key Features

- **Dual Interface**: CLI for automation and batch operations, Web UI for visualization
- **Flexible Data Ingestion**: Supports multiple data sources with normalization (commission history, manual input, stock mapping)
- **Lightweight & Portable**: SQLite-based design with minimal setup
- **Built-in Visualization**: Trading insights powered by Chart.js

## Preview

### Web Dashboard
<img src="https://meee.com.tw/gPZFWDs.png" width="500">

### CLI Usage

<img src="https://meee.com.tw/LuS66rQ.png" width="400"> <img src="https://meee.com.tw/B1gn8WZ.png" width="360">



```bash
./hermInvestCli add --symbol AAPL --price 180 --qty 10
./hermInvestCli query --all
```

## Architecture

This project follows a lightweight three-tier architecture, separating interface, business logic, and data access layers while minimizing external dependencies.

```
[ Interface Layer ]
- CLI (Cobra)
- Web (Go Gin + HTML + Chart.js)
        ↓
[ Service Layer ]
- Business logic
- Data processing & normalization
        ↓
[ Data Layer ]
- GORM ORM
- SQLite
```

- **Interface Layer**: Provides both CLI and Web interfaces for different usage scenarios
- **Service Layer**: Centralized business logic, shared by CLI and Web
- **Data Layer**: Handles persistence using GORM with SQLite backend

The system is designed with minimal dependencies and avoids heavy frameworks to maintain simplicity, transparency, and control.

## Tech Stack

- **Backend**: Go (Gin, GORM)
- **Database**: SQLite
- **Frontend**: HTML, Bootstrap, jQuery, Chart.js
- **CLI**: Cobra
- **Build Tool**: Makefile
- **Scripts**: Shell scripts for data conversion


## Installation and Usage

1. Ensure Go 1.20+ is installed
2. Clone the repository: `git clone <repo>`
3. Install dependencies: `go mod tidy`
4. Initialize database: `go run ./internal/createDBSchema/createDBSchema.go`
5. Seed sample data (optional): `go run ./internal/seedSampleData/seedSampleData.go`
6. Build: `make build`
7. Run CLI: `./hermInvestCli --help`
8. Start Web: `./hermInvestCli stock web`
