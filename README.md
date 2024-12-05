# Go Trading Ace

## Introduction
Welcome to Go Trading Ace – a Go-based platform that I built to interact with decentralized finance (DeFi) protocols. The goal of this project is to leverage Go to fetch data from the Uniswap smart contracts and create a point-based rewards system for liquidity providers, turning liquidity provision into an engaging incentive activity.

In this system, I extract data from Uniswap’s smart contracts, focusing on liquidity pools and trading volumes. Based on the liquidity provided by users, I calculate rewards and structure them into a "share pool" system. Users earn points for their participation, which are stored and tracked, creating a dynamic incentive system to encourage liquidity provision.

Key Features:
- Uniswap Integration: The system interacts with Uniswap's Ethereum-based smart contracts to fetch real-time data on liquidity pools and trading volumes.
- Share Pool Rewards: A point-based reward system where users earn points based on the liquidity they provide. These points are recorded in a decentralized ledger for transparency.
- Leaderboard System: Users can see where they rank on the leaderboard, showcasing the top contributors to liquidity pools based on their accumulated rewards.

With Go Trading Ace, I’ve built a platform where developers and DeFi enthusiasts can easily integrate Uniswap data into their applications, develop liquidity incentives, and create engaging user experiences that promote participation in DeFi protocols.

## Getting Started

These instructions will help you set up the project on your local machine for development and testing purposes.

### Prerequisites

Before getting started, ensure that you have the following software installed on your machine:

- **Go** (recommended version 1.18 or higher). You can download Go from [Go's official website](https://golang.org/dl/).
- **Docker** and **Docker Compose** for running the application in containers. You can download these from:
  - [Docker](https://www.docker.com/get-started)
  - [Docker Compose](https://docs.docker.com/compose/install/)
- **migrate** (for handling database migrations). You can install it using:
  ```
  brew install golang-migrate
  ```

  or you can manually install it:
  ```
  curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz
  sudo mv migrate /usr/local/bin
  ```
- **Infura Key**
  You can get Infura key from [Infura](https://www.infura.io/).

### Installing and Running the Project

1. Clone the Repository First, clone the repository to your local machine.
    ```
    git clone https://github.com/tsen1220/go-trading-ace.git
    ```
2. Configure the database, Redis, infura key to connect ethereum using `/config/config.yml`. And here is example key.
    ```
    database:
      host: "postgres"
      user: "root"
      password: "root"
      port: 5432
      name: "trading-ace"
      sslmode: "disable"

    redis:
      prefix: "trading-ace:"
      host: "redis"
      port: 6379
    
    infura:
      key: "<your-infura-project-id>"
    ```
3. Build and Run with Docker Compose
    ```
    docker-compose up --build
    ```
4. Run the Project Locally (Without Docker)  
   If you prefer to run the project locally, follow these steps:  

   - Download the necessary dependencies:  
     ```bash
     go mod tidy
     ```  

   - Run the application:  
     ```bash
     go run main.go
     ```

### Database Migration

1. **Configure Database Connection**
   
    Ensure that your database connection details are correctly set in the config.yml file or your environment variables.
2. **Run Migrations**
    Run the database migrations to set up the schema in the database:
    ```
    migrate -path=migrations -database "postgres://root:root@localhost:5432/trading-ace?sslmode=disable" up
    ```
3. **Rollback Migrations**
    If you need to rollback the migrations, you can use:
    ```
    migrate -path=migrations -database "postgres://root:root@localhost:5432/trading-ace?sslmode=disable" down
    ```

### Running Tests
To run the tests for the project, execute the following command:
```
go test ./...
```

If you want to get coverage info:
```
go test ./... -cover
```
