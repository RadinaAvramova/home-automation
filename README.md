# Home Automation

Distributed home automation system. Largely a learning opportunity rather than a production-ready system.

The home automation system is made up of separate microservices which run on various devices that are distributed physically around the home.
Most of the core services run within a Kubernetes cluster, however there are some periphery services that sit outside of the cluster.
The hardware used is a combination of Raspberry Pis and a Synology NAS.

## Getting started

There are various tools that can be installed to aid development.

```shell
./tools/install
``` 

## Project structure
- `docs/`
  - Documentation about the system
- `libraries/`
  - Library code shared between all services
- `private/`
  - A git submodule containing mostly private configuration
- `services/`
  - The services that make up the system
- `tools/`
  - Useful tools for working with the system
- `web.x`
  - A web-based application
