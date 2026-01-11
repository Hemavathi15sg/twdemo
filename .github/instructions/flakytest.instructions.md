---
applyTo: '**'
---


## Flaky Test Self-Healing Pipeline Instructions

### Context
This project implements self-healing mechanisms for flaky tests in CI/CD pipelines to automatically detect, analyze, and remediate unstable tests.

### Flaky Test Detection
- Implement retry logic with exponential backoff (max 3 retries)
- Track test execution history and failure patterns
- Calculate flakiness score based on pass/fail ratio over last 10 runs
- Mark tests as flaky if failure rate is between 10-90%
- Store flaky test metadata in `.flaky-tests.json`

### Self-Healing Strategies
1. **Timing Issues**: Add explicit waits, increase timeouts progressively
2. **Race Conditions**: Implement proper synchronization, use locks/semaphores
3. **Resource Contention**: Isolate test resources, use unique identifiers
4. **Network Flakiness**: Mock external dependencies, implement circuit breakers
5. **Environment Issues**: Reset state between tests, cleanup resources

### Pipeline Integration
- Run flaky test detection before main test suite
- Automatically quarantine tests exceeding flakiness threshold
- Generate reports with root cause analysis suggestions
- Create GitHub issues automatically for persistent flaky tests
- Send notifications to relevant teams via webhooks

### Code Generation Guidelines
- Always add retry decorators to potentially flaky tests
- Include detailed logging for failure analysis
- Implement proper teardown and cleanup in test fixtures
- Use deterministic test data and avoid time-based assertions
- Add comments explaining flakiness mitigation strategies

### Monitoring & Reporting
- Track metrics: flaky test count, auto-fix success rate, MTTR
- Generate weekly flaky test dashboards
- Maintain historical trends of test stability
- Alert on spike in flaky test detections