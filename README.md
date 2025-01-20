## Dashboard Request Per second Per Path
Dahboard Type: Timeseries
```promql
histogram_quantile(0.95, sum by(le, path) (rate(api_request_duration_seconds_bucket[$__rate_interval])))
```