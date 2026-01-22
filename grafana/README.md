# Grafana Monitoring Setup - Grade Management API

## Overview

Complete observability stack for monitoring the Grade Calculation Engine (TEC-31) with:
- **Grafana** - Visualization and dashboards
- **Prometheus** - Metrics collection and storage
- **Loki** - Log aggregation
- **Redis** - Cache metrics

## Quick Start

### 1. Start Monitoring Stack

```bash
docker-compose up -d
```

This starts:
- Redis (port 6379)
- Prometheus (port 9090)
- Grafana (port 3000)
- Loki (port 3100)

### 2. Access Grafana

**URL:** http://localhost:3000

**Default Credentials:**
- Username: `admin`
- Password: `admin`

### 3. View Dashboards

Pre-configured dashboard available:
- **Grade Calculation Engine - Performance Dashboard**
  - Response time monitoring (P50, P95, P99)
  - API request rates
  - Grade distribution by letter grade
  - Error rate tracking
  - Cache hit rate
  - Performance SLA validation (200ms threshold)

## Dashboard Features

### Performance Metrics (TEC-31)

1. **Response Time Gauge**
   - P95 response time
   - Thresholds: Green (<150ms), Yellow (150-200ms), Red (>200ms)
   - Validates TEC-31 requirement: <200ms

2. **API Request Rate**
   - Real-time request rate for all grade endpoints
   - `/api/grades/calculate` highlighted
   - 5-minute rolling average

3. **Grade Distribution**
   - Donut chart with Figma design tokens
   - Distribution by letter grade (A+, A, B, C, D, F)
   - Color-coded by status (success, info, warning, critical)

4. **Error Rate**
   - Percentage of 5xx errors
   - Threshold: <1% (green)

5. **Cache Hit Rate**
   - Redis cache effectiveness
   - Threshold: >85% (green)

6. **Performance SLA Timeline**
   - P50, P95, P99 response times over time
   - 200ms threshold line (TEC-31 requirement)
   - Mean, max, min calculations

7. **Grade Status Distribution**
   - Figma alert states: success, info, warning, critical
   - Mapped to grade color tokens

8. **Grade Curve Usage**
   - Tracks how often grade curves are applied
   - Bar gauge visualization

## Metrics Endpoints

### Application Metrics

The Go application should expose metrics at:
```
http://localhost:8080/metrics
```

**Required Metrics for Dashboard:**

```prometheus
# Response time histogram
http_request_duration_milliseconds_bucket{endpoint="/api/grades/calculate"}

# Request counter
http_requests_total{endpoint="/api/grades/calculate", status_code="200"}

# Grade distribution
grade_letter_distribution{letter_grade="A+"}
grade_letter_distribution{letter_grade="A"}
grade_letter_distribution{letter_grade="B"}
# ... etc

# Grade status distribution (Figma tokens)
grades_calculated_total{grade_status="success"}
grades_calculated_total{grade_status="info"}
grades_calculated_total{grade_status="warning"}
grades_calculated_total{grade_status="critical"}

# Curve usage
grades_with_curve_total{curve_applied="true"}
grades_with_curve_total{curve_applied="false"}

# Cache metrics
redis_cache_hits
redis_cache_misses
```

## Prometheus Configuration

### Scrape Targets

1. **Grade Management API** (every 5s)
   - Target: `host.docker.internal:8080`
   - Metrics path: `/metrics`

2. **Redis** (every 15s)
   - Target: `redis:6379`

3. **Prometheus** (self-monitoring)
   - Target: `localhost:9090`

### Configuration File

Location: `grafana/prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'grade-management-api'
    static_configs:
      - targets: ['host.docker.internal:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s
```

## Grafana Configuration

### Datasources

Auto-provisioned from `grafana/provisioning/datasources/datasources.yml`:
- Prometheus (default)
- Loki
- Redis

### Dashboards

Auto-loaded from `grafana/dashboards/`:
- `grade-calculation-performance.json`

## Directory Structure

```
grafana/
├── prometheus.yml                    # Prometheus config
├── provisioning/
│   ├── datasources/
│   │   └── datasources.yml          # Auto-provision datasources
│   └── dashboards/
│       └── dashboards.yml           # Dashboard provider config
└── dashboards/
    └── grade-calculation-performance.json  # Main dashboard
```

## Monitoring Alerts (Optional)

Add alerting rules to `prometheus.yml`:

```yaml
rule_files:
  - "alerts.yml"
```

Example alerts:
- Response time > 200ms for 5 minutes
- Error rate > 1% for 5 minutes
- Cache hit rate < 70% for 10 minutes

## Performance Validation

### TEC-31 Requirements

✅ **Response Time:** <200ms
- Monitored via P95 response time gauge
- Historical trend in performance timeline
- Alert if threshold breached

✅ **100 Student Capacity:**
- Request rate monitoring
- Performance under load tracking

✅ **Design Token Compliance:**
- Grade distribution by Figma colors
- Status mapping (success, info, warning, critical)

## Troubleshooting

### Grafana not starting
```bash
docker-compose logs grafana
```

### Metrics not appearing
1. Check Prometheus targets: http://localhost:9090/targets
2. Verify API is exposing `/metrics` endpoint
3. Check Prometheus logs:
   ```bash
   docker-compose logs prometheus
   ```

### Dashboard not loading
1. Verify dashboards.yml configuration
2. Check Grafana logs for provisioning errors
3. Manually import dashboard JSON via Grafana UI

## Development

### Reload Prometheus Config
```bash
curl -X POST http://localhost:9090/-/reload
```

### Update Dashboard
1. Edit `grafana/dashboards/grade-calculation-performance.json`
2. Grafana auto-reloads every 10 seconds

### Add New Metrics
1. Instrument Go code with Prometheus client
2. Update dashboard JSON with new panels
3. Restart containers if needed

## Production Considerations

1. **Security:**
   - Change default Grafana credentials
   - Enable authentication on Prometheus
   - Use environment variables for secrets

2. **Persistence:**
   - All data stored in Docker volumes
   - Backup volumes regularly

3. **Scaling:**
   - Consider Prometheus federation for multiple instances
   - Use Grafana Cloud for hosted solution

4. **Retention:**
   - Default Prometheus retention: 15 days
   - Adjust with `--storage.tsdb.retention.time` flag

## Resources

- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Redis Exporter](https://github.com/oliver006/redis_exporter)

## Support

For issues related to:
- **Dashboard Design:** See Figma design tokens documentation
- **Performance Issues:** Check TEC-31 requirements
- **Metrics Collection:** Verify Prometheus scrape configuration
