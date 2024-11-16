# go-proxy

Simple forward proxy. Utilizing redis on local caching

## Getting Started

To run on your machine:

- run this command to install :

```bash
sudo curl -fsSL https://raw.githubusercontent.com/frsfahd/go-proxy/refs/heads/main/install.sh | bash
```

- create config file (`.yaml`) for defining redis connection (if you don't have redis instance, spin one up using `docker-compose.yml`). use this structure for config file :

```yaml
redis:
  host:
  port:
  username:
  password:
```

- run the app :

```bash
go-proxy --port [PORT] --origin [target-server] --config [config-file]
```

## MakeFile

Run build make command

```bash
make all
```

Build the application

```bash
make build
```

Build the application in windows

```bash
make build-windows
```

Run the application

```bash
make run
```

Create redis container

```bash
make docker-run
```

Shutdown redis Container

```bash
make docker-down
```

Live reload the application:

```bash
make watch
```

Clean up binary from the last build:

```bash
make clean
```
