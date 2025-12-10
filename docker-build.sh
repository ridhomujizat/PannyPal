#!/bin/bash

echo "🐳 Building Cash Flow API Docker Image with Swagger..."
echo ""

# Generate Swagger docs first
echo "📚 Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs

if [ $? -ne 0 ]; then
    echo "❌ Failed to generate Swagger docs!"
    exit 1
fi

echo "✅ Swagger docs generated!"

# Build Docker image
echo "🔧 Building Docker image..."
docker build -t pannypal-api:latest .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully!"
    echo ""
    echo "🚀 To run the container:"
    echo "   docker run -p 9001:9001 pannypal-api:latest"
    echo ""
    echo "🐳 Or use Docker Compose for full stack:"
    echo "   docker-compose up"
    echo ""
    echo "📚 Swagger will be available at:"
    echo "   http://localhost:9001/swagger/index.html"
else
    echo "❌ Docker build failed!"
    exit 1
fi