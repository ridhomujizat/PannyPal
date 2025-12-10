#!/bin/bash

echo "🚀 Starting Cash Flow API with Swagger Documentation..."
echo ""
echo "📊 API Endpoints:"
echo "   - Health Check: http://localhost:9001/health"
echo "   - Swagger UI:   http://localhost:9001/swagger/index.html"
echo "   - Swagger JSON: http://localhost:9001/swagger/doc.json"
echo "   - API Base:     http://localhost:9001/api"
echo ""
echo "🔧 Choose deployment method:"
echo "   1) Local build and run"
echo "   2) Docker build and run"
echo "   3) Docker Compose (with PostgreSQL, Redis, RabbitMQ)"
echo ""

read -p "Enter choice [1-3]: " choice

case $choice in
  1)
    echo "🔧 Building application locally..."
    go build -o pannypal-api ./cmd/api
    
    if [ $? -eq 0 ]; then
        echo "✅ Build successful!"
        echo ""
        echo "🎯 Starting application..."
        ./pannypal-api
    else
        echo "❌ Build failed!"
        exit 1
    fi
    ;;
  2)
    echo "🐳 Building Docker image..."
    docker build -t pannypal-api .
    
    if [ $? -eq 0 ]; then
        echo "✅ Docker build successful!"
        echo ""
        echo "🎯 Starting container..."
        docker run -p 9001:9001 pannypal-api
    else
        echo "❌ Docker build failed!"
        exit 1
    fi
    ;;
  3)
    echo "🐳 Starting with Docker Compose..."
    echo "📝 This will start:"
    echo "   - Cash Flow API (port 9001)"
    echo "   - PostgreSQL (port 5432)"
    echo "   - Redis (port 6379)"
    echo "   - RabbitMQ (port 5672, management: 15672)"
    echo ""
    
    docker-compose up --build
    ;;
  *)
    echo "❌ Invalid choice!"
    exit 1
    ;;
esac

echo ""
echo "📚 Access Swagger Documentation:"
echo "   http://localhost:9001/swagger/index.html"