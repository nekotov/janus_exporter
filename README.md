# Janus exporter

This is a simple Prometheus exporter for Janus WebRTC Gateway. It uses the Janus REST API to get the metrics.

## Usage

```bash
janus_exporter -h
```

## Metrics

- janus_sessions
- janus_handlers
- janus_subscribers
- janus_packets_in
- janus_packets_out
- janus_bytes_in
- janus_bytes_out
- janus_clients_ips

Other metrics will be added in the future.