# Go Uniswap

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

### Installing and Running the Project

1. Clone the Repository First, clone the repository to your local machine.
    ```
    git clone https://github.com/tsen1220/go-trading-ace.git
    ```
2. Build and Run with Docker Compose
    ```
    docker-compose up --build
    ```
3. Run the Project Locally (Without Docker)
    If you prefer to run the project locally, follow these steps:
    - Download the necessary dependencies:  
        ```
        go mod tidy
        ```

    - Run the application:
        ```
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
go test ./tests
```
