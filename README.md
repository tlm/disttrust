# Overview

disttrust is a program for requesting tls certificates from a provider and
maintaining their validity on the local system. It currently only implements the
**vault** pki provider using approle based auth. The design is such that new
providers and auth methods can be added with minimal work.

# Badges
[![Go Report Card](https://goreportcard.com/badge/github.com/tlmiller/disttrust)](https://goreportcard.com/report/github.com/tlmiller/disttrust)

# Building & Installation

## From Source

### Requirements
- go >= 1.8

Installation on the local system can be done with a normal go get

`go get -u github.com/tlmiller/disttrust`

# Usage

disttrust is configured from one or more files supplied by a flag. For ease
of deployments many seperate config files can be provided each with different
providers and anchors that will be aggregated together at run time.

`disttrust -c <cfile.json1> -c <cfile.json2> ...`

See [config](#Config) for documentation on the configuration properties.

# Config
disttrust is configured using json files. The core of the file is a json object
with keys to the various parts.

```json
{
    "providers": [...],
    "anchors": [...]
}
```

## Providers
The providers array of the config file defines one or more providers in use as a
json object. Each supported key in the object is defined below.

**example:**
```json
{
    "id": "vault",
    "name": "etcd-server",
    "options": {
    	"address": "http://127.0.0.1:8200",
    	"authMethod": "approle",
    	"authOpts": {
    		"roleId": "a8e6327d-5585-ba52-330b-7fcfd8baa33c",
    		"secretId": "517a8361-38c8-7cb8-24ee-a94544163ad0"
    	},
    	"path": "klust/tst/etcd/pki",
    	"role": "server"
}
```

`id` - required

Identifies the type of the provider. Current supported values are:
* `vault`

`name` - required

User defined name to describe this provider. Name must be unique.

`options` - required

Provider specific options object for configuring the type of provider. See
individual providers below for more information.

### Vault Options

`address` - defaults to http://127.0.0.1:8200

Full URI to the vault cluster to use

`authMethod` - required

Defines the authMethod to use for vault. Supported values are:
* approle

`authOpts` - required
Defines an object map of properties for the auth method. Below value are
supported based on auth method.

* approle

    `roleId` - required - the approle role-id to use
    
    `secretId` - required - the approle secret-id to use

`path` - required
The path to where the pki backend is mounted with in vault. Path needs to be the
root of the backend.

`role` - required
Vault pki backend role to use when issuing certificates

## Anchors

The anchors array of the config file defines one or more certificates to
issue on the local system from defined providers.

**example:**
```json
{
	"provider": "etcd-server",
	"cn": "etcd-1.example.com",
	"altNames": [
		"etcd-all.example.com"
	],
	"dest": "file",
	"destOpts": {
        "caFile": "./out/test-ca.pem",
        "certFile": "./out/test-cert.pem",
        "certBundleFile": "./out/test-bundle.pem",
        "privKeyFile": "./out/test-key.pem"
	},
    "action": {
        "command": [
			"/opt/local/bin/bash",
			"-c",
			"echo ff >> ./out/out.txt"
		]
    }
}
```

`provider` - required

The name of previously defined provider to use

`cn` - required

Common name attribute for the certificate to issue

`atlNames` - optional

List of alternative names to use in the issued certificate

`dest` - required

The destination location for issues certificates. Currently supported values are:
* `file`

`destOpts` - required

dest specific options. Options specific to the dest type choosen.

* `file`

    `caFile` - optional - location to output ca to
    `certFile` - optional - location to output issued certificate to
    `certBundleFile` - optional location to output issued certificate bundle to
    `privKeyFile` - optional location to output private key to
