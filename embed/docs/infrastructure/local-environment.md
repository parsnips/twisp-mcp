---
title: Local Environment
description: The local environment helps streamline development and CICD integration testing with Twisp running on your own machine.
showNext: true
showPrev: true
---

## Starting the environment
The following shell commands can be used to get up and running:

```sh
docker pull public.ecr.aws/twisp/local:latest
docker run -p 3000:3000 -p 8080:8080 -p 8081:8081 public.ecr.aws/twisp/local:latest
```

You should see the startup logs in your terminal. The console should now be available in the browser at <http://localhost:3000>.

## Work through the Twisp tutorial
After starting up your local environment check out the Twisp [tutorial](/tutorials/twisp-101) for more information on using Twisp.

## Environment details
The local Twisp container image includes a number of services and data persistence options. Container images are available for both `amd64` and `arm64` architectures.

### Available services
- HTTP console port `3000`
- HTTP API port `8080`
- GRPC API port `8081`
- Health check command `/healthcheck`

> Note: the Twisp Local Instance does not require HTTP `Authorization` headers for access.

### Data persistence
- Set `DB_PATH` environment variable to `-e DB_PATH=/data`
- Mount a volume: `-v volume:/data`

### Account management
- Execute multi-tenant API calls by posting with a specific `X-Twisp-Account-Id` HTTP header
- If an `X-Twisp-Account-Id` header is not provided, the system will default to `X-Twisp-Account-Id: 000000000000`
- The console invokes the API on behalf of the default account
