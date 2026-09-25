# Toxiproxy

Sits between the clients (`../clients/`) and the RabbitMQ cluster (`../rabbitmq-server/`)
so network issues (latency, timeouts, bandwidth caps, connection resets, full outages) can be
injected on demand for resilience testing.

Deploy:

```
kubectl apply -f toxiproxy-configmap.yaml
kubectl apply -f toxiproxy-deployment.yaml
```

`toxiproxy-configmap.yaml` preloads two proxies at startup (see `toxiproxy-server -config`):

- `rabbitmq-amqp`: `0.0.0.0:5672` -> `tls.default.svc.cluster.local:5672`
- `rabbitmq-stream`: `0.0.0.0:5552` -> `tls.default.svc.cluster.local:5552`

The client manifests in `../clients/` already point at `toxiproxy.default.svc.cluster.local`
instead of the `tls` service directly, so traffic always flows through these proxies.

## Injecting issues

Port-forward the Toxiproxy API and use [`toxiproxy-cli`](https://github.com/Shopify/toxiproxy)
to add/remove toxics.

### Installing `toxiproxy-cli` on Linux

Download the prebuilt binary matching the server version used in
`toxiproxy-deployment.yaml` (`2.9.0`) and the machine architecture:

```
# amd64
curl -Lo toxiproxy-cli https://github.com/Shopify/toxiproxy/releases/download/v2.9.0/toxiproxy-cli-linux-amd64

# arm64
curl -Lo toxiproxy-cli https://github.com/Shopify/toxiproxy/releases/download/v2.9.0/toxiproxy-cli-linux-arm64

chmod +x toxiproxy-cli
sudo mv toxiproxy-cli /usr/local/bin/toxiproxy-cli
```

(Alternatively, with a Go toolchain installed: `go install github.com/Shopify/toxiproxy/v2/cli@v2.9.0`,
which installs to `$(go env GOPATH)/bin/cli` — rename/symlink it to `toxiproxy-cli`.)

Then port-forward the API and confirm it's reachable:

```
kubectl port-forward svc/toxiproxy 8474:8474
export TOXIPROXY_URL=http://localhost:8474

toxiproxy-cli list
```

Examples (target either `rabbitmq-amqp` or `rabbitmq-stream`):

```
# Add 2s latency with 500ms jitter
toxiproxy-cli toxic add -t latency -a latency=2000 -a jitter=500 rabbitmq-amqp

# Cap bandwidth to 64 KB/s
toxiproxy-cli toxic add -t bandwidth -a rate=64 rabbitmq-stream

# Drop the connection after 1s (simulates a dead peer)
toxiproxy-cli toxic add -t timeout -a timeout=1000 rabbitmq-amqp

# Reset the TCP connection immediately (simulates a crashed broker)
toxiproxy-cli toxic add -t reset_peer -a timeout=0 rabbitmq-stream

# Simulate a full outage: disable the proxy so connections are refused
toxiproxy-cli toggle rabbitmq-amqp

# Remove a toxic / restore normal behavior
toxiproxy-cli toxic remove -n latency_downstream rabbitmq-amqp
toxiproxy-cli toggle rabbitmq-amqp   # re-enable after a disable
```

Without `toxiproxy-cli`, the same operations work via the REST API, e.g.:

```
curl -X POST http://localhost:8474/proxies/rabbitmq-amqp/toxics \
  -d '{"name":"latency_downstream","type":"latency","stream":"downstream","attributes":{"latency":2000,"jitter":500}}'
```

## Simulating a periodic reset_peer (e.g. every 60s)

`reset_peer` isn't a repeating timer — once added it stays enabled (resetting every
connection that flows through it) until removed, and it only fires when traffic is
actually in flight. To get a one-off reset roughly every 60 seconds, add the toxic,
wait briefly for it to hit any live connection, then remove it again, on a loop.

**In-cluster (recommended):** `reset-peer-cronjob.yaml` runs this cycle every minute
against both proxies via a `CronJob` (schedule `* * * * *`), calling the Toxiproxy API
directly at `http://toxiproxy:8474`:

```
kubectl apply -f reset-peer-cronjob.yaml

# stop it
kubectl delete cronjob toxiproxy-reset-peer
# or pause without deleting
kubectl patch cronjob toxiproxy-reset-peer -p '{"spec":{"suspend":true}}'
```

**Ad hoc, from your machine:** with the API port-forwarded (see above), run a loop with
`toxiproxy-cli`:

```
while true; do
  toxiproxy-cli toxic add -t reset_peer -a timeout=0 -n periodic_reset rabbitmq-stream
  sleep 1
  toxiproxy-cli toxic remove -n periodic_reset rabbitmq-stream
  sleep 59
done
```
