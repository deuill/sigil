# Sigil - A simple web application server

## Introduction

Sigil is a web application server, consisting of several "services", which provide Sigil's built-in functionality, and "engines", which provide scripting capabilities to Sigil by embedding external runtimes. It is written in Go and is designed to be simple to build, deploy and use.

## Building

Assuming you have all already installed build dependancies required via `go get`, building Sigil is simply a matter of running `make` in the project root. You may install Sigil by running `make install` or build a redistributable package using `make package`.

## Running

You may either run Sigil directly using the 'sigil' binary, or use the supplied init script, which will also handle permissions and locking. Sigil contains several services which may require elevated permissions to execute, so running through the init script is highly recommended for production environments.

## Configuration

Sigil does not require any initial configuration and relies on a good set of defaults for operation. However, it can either use environment variables, or a local configuration file (which is created automatically if it does not exist) located in `~/.config/sigil/config.ini` for overriding default values, using the following semantics:

Configuration values are namespaced under their service name and option key. Environment variables use an 'SIGIL_' prefix, and are uppercase, while `config.ini` variables are placed in sections using the service name as a key, and are lowercase. So, for an option 'port' under service 'http', the following methods could be used to set the corresponding variable to `8080`:

```shell
export SIGIL_HTTP_PORT=8080
```

set in the environment in which `sigil` is launched, or:

```ini
[http]
port = 8080
```

set as a persistent value in the default `~/.config/sigil/config.ini` file. Environment variables override file variables, which in turn override defaults.