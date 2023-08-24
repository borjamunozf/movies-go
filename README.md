# Introduction

Toy project based in the book Microservices with Go.
- Setup multiple microservices.
- See layout.
- Check REST API & gRPC
- Think about tracing, observability & so in a distributed env


**Ongoing**

# Service discovery

- Client side first (Consul)

```bash
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0
```