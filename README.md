# Messenger [Service]

## Features
1. User management - Register, Update, Authenticate, Remove
2. Messaging supported:
   1. personal messages between any two users
   2. group messages among members of the group
3. If a user is connected, then they receive personal or group messages immediately
4. Offline users get both type of messages, when they come online
5. Messages are end-to-end encrypted using SSL/TLS encryption

## In-progress
1. Frontend application to consume these backend messenger service
2. Refactoring the repo
3. Adding database migrations and seeder for the backend service

## Architecture diagram
1.

## Low level design
1.

## Flow diagram
1.

## Steps to deploy backend service on GCP
1.

## Steps to deploy docker container on local
1.

## Steps to utilize postman collection
1.

## Known issues
1. Both live messages don't contain sender user information (sender is required to share his/her info along with each message)
2. Need to manually add VM instance public IP to allow connectivity with DB instance

TODO:
# GCP infra launch config
1. config.yaml in DockerFile
2. commit files in git repo, add steps in READme
