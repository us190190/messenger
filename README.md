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
      ┌─────────────┐
      │   Clients   │
      └─────────────┘
             │
             ▼
      ┌─────────────┐
      │    Load     │
      │   Balancer  │
      └─────────────┘
             │
       ┌─────┼─────┐
       ▼     ▼     ▼
      ┌─────────────┐
      │  WebSocket  │
      │   Servers   │
      └─────────────┘
             │
       ┌─────┼─────┐
       ▼     ▼     ▼
      ┌─────────────┐
      │   Backend   │
      │  Services   │
      └─────────────┘
       │           │
       ▼           ▼
    ┌──────┐   ┌──────┐
    │  DB  │   │ Cache│
    └──────┘   └──────┘

## Low level design
1. TODO

## Flow diagram
1. TODO

## Steps to deploy backend service on GCP DM
   gcloud deployment-manager deployments create messenger-app-infra --config messenger-deployment.yaml

## Steps to deploy docker container on local
   - sudo apt update
   - sudo apt install -y docker.io
   - sudo systemctl start docker
   - sudo systemctl enable docker
   - docker pull us190190/messenger-app:3.0
   - docker run -e DB_HOST=XYZ -e DB_PASSWORD=XYZ -e DB_PORT=XYZ -e DB_USER=XYZ -p 443:443 messenger-app:3.0

## Steps to utilize postman collection
1. User management APIs collection: [User management APIs.postman_collection.json]
2. Messenger connection: [wss://<SERVER_IP_ADDRESS>:443/start]

## Known issues
1. Both live messages don't contain sender user information (sender is required to share his/her info along with each message)
2. Need to manually add VM instance public IP to allow connectivity with DB instance
