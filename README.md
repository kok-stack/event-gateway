# event-gateway
event gateway

## test command
```shell
docker run --rm -p 6379:6379 redis
curl -v "http://localhost:8080"  -X POST  -H "Ce-Id: 1"  -H "Ce-Specversion: 1.0"  -H "Ce-Type: greeting"  -H "Ce-Source: not-sendoff"  -H "Content-Type: application/json"  -d '{"msg":"Hello Knative!"}'

```