# disttrust

A program for requesting tls certificates from a provider and
maintaining their validity on the local system. It currently only implements the
**vault** pki provider using approle based auth. The design is such that new
providers and auth methods can be added with minimal work.

disttrust was originally created as a tool to help maintain short lived
certificates in an automated fashion for Kubernetes clusters but can be used for
any scenario that requires maintaining valid certificates on a system.

The idea for this application was inspired by
[Digital Oceans](https://blog.digitalocean.com/vault-and-kubernetes/) work in
automating the same problem with consul-template. Disttrust aims to remove some
of the rough edges and provide better handling and monitoring overall.

# Version
disttrust is still in pre 1.0 releases, it is currently in use for production
env's but will remain pre 1.0 till the initial feature set has been finalised.
Plans include
  - More unit testing
  - Prometheus support
  - Fine grain health checks
  - Better error handling support on SIGHUP where the program will continue to
    run with its old configuration if the new confiration contains an error.
  - Possible support for KV db config support such as consul
  - CICD pipeline & linux package building

# Badges
[![Go Report Card](https://goreportcard.com/badge/github.com/tlmiller/disttrust)](https://goreportcard.com/report/github.com/tlmiller/disttrust)

# Building & Installation

## From Source

### Requirements
- go >= 1.11

Installation on the local system can be done with a normal go get

`go get -u github.com/tlmiller/disttrust`

# Usage

disttrust is configured from one or more files supplied by a flag. For ease
of deployments many seperate config files can be provided each with different
providers and anchors that will be merged together at run time.

`disttrust -c <cfile.yml> -c <cfile.json2> ... -c <config_dir>`

# Config
See [config](#Config) for documentation on the configuration properties.
