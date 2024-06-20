
# Chat App

This is a chat application built with Go. It connects to a MongoDB database and provides functionality for user authentication, messaging, and room management.

## Features

- User authentication
- Room creation and management
- Real-time messaging
- MongoDB integration

## Prerequisites

- [Go](https://golang.org/dl/) 1.18 or later
- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

### Clone the repository

```sh
git clone https://gitlab.com/yanlkm/chat-app.git
cd chat-app
``` 

### Install dependencies

```sh
go mod tidy
```

### Build the application

```sh
go build -o chat-app
```

### Access the application

```sh
./chat-app
```

## Docker

### Build the Docker image

```sh
docker build -t yanlkm/chat-app:latest .
```

### Run the Docker container

```sh
docker run -d -p 8080:8080 --name chat-app --env-file .env yanlkm/chat-app:latest
```

## Docker Compose

### Create a `compose.yml` file:

```yaml
services:
  mongodb:
    image: yanlkm/mongo:latest
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db

  chat-app:
    image: yanlkm/chat-app:latest
    build:
      context: .
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
    volumes:
      - chat_app_data:/app/data

volumes:
  mongodb_data:
    driver: local
  chat_app_data:
    driver: local
```

### Run the Docker container:

```sh
docker-compose -f compose.yaml -p chat-app up -d 

```

## API Endpoints

### Authentication

- **POST /auth/login**: User login
- **GET /auth/logout**: User logout

### Users

- **GET /users**: Get all users
- **POST /users**: Create a user
- **GET /users/:id**: Get user by ID
- **PUT /users/:id**: Update user by ID
- **DELETE /users/:id**: Delete user by ID
- **PUT /users/:id/password**: Update password for a specific user
- **GET /users/ban/:id/:idBanned**: Ban a user from a room
- **GET /users/unban/:id/:idBanned**: Unban a user from a room

### Rooms

- **GET /rooms**: Get all rooms
- **POST /rooms**: Create a room
- **GET /rooms/:id**: Get room by ID
- **DELETE /rooms/:id**: Delete a room
- **GET /rooms/user/:id**: Get rooms of a user
- **PUT /rooms/add/:id**: Add a user to a room
- **PUT /rooms/remove/:id**: Remove a user from a room
- **PATCH /rooms/add/hashtag/:id**: Add a hashtag to a room
- **PATCH /rooms/remove/hashtag/:id**: Remove a hashtag from a room
- **GET /rooms/members/:id**: Get members of a room

### Messages

- **POST /messages**: Send a new message in a room
- **GET /messages/{id}**: Get messages of a specific room
- **DELETE /messages/{id}**: Delete a message in a room

### Codes

- **POST /codes**: Create an authentication code (admin only)

## Author

Yan [yanlkm](https://github.com/yanlkm)
