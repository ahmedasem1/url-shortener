# URL Shortener

## Overview
This is a distributed URL shortener service built using **GoLang, Redis, MongoDB, and Kubernetes**. The service allows users to shorten URLs, track analytics, and enforce rate limiting.

## Features
- **Shorten URLs**: Generate short links for long URLs.
- **Redirect**: Redirect users to the original URL via the short link.
- **Rate Limiting**: Prevent abuse using Redis-based rate limiting.
- **Analytics Logging**: Track visits and log analytics asynchronously.
- **Distributed Deployment**: Scalable with Kubernetes.


## Installation & Setup

### **Prerequisites**
Ensure you have the following installed:
- **Docker & Docker Compose**

### **Clone the Repository**
```sh
git clone https://github.com/ahmedasem1/url-shortener.git
cd url-shortener
```

### **Environment Variables**
The service automatically sets default values if environment variables are not provided.However, you can create a `.env` file (optional) with custom values:
```ini
REDIS_ADDR=127.0.0.1:6379
MONGO_URI=mongodb://localhost:27017
SERVER_PORT=8080
RATE_LIMIT=10 # Requests per minute per IP
```

### **Run Locally (Docker Compose)**
```sh
docker-compose up --build
```

### **Run Tests**
```sh
docker compose up  test
```

---


## CI/CD Pipeline
This project includes a **GitHub Actions** CI/CD pipeline that automates deployment.

### **Build & Test**
- Runs on every push to `main` or `develop` branches and on pull requests.
- Sets up Go, installs dependencies, and runs tests.
- Starts MongoDB and Redis containers for integration testing.

### **Build & Push Docker Image**
- Runs only on pushes to `main`.
- Builds and pushes the Docker image to **Docker Hub**.

### **Deploy to Kubernetes**
- Runs only on pushes to `main`.
- Deploys the application using `kubectl apply`.

---

## API Endpoints

### **Swagger API Documentation**

After running the service, visit:

```
http://localhost:8080/swagger/index.html
```

to explore the available API endpoints.

### **1. Shorten URL**
**Endpoint:** `POST /shorten`
```json
{
  "url": "https://example.com"
}
```
**Response:**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

### **2. Redirect to Original URL**
**Endpoint:** `GET /:shortID`
- Redirects to the original URL.

### **3. Rate Limit Test**
**Endpoint:** `GET /test`
- Returns `200 OK` if under limit, `429 Too Many Requests` if exceeded.

---

## Kubernetes Deployment

### **Start Minikube**
```sh
minikube start --driver=docker
```

### **Enable Ingress (Optional, if using Ingress for routing)**
```sh
minikube addons enable ingress
```

### **Deploy to Kubernetes**
```sh
kubectl apply -f k8s/
```

### **Check if Pods are Running**
```sh
kubectl get pods
```

### **Get Service URL (If using LoadBalancer)**
```sh
minikube service list
```
Or, for a specific service:
```sh
minikube service api --url
```
