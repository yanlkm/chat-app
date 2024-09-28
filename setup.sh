#!/bin/bash

echo "Setting up our environment..."
echo MONGO_URI=$MONGO_URI >> .env
echo PORT=:8080 >> .env
echo JWT_SECRET=$JWT_SECRET >> .env
echo GITHUB_CLIENT_SECRET=$GITHUB_CLIENT_SECRET >> .env
echo GITHUB_CLIENT_ID=$GITHUB_CLIENT_ID >> .env
echo BASE_URL=$BASE_URL >> .env
echo CORS_ORIGIN=$CORS_ORIGIN >> .env
echo "Environment setup complete."