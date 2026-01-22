# Grafana Monitoring Setup Script
# Starts the complete observability stack

Write-Host "🚀 Starting Grade Management Monitoring Stack..." -ForegroundColor Green
Write-Host ""

# Check if Docker is running
try {
    docker ps | Out-Null
    Write-Host "✅ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker Desktop." -ForegroundColor Red
    exit 1
}

# Start services
Write-Host ""
Write-Host "📦 Starting services with docker-compose..." -ForegroundColor Cyan
docker-compose up -d

# Wait for services to be healthy
Write-Host ""
Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Check service status
Write-Host ""
Write-Host "📊 Service Status:" -ForegroundColor Cyan
docker-compose ps

# Service URLs
Write-Host ""
Write-Host "🎉 Monitoring stack is ready!" -ForegroundColor Green
Write-Host ""
Write-Host "📍 Access Points:" -ForegroundColor Cyan
Write-Host "   🔴 Redis:       redis://localhost:6379"
Write-Host "   🟠 Prometheus:  http://localhost:9090"
Write-Host "   🟢 Grafana:     http://localhost:3000"
Write-Host "   🔵 Loki:        http://localhost:3100"
Write-Host ""
Write-Host "🔑 Grafana Credentials:" -ForegroundColor Yellow
Write-Host "   Username: admin"
Write-Host "   Password: admin"
Write-Host ""
Write-Host "📈 Pre-configured Dashboard:" -ForegroundColor Cyan
Write-Host "   Grade Calculation Engine - Performance Dashboard"
Write-Host ""
Write-Host "📝 Next Steps:" -ForegroundColor Magenta
Write-Host "   1. Start your Go API: go run main.go"
Write-Host "   2. Open Grafana: http://localhost:3000"
Write-Host "   3. Navigate to Dashboards > Grade Management"
Write-Host "   4. Generate traffic to see metrics"
Write-Host ""
Write-Host "🛑 To stop: docker-compose down" -ForegroundColor Red
Write-Host ""
