# sockssh

`sockssh` is a command-line tool that SSH connections through a SOCKS proxy by automatically configuring the SSH ProxyCommand.

## Features

- Configure multiple SOCKS proxy environments
- YAML-based configuration
- Support for additional SSH options
- Environment-specific proxy settings

### Prerequisites

- OpenSSH client
- netcat (nc) command
- Running SOCKS proxy server

## Installation

Build the tool:

```bash
$ go build
```

## Configuration

Create a configuration file at `~/.config/sockssh.yaml` with the following structure:

```yaml
defaults:
  socks_server: 127.0.0.1
  port: 1080

environments:
  dev:
    socks_server: dev-proxy.example.com
    port: 1080
  stg:
    socks_server: stg-proxy.example.com
    port: 1080
  prod:
    socks_server: prod-proxy.example.com
    port: 1080
```

## Usage

The basic syntax for using `sockssh` is:

```
sockssh [-env <environment>] [--] <target> [ssh-options...]
```

### Examples

Connect to a host using the default SOCKS proxy settings:

```bash
$ sockssh -- user@example.com
```

Connect using a specific environment's proxy settings:

```bash
$ sockssh -env prod -- user@example.com
```

Use additional SSH options:

```bash
$ sockssh -- user@example.com -i ~/.ssh/custom_key
$ sockssh -- user@example.com -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
