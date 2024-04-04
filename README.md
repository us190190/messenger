# Messenger Service
The service is a primitive platform for real-time communication between users. It provides an infrastructure for sending and receiving messages, supporting features such as text messaging, group chats, and more. Built with performance in mind, it ensures seamless communication across various devices and platforms.

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
      │     and     │
      │ Application │
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
    ┌────────┐ ┌────────┐
    │Database│ │ Cache  │
    └────────┘ └────────┘

## Low level design
![image](https://github.com/us190190/messenger/assets/3051295/cb0f4bb8-b909-44b3-8da2-107c903a5805)

## Key Features
1. User management: Register, Update, Authenticate, Remove a user
2. Real-time Messaging: Instantly send and receive messages in real-time.
3. Group Chats: Create and participate in group conversations with multiple users.
4. Notifications: Receive push notifications for new messages and updates.
5. Security: Ensure the privacy and security of messages with end-to-end encryption.

## Technologies Used

1. Backend: Go, MySQL
2. Frontend: WIP
3. Messaging Protocol: WebSocket
4. Deployment: Docker, GCPDM

## Steps to deploy backend service on GCP DM
```
   gcloud deployment-manager deployments create messenger-app-infra --config messenger-deployment.yaml
```

## Steps to deploy docker container on local
```
   sudo apt update
   sudo apt install -y docker.io
   sudo systemctl start docker
   sudo systemctl enable docker
   docker pull us190190/messenger-app:3.0
   docker run -e DB_HOST=XYZ -e DB_PASSWORD=XYZ -e DB_PORT=XYZ -e DB_USER=XYZ -p 443:443 messenger-app:3.0
```

## Steps to utilize postman collection
1. User management APIs collection: [User management APIs.postman_collection.json]
2. Messenger connection: [wss://<SERVER_IP_ADDRESS>:443/start]

## In-progress
1. Frontend application to consume these backend messenger service
2. Refactoring the repo
3. Adding database migrations and seeder for the backend service
4. Scalability: Scale the messaging service to accommodate growing user bases
5. Reliability: Ensure high availability and reliability with fault-tolerant architecture.
   
## Known issues
1. Both live messages don't contain sender user information (sender is required to share his/her info along with each message)
2. Need to manually add VM instance public IP to allow connectivity with DB instance
