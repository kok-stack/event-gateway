# event-gateway
event gateway

## test command
```shell
docker run --rm -it -p 2113:2113 -p 1113:1113 eventstore/eventstore:20.10.2-buster-slim --insecure --run-projections=All --enable-external-tcp --enable-atom-pub-over-http

curl -v "http://localhost:8080"  -X POST  -H "Ce-Id: 1"  -H "Ce-Specversion: 1.0"  -H "Ce-Type: greeting"  -H "Ce-Source: not-sendoff"  -H "Content-Type: application/json"  -d '{"msg":"Hello Knative!"}'

```