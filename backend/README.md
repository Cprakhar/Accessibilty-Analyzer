# Backend Environment Setup

## Prerequisites
- Go 1.20+ installed
- MongoDB running locally or accessible via URI
- (Optional) Redis if you plan to use rate limiting or background jobs

## Environment Variables
Create a `.env` file in the `backend/` directory (or set these in your shell):

```
MONGODB_URI=mongodb://localhost:27017/accessibility_analyser
PORT=8080
JWT_SECRET=your_secret_key_here
```

- `MONGODB_URI`: MongoDB connection string
- `PORT`: Port for the backend server (default: 8080)
- `JWT_SECRET`: Secret key for signing JWT tokens

## Install Go Dependencies
Run this in the `backend/` directory:

```
go mod init backend

go get github.com/gin-gonic/gin

go get go.mongodb.org/mongo-driver/mongo

go get github.com/golang-jwt/jwt/v4

go get golang.org/x/crypto/bcrypt
```

## Running the Server
From the `backend/` directory:

```
go run server.go
```

The server should start on `http://localhost:8080`.

## Notes
- Make sure MongoDB is running and accessible.
- Update `JWT_SECRET` in your code/config to use the value from the environment variable for better security.
- For production, use a process manager (like systemd or pm2) and serve over HTTPS.
