Advertising Service

This project is a REST API for posting and retrieving classified ads, built with Go (Golang) and following Clean Architecture principles.

🚀 Features

Create new advertisements with title, description, photo URLs, and price.

Get ad by ID (basic fields or full info via fields query parameter).

List ads with pagination (10 items per page) and sorting by price or creation date (ascending/descending).

Configuration via .env and config.yaml using Viper.

Graceful shutdown support.

Fully containerized with Docker and Docker Compose.

Unit testing and linting support with Makefile.

🗂 Project Structure

├── cmd                  # Entry point (main.go)    
├── pkg                  # Clean Architecture layers    
│   ├── handler          # HTTP request handlers    
│   ├── service          # Business logic   
│   ├── repository       # Database access via sqlx     
│   ├── model            # Domain models    
│   └── error_message    # Error definitions        
├── configs              # Configuration loaders (.env and config.yaml) 
├── migrations           # SQL schema migrations    
├── scripts              # SQL test data scripts    
├── test                 # Test database initialization     
├── e2e_test.go          # End-to-end test placeholder  
├── Dockerfile           # Multi-stage Docker build     
├── docker-compose.yml   # Docker Compose setup (app + db)  
├── Makefile             # Build, lint, test, and Docker shortcuts  
├── .env                 # Environment variables (DATABASE_URL, PORT)   
└── README.md            # Project documentation    

🛠 Technologies Used

Go 1.21+

Echo framework

PostgreSQL

sqlx

Viper

Testify and mock

golangci-lint

Docker & Docker Compose

Make

📦 Getting Started

Prerequisites

Docker & Docker Compose

Go 1.21+ (for local builds)

Make

Run with Docker Compose

docker-compose up --build

Your application will be available at http://localhost:8080.

Run Tests and Lint

make test
make lint

Build and Run Locally

make build
./advertising

🧑‍💻 Author

Alexandr Petrenko

✍️ This project was built as a technical test assignment demonstrating backend development skills, clean code, testing, and containerization.

