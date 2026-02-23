## Postgres

The named volume `postgres_data` is important — without it, your data is destroyed every time the container stops. Docker manages this volume on the host and it persists across restarts.

