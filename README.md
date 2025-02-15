# URL Shortener

## Overview
This is a distributed URL shortener service built using **GoLang, Redis, MongoDB, and Kubernetes**. The service allows users to shorten URLs, track analytics, and enforce rate limiting.

## Features
- **Shorten URLs**: Generate short links for long URLs.
- **Redirect**: Redirect users to the original URL via the short link.
- **Rate Limiting**: Prevent abuse using Redis-based rate limiting.
- **Analytics Logging**: Track visits and log analytics asynchronously.
- **Distributed Deployment**: Scalable with Kubernetes.

---

## Tech Stack
- **Backend**: Go (Gin framework)
- **Storage**: Redis (caching) + MongoDB (persistent storage)
- **Queue Processing**: Worker-based logging system
- **Deployment**: Kubernetes + Docker

---

## Installation & Setup

### **Prerequisites**
Ensure you have the following installed:
- **Go** (1.18+)
- **Docker & Docker Compose**
- **Kubernetes (Minikube or K3s for local testing)**
- **Redis**
- **MongoDB**

### **Clone the Repository**
```sh
git clone https://github.com/your-username/url-shortener.git
cd url-shortener
```

### **Environment Variables**
Create a `.env` file with the following variables:
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

## API Endpoints

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