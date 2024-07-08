# blocky-ui

Simple web interface for [blocky](https://github.com/0xERR0R/blocky), a lightweight DNS proxy and DNS ad-blocker.

## Features

- web interface for all blocky API functions:
    - toggle or pause blocking
    - execute DNS queries
    - refresh allow and deny lists
    - flush DNS cache
- uses [HTMX](https://htmx.org/) to get status and query updates
- **disabled Javascript:** no lost functionality, uses `<form>` for full page reloads

### Screenshot

![screenshot: status and query results](screenshot.png)

## Config

All config is done via environment variables:

| Name             | Info                                                                   | Default                  |
|------------------|------------------------------------------------------------------------|--------------------------|
| `HOST`           | HTTP host                                                              | `0.0.0.0`                |
| `PORT`           | HTTP port                                                              | `3000`                   |
| `API_BASE_URL`   | blocky API URL, no trailing slash                                      | `http://blocky:4000/api` |
| `PAUSE_DURATION` | [Duration](https://pkg.go.dev/time#ParseDuration) for the pause button | `5m`                     |

## Run with docker

*Not pushed to any registry yet, need to build it yourself!*

Use `-e ...` or `--env-file` to set the variables if needed.

```shell
docker buildx build --build-arg VERSION=latest . -t blocky-ui
docker run --name blocky-ui -p 3000:3000 blocky-ui
```

### compose.yml

*Not pushed to any registry yet, see above for build command.*

```yaml
services:
  blocky:
    image: spx01/blocky
    container_name: blocky
    # ... other blocky config
    # https://0xerr0r.github.io/blocky/latest/installation/#run-with-docker-compose

  blocky-ui:
    image: blocky-ui
    container_name: blocky-ui
    restart: unless-stopped
    depends_on:
      - blocky
    environment:
      API_BASE_URL: http://blocky:4000/api
    ports:
      - "3000:3000/tcp"
```

## Run locally for tinkering

This project is set up for [live-reloading with air](https://github.com/air-verse/air), so just run `air` and
open `http://127.0.0.1:8080/`

## Notes

This is my first Go project, so issues and other comments are welcome!
