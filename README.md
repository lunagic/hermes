# hermes

Shell Command Delivery Tool

## Installation

```shell
go install github.com/aaronellington/hermes@latest
```

## Command Line Usage

```shell
hermes --help
```

| Env Var                    | Default | Description                                                                                        |
| -------------------------- | ------- | -------------------------------------------------------------------------------------------------- |
| `HERMES_WORKING_DIRECTORY` | `"."`   | Where to look for the `hermes.yaml` file only if one is not found in the current working directory |

## Full Example

```shell
# Run the "update" operation
hermes run update
# Run the "reboot" operation
hermes run reboot
```

```yaml
# hermes.yaml
hosts:
  - hostname: macbook-pro-01
    tags: [mac, docker, primary]
  - hostname: server-01
    tags: [debian, docker, zfs]
  - hostname: server-02
    tags: [debian, docker, zfs]

tasks:
  update_apt:
    user: root
    host_tags:
      - debian
    commands:
      - systemctl mask systemd-networkd-wait-online.service
      - DEBIAN_FRONTEND=noninteractive apt-get update
      - DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade
      - DEBIAN_FRONTEND=noninteractive apt-get -y autoremove

  update_brew:
    host_tags:
      - mac
    commands:
      - brew update
      - brew upgrade
      - brew cleanup

  check_zpool:
    host_tags:
      - zfs
    commands:
      - zpool status -x 2>&1 | grep "all pools are healthy"

  cleanup_docker:
    host_tags:
      - docker
    commands:
      - docker image prune -f
      - docker container prune -f
      - docker volume prune -f
      - docker network prune -f

  check_reboot:
    user: root
    host_tags:
      - debian
      - fedora
    commands:
      - test ! -f /var/run/reboot-required

  reboot:
    user: root
    if: test -f /var/run/reboot-required
    host_tags:
      - debian
    commands:
      - reboot

operations:
  update:
    description: update my endpoints
    tasks:
      - update_apt
      - update_brew
      - cleanup_docker
      - check_zpool
      - check_reboot
  reboot:
    description: reboot endpoints that need it
    tasks:
      - reboot
```
