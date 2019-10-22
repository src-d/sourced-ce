# Install source{d} Community Edition Dependencies

## Install Docker

_Please note that Docker Toolbox is not supported. In case that you're running Docker Toolbox, please consider updating to newer Docker Desktop for Mac or Docker Desktop for Windows._

_For Linux installations, using Docker installed from a snap package is not supported._

Follow the instructions based on your OS:

- [Docker for Ubuntu Linux](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce-1)
- [Docker for Arch Linux](https://wiki.archlinux.org/index.php/Docker#Installation)
- [Docker for macOS](https://store.docker.com/editions/community/docker-ce-desktop-mac)
- [Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows) (Make sure to read the [system requirements for Docker on Windows](https://docs.docker.com/docker-for-windows/install/).)

Minimal supported version 18.02.0.

## Docker Compose

**source{d} CE** is deployed as Docker containers, using [Docker Compose](https://docs.docker.com/compose), but it is not required to have a local installation of Docker Compose; if it is not found it will be downloaded from docker sources, and deployed inside a container.

If you prefer a local installation of Docker Compose, or you have no access to internet to download it, you can follow the [Docker Compose install guide](https://docs.docker.com/compose/install)

Minimal supported version 1.20.0.

## Internet Connection

source{d} CE automatically fetches some resources from the Internet when they are not found locally, so in order to use all source{d} CE capacities, Internet connection is needed.

For more details, you can refer to [Why Do I Need Internet Connection?](../learn-more/faq.md#why-do-i-need-internet-connection)
