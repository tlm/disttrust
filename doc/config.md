# Configuration

disttrust is configured by supplying one or more json/yaml files in the command
line args and also the inclusion of a directory where all yaml/json files will
be processed.

# Hierarchy

All config files have the following hierachy that can be merged from multiple
sources.

```yaml
providers: []
anchors: []
```

# Providers
The providers array of the config file defines one or more providers in use as
an object. Each supported key in the object is defined below.

**example:**
```yaml
providers:
- id: vault
  name: example-provider
  options: {} # provider specific options. In this case they would be vault options
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

* approle:

    `roleId` - required - the approle role-id to use

    `secretId` - required - the approle secret-id to use

`path` - required
The path to where the pki backend is mounted with in vault. Path needs to be the
root of the backend.

`role` - required
Vault pki backend role to use when issuing certificates

`request` - optional
Defines a set of config options that controls how disttrust makes signing
requests to [vault pki](https://www.vaultproject.io/api/secret/pki/index.html)
backend.


* `csr` - optional default false - controls if disttrust makes a csr and ask's vault to sign it through the sign-verbatim endpoint. Using this method disttrust also generates private keys used for the certificate
* `keyType` - optional default rsa - One of rsa or ecdsa to pick the algo used for making public/private key pairs
* `rsaBits` - optional default 2048 - Numbers of bits to use for rsa keys when generating csr's
* `ecdsaCurve` - optional default p246 - ecdsa curve to use. Valid values are p224, p256, p384, p521

# Anchors

The anchors array of the config file defines one or more certificates to
issue on the local system from defined providers.

**example:**
```yaml
anchors:
- provider: example-provider
  cn: example-common-name
  altNames:
    dnsNames:
    - example.com
    ipAddresses:
    - fe80::feab:123
  dest: file
  destOpts:
    caFile: ./out/test-ca.pem
    certFile: ./out/test-cert.pem
    certBundleFile: ./out/test-bundle.pem
    privKeyFile: ./out/test-key.pem
    privKeyFileMode: 0600
  action:
    command:
    - /opt/local/bin/bash
    - -c
    - echo ff >> ./out/out.txt
```

`provider` - required

The name of previously defined provider to use

`cn` - required

Common name attribute for the certificate to issue

`atlNames` - optional

Object of alternative names to use in the certificate. Supported object keys
are.

* `DNSNames` - optional - list of string dns names to add to the certificate
* `IPAddresses` - optional - list of ip address alt names to add to the
  certificate

`organization` - optional

List of string organizations to put in the certificate.

`organizationalUnit` - optional

List of organizational units to put in the certificate

`action` - required

Object describing the actions to take when certificates are issues/updated. Supported values are:

* `command` - required - List of command with arguments to execute

`dest` - required

The destination location for issued certificates. Currently supported values are:
* `file`
* `aggregate`
* `template`

`destOpts` - required

dest specific options based on the choosen dest type. See below for dest
specific options.

### File Options

`caFile` - optional

file path of output

`caFileMode` - optional

output file mode

`caFileGid` - optional

group id to set on output file

`caFileUid` - optional

user id to set on output file

`certFile` - optional

file path of output

`certFileMode` - optional

output file mode

`certFileGid` - optional

group id to set on output file

`certFileUid` - optional

user id to set on output file

`certBundleFile` - optional

location to output CA to

`certBundleFileMode` - optional

output file mode

`certBundleFileGid` - optional

group id to set on output file

`certBundleFileUid` - optional

user id to set on output file

`privKeyFile` - optional

location to output CA to

`privKeyMode` - optional

output file mode

`privKeyGid` - optional

group id to set on output file

`privKeyUid` - optional

user id to set on output file

### Template Options
Allows making custom output templates with certificate data

`source` - required

Source input file template. See https://golang.org/pkg/text/template/ for
information on how to use go templates. Additional values can be used.
* `.CA` - Certifcate authority
* `.Certificate` - Generated certificate in pem
* `.CABundle` - Certificate authority bundle
* `.PrivateKey` - Generated private key
* `.Serial` - Serial number of the certificate

`out` - required

location to output compiled template to

`gid` - optional

group id to set on output file

`uid` - optional

user id to set on output file

`mode` - optional

output file mode

### Aggregate options

Allows grouping zero or more dests together

`dests` - optional

Array of dest objects that match the root. Takes the form of:
```yaml
dests:
- dest: "..."
  destOpts: "..."
```
