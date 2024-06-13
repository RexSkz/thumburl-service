# ThumbURL Service

A self-hosted API service written in Golang to generate thumbnails of given URLs using headless Chromium.

## Prerequisites

You should have an HTTP service that supports [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/). You can provide one for testing purposes using the following command (replace `/path/to/chromium` with the path to your Chrome/chromium executable):

```bash
# Start chromium with remote debugging port,
# you can see everything through the familiar browser GUI
/path/to/chromium --remote-debugging-port=9222 --disable-gpu

# Start "headless" chromium, without GUI
/path/to/chromium --remote-debugging-port=9222 --disable-gpu --headless
```

For production use, you may want to use a Docker version of chromium, e.g. `nextools/chromium`.

```bash
docker run \
  --rm \
  --name chromium \
  -p 9222:9222 \
  nextools/chromium:latest
```

Through docker-compose:

```yaml
version: "3"
services:
  chromium:
    image: nextools/chromium:latest
    container_name: chromium
    hostname: chromium
    volumes:
      # Mount the fonts directory to the container.
      # https://github.com/nextools/images/tree/master/chromium#add-custom-fonts
      - ./fonts:/home/chromium/.fonts
    ports:
      - 9222:9222
    restart: unless-stopped
```
