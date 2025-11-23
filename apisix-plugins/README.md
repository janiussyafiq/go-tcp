# APISIX Custom Plugins

Three custom Lua plugins for Apache APISIX demonstrating API gateway middleware patterns.

> **⚠️ Security Notice**: This project uses the default APISIX admin key for local development only. Never use default credentials in production. Change the `admin_key` in `apisix_conf/config.yaml` before deploying.

## Plugins

1. **request-logger** - Logs HTTP method and URI for every request
2. **ip-blocker** - Blocks requests from specified IP addresses
3. **header-injector** - Adds custom headers to requests before proxying

## Prerequisites

- Docker & Docker Compose
- curl (for testing)

## Quick Start

### 1. Start APISIX
```bash
docker-compose up -d
```

Wait 10 seconds for services to initialize.

### 2. Copy Plugins to APISIX
```bash
docker exec -it apisix-plugins-apisix-1 cp /opt/apisix/plugins/request-logger.lua /usr/local/apisix/apisix/plugins/
docker exec -it apisix-plugins-apisix-1 cp /opt/apisix/plugins/ip-blocker.lua /usr/local/apisix/apisix/plugins/
docker exec -it apisix-plugins-apisix-1 cp /opt/apisix/plugins/header-injector.lua /usr/local/apisix/apisix/plugins/
docker-compose restart apisix
```

### 3. Set Admin Key
```bash
export ADMIN_KEY=edd1c9f034335f136f87ad84b625c8f1  # Default dev key
```

### 4. Create Test Routes

**Request Logger:**
```bash
curl -X PUT http://127.0.0.1:9180/apisix/admin/routes/1 \
-H "X-API-KEY: $ADMIN_KEY" \
-H 'Content-Type: application/json' \
-d '{
  "uri": "/test",
  "plugins": {"request-logger": {}},
  "upstream": {"type": "roundrobin", "nodes": {"httpbin.org:80": 1}}
}'
```

**IP Blocker:**
```bash
curl -X PUT http://127.0.0.1:9180/apisix/admin/routes/2 \
-H "X-API-KEY: $ADMIN_KEY" \
-H 'Content-Type: application/json' \
-d '{
  "uri": "/blocked",
  "plugins": {"ip-blocker": {"blocked_ips": ["192.168.65.1"]}},
  "upstream": {"type": "roundrobin", "nodes": {"httpbin.org:80": 1}}
}'
```

**Header Injector:**
```bash
curl -X PUT http://127.0.0.1:9180/apisix/admin/routes/4 \
-H "X-API-KEY: $ADMIN_KEY" \
-H 'Content-Type: application/json' \
-d '{
  "uri": "/headers",
  "plugins": {
    "header-injector": {
      "headers": {
        "X-Custom-Token": "secret-123",
        "X-API-Version": "v1"
      }
    }
  },
  "upstream": {"type": "roundrobin", "nodes": {"httpbin.org:80": 1}}
}'
```

## Testing

**Request Logger:**
```bash
curl http://127.0.0.1:9080/test
docker exec -it apisix-plugins-apisix-1 tail -10 /usr/local/apisix/logs/error.log | grep "REQUEST:"
```

**IP Blocker:**
```bash
curl http://127.0.0.1:9080/blocked
# Expected: {"message":"Access denied"}
```

**Header Injector:**
```bash
curl http://127.0.0.1:9080/headers
# Check response for injected headers
```

## Project Structure
```
apisix-plugins/
├── docker-compose.yml          # APISIX and etcd services
├── apisix_conf/
│   └── config.yaml            # APISIX configuration
├── plugins/
│   ├── request-logger.lua     # Plugin 1: Request logger
│   ├── ip-blocker.lua         # Plugin 2: IP blocker
│   └── header-injector.lua    # Plugin 3: Header injector
└── apisix_logs/               # APISIX logs (gitignored)
```

## Security Notes

- The default admin key (`edd1c9f034335f136f87ad84b625c8f1`) is for **local development only**
- For production, generate a secure key and update `apisix_conf/config.yaml`
- Restrict admin API access with `allow_admin` configuration
- Use environment variables for sensitive configuration

## Tech Stack

- **Apache APISIX 3.7.0** - API Gateway
- **OpenResty** - Nginx + Lua
- **etcd 3.5** - Configuration storage
- **Lua** - Plugin language

## Cleanup
```bash
docker-compose down -v
```