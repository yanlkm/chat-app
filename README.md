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
git clone https://github.com/yanlkm/chat-app.git
cd chat-app
``` 

## Install dependencies

```sh
go mod tidy
```

## Build the application

```sh
go build -o chat-app
```
## Access the application
### localhost:8080
```sh
./chat-app
```

# Docker

## Build the Docker image

```sh
docker build -t yanlkm/chat-app:latest .
```
## Run the Docker container

```sh
docker run -d -p 8080:8080 --name chat-app --env-file .env yanlkm/chat-app:latest
```

# Docker compose

## Run the Docker container

###  Create a `compose.yml` file:

```sh
version: '3.8'

services:
  mongodb:
    image: mongo:latest
    container_name: mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongo-data:/data/db

  chat-app:
    image: yanlkm/chat-app:latest
    container_name: chat-app
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - mongodb

volumes:
  mongo-data:
```

### Run the Docker container:

```sh
docker-compose up -d
```
# API Endpoints
## Authentication

    POST /api/auth/register: Register a new user
    POST /api/auth/login: Login an existing user

## Users

    GET /api/users: Get all users
    GET /api/users/:id: Get user by ID

## Rooms

    POST /api/rooms: Create a new room
    GET /api/rooms: Get all rooms
    GET /api/rooms/:id: Get room by ID

## Messages

    POST /api/messages: Send a message
    GET /api/messages: Get all messages
    GET /api/messages/:room_id: Get messages by room ID
## Authors

 Yan [yanlkm](https://github.com/yanlkm)