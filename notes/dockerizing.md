# From Claude
Here's a complete setup with your Go app, Postgres, and pgAdmin:

**docker-compose.yml**
```yaml
services:
  db:
    image: postgres:16
    restart: unless-stopped
    environment:
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: myapp_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myapp"]
      interval: 5s
      timeout: 5s
      retries: 5

  pgadmin:
    image: dpage/pgadmin4
    restart: unless-stopped
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@admin.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "5050:80"
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - pgadmin_data:/var/lib/pgadmin

  app:
    build: .
    restart: unless-stopped
    depends_on:
      db:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://myapp:secret@db:5432/myapp_dev?sslmode=disable
    ports:
      - "8080:8080"

volumes:
  postgres_data:
  pgadmin_data:
```

**Dockerfile** for your Go app:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

A minimal **main.go** to verify the DB connection on startup:
```go
package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("failed to open db:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("failed to connect to db:", err)
    }
    log.Println("connected to postgres")

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            http.Error(w, "db unreachable", http.StatusServiceUnavailable)
            return
        }
        w.Write([]byte("ok"))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Starting everything:**
```bash
docker compose up --build
```

**Connecting pgAdmin to your database:**

Once running, open `http://localhost:5050` and log in with `admin@admin.com` / `admin`. Then register a server with these details:

- **Host:** `db` (the service name, not localhost)
- **Port:** `5432`
- **Database:** `myapp_dev`
- **Username:** `myapp`
- **Password:** `secret`

The hostname has to be `db` rather than `localhost` because pgAdmin is inside the Docker network and resolves the Postgres container by its service name.

**Tip:** You can pre-configure the server connection so pgAdmin connects automatically on startup by mounting a servers.json file:

```yaml
pgadmin:
  volumes:
    - pgadmin_data:/var/lib/pgadmin
    - ./pgadmin/servers.json:/pgadmin4/servers.json
```

```json
{
  "Servers": {
    "1": {
      "Name": "Local Dev",
      "Group": "Servers",
      "Host": "db",
      "Port": 5432,
      "MaintenanceDB": "myapp_dev",
      "Username": "myapp",
      "PassFile": "/pgpassfile"
    }
  }
}
```

This saves you from manually registering the server every time the pgAdmin container is recreated.