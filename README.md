## How to run the project!
  - this project for testing idenpotency messaging to consumer queue.
### you can run the docker for dependency
  * redis
  * rabbitmq (http://localhost:15672)

```bash
docker run -d -p 6379:6379 -e ALLOW_EMPTY_PASSWORD=yes bitnami/redis:latest
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:3.8.4-management
```

