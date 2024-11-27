# Go Trading Ace

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
2. Configure the database, Redis, infura key to connect ethereum using `/config/config.yml`
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
      key: "your-key"
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
