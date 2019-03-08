# consul-release

This is a [BOSH](http://bosh.io) release for [consul](https://github.com/hashicorp/consul).

* [CI](https://wings.concourse.ci/teams/cf-infrastructure/pipelines/consul)
* [Roadmap](https://www.pivotaltracker.com/n/projects/1488988)

### Contents

* [Using Consul](#using-consul)
* [Deploying](#deploying)
* [Configuring](#configuring)
* [Known Issues](#known-issues)
* [Failure Recovery](#failure-recovery)
* [Contributing](#contributing)
* [Confab Tests](#confab-tests)
* [Acceptance Tests](#acceptance-tests)

## Using Consul

Consul is a distributed key-value store that provides a host of applications.
It can be used to provide service discovery, key-value configuration,
and distributed locks within cloud infrastructure environments.

### Within CloudFoundry

Principally, Consul is used to provide service discovery for many of
the components. Components can register services with Consul, making these
services available to other CloudFoundry components. A component looking to
discover other services would run a Consul agent locally, and lookup services
using DNS names. Consul transparently updates the DNS records across the cluster
as services start and stop, or pass/fail their health checks.

Additionally, Consul is able to store key-value data across its distributed
cluster. Cloud Foundry makes use of this feature by storing some simple
configuration data, making it reachable across all nodes in the cluster.

CloudFoundry also makes some use of Consul's distributed locks.
This feature is used to ensure that one, and only one, component is able to
perform some critical action at a time.

### Fault Tolerance and Data Durability

Consul is a distributed data-store and as such, conforms to some form of fault
tolerance under disadvantageous conditions. Broadly, these tolerances are
described by the [CAP Theorem](https://en.wikipedia.org/wiki/CAP_theorem),
specifying that a distributed computer system cannot provide all three of the
guarantees outlined in the theorem (consistency, availability,
and partition tolerance). In the default configuration, Consul has a preference
to guarantee consistency and partition tolerance over availability. This means
that under network partitioning, the cluster can become unavailable. The
unavailability of the cluster can result in the inability to write to the
key-value store, maintain or acquire distributed locks, or discover other
services. Consul makes this tradeoff with a preference for consistency of the
stored data in the case of network partitions. The Consul team has published
some [results](https://www.consul.io/docs/internals/jepsen.html) from their
testing of Consul's fault tolerance.

This behavior means that Consul may not be the best choice for persisting
critically important data. Not having explicitly supported backup-and-restore
workflows also makes guaranteeing data durability difficult.

## Deploying

In order to deploy consul-release you must follow the standard steps for deploying software with BOSH.

We assume you have already deployed and targeted a BOSH director. For more instructions on how to do that please see the [BOSH documentation](http://bosh.io/docs).

### 1. Uploading a stemcell

Find the "BOSH Lite Warden" stemcell you wish to use. [bosh.io](https://bosh.io/stemcells) provides a resource to find and download stemcells.  Then run `bosh upload stemcell STEMCELL_URL_OR_PATH_TO_DOWNLOADED_STEMCELL`.

### 2. Creating a release

From within the consul-release directory run `bosh create-release --force` to create a development release.

### 3. Uploading a release

Once you've created a development release run `bosh upload-release` to upload your development release to the director.

### 4. Using a sample deployment manifest

We provide a set of sample deployment manifests that can be used as a starting point for creating your own manifest, but they should not be considered comprehensive. They are located in the `manifests` directory.

### 5. Deploy.

Run `bosh -d DEPLOYMENT_NAME deploy DEPLOYMENT_MANIFEST_PATH`.

## Configuring

### Generating Keys and Certificates

We only support running Consul in secure mode, you will need to provide
certificates and keys for Consul.

1. Generate SSL Certificates and Keys:
To generate the certificates and keys that you need for Consul, we recommend
using certstrap. This repository contains a helper script, `scripts/generate-certs`.
This script uses certstrap to initialize a certificate authority (CA), and
generate the certificates and keys for Consul.

2.  All servers must have a certificate valid for `server.<datacenter>.<domain>` or 
the client will reject the handshake.
For a default consul configuration, this means that a server certificate with the common name `server.dc1.cf.internal` will need to be created.
Further documentation concerning TLS encryption may be found on the official consul [documentation](https://www.consul.io/docs/agent/encryption.html).

If you already have a CA, you may have an existing workflow. You can modify
the `generate-certs` script to use your existing CA instead of generating a new one.

The `generate-certs` script outputs files to the `./consul-certs` directory.


2. Create Gossip Encryption Keys:
To create an encryption key for use in the serf gossip protocol, provide an
arbitrary string value. The consul agent job template transforms this string
into a 16-byte Base64-encoded value for consumption by the consul process.

3. Update your manifest:
Copy the contents of each file in the `./consul-certs` directory, as well as the
value for your Gossip encryption key, into the proper sections of your manifest.

For reference see below:

```
properties:
  consul:
    encrypt_keys:
    - RANDOM-SECRET-VALUE
    ca_cert: |
      -----BEGIN CERTIFICATE-----
      ###########################################################
      #######           Your New CA Certificate           #######
      ###########################################################
      -----END CERTIFICATE-----
    agent_cert: |
      ----BEGIN CERTIFICATE----
      ###########################################################
      #######           Your New Agent Certificate        #######
      ###########################################################
      ----END CERTIFICATE----
    agent_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ###########################################################
      #######           Your New Agent Key                #######
      ###########################################################
      -----END RSA PRIVATE KEY-----
    server_cert: |
      ----BEGIN CERTIFICATE----
      ###########################################################
      #######           Your New Server Certificate       #######
      ###########################################################
      ----END CERTIFICATE----
      ----BEGIN CERTIFICATE----
      ###########################################################
      #######           Your New CA Certificate           #######
      ###########################################################
      ----END CERTIFICATE----
    server_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ###########################################################
      #######           Your New Server Key               #######
      ###########################################################
      -----END RSA PRIVATE KEY-----
```

### Rotating Certificates and Keys

Reference the [Security Configuration for Consul](https://docs.cloudfoundry.org/deploying/common/consul-security.html#rotating-certs)
in the CloudFoundry Docs for steps on rotating certificates and keys.

### Defining a Service

This Consul release allows consumers to declare services provided by jobs that
should be discoverable over DNS. Consul achieves this behavior by consuming a
service definition. A service definition can be given to consul by providing
some configuration information in the manifest properties for the given job.

Below is an example manifest snippet that provides a service for a hypothetical
database:

```
1 jobs:
2 - name: database
3   instances: 3
4   networks:
5   - name: default
6   resource_pool: default
7   templates:
8   - name: database
9     release: database
10  - name: consul_agent
11    release: consul
12  properties:
13    consul:
14      agent:
15        services:
16          big_database:
17            name: big_database
18            tags:
19            - db
20            - persistence
21            check:
22              script: /bin/check_db
23              interval: 10s
```

In this example we are defining a "database" service that we want to make
available through Consul's service discovery mechanism. The first step is
to include the `consul_agent` template on the "database" job. We can see this
addition on lines 10-11. Once this template has been added, we can include
the service definition (lines 15-23). We define a service called "big_database"
that defines its service health check. The health check is used to determine
the health of a service and will automatically register/deregister that
service with the discovery system depending upon the status of that check. The
structure of a service definition follows the same structure as they would be
defined in JSON, but translated into YAML to fit into a manifest. More
information about service registration can be found
[here](https://www.consul.io/docs/agent/services.html).

Note: When running on Windows VMs, we use powershell. To honor the exit code of the
health check script, we wrap the script as follows:
`powershell -Command some_script.ps1; Exit $LASTEXITCODE`.
Any user provided script should do the same.

### Health Checks

Health checks provide another level of functionality to the service discovery
mechanism of Consul. When a service is defined with a health check, it can be
registered/deregistered from the service discovery system. This means that were
a service like a hypothetical "database" to experience some loss in
availability, Consul would notice and update the service discovery entries to
route traffic around or away from that service. Defining health checks can be
done in several ways. The
[documentation](https://www.consul.io/docs/agent/checks.html) provides
several examples of health check definitions, including script, HTTP, TCP, TTL,
and Docker-based examples.

When a service is defined without an explicit health check, the consul_agent job
will provide a default check. That check is the equivalent of the following:

```
jobs:
- name: database
  instances: 3
  networks:
  - name: default
  resource_pool: default
  templates:
  - name: database
    release: database
  - name: consul_agent
    release: consul
  properties:
    consul:
      agent:
        services:
          database:
            check:
              name: dns_health_check
              script: /var/vcap/jobs/database/bin/dns_health_check
              interval: 3s
```

In the `check` section of that definition, we can see that it assumes a script
called `dns_health_check` is located in the `/var/vcap/jobs/SERVICE_NAME/bin`
directory. Not providing this script, and not explicitly defining some other
check in your service definition will result in a failing health check for the
service. If you wish to disable the health check for a service, you can specify
an empty `check` section. For example:

```
jobs:
- name: database
  ...
  properties:
    consul:
      agent:
        services:
          database:
            check: {}
```

## Known Issues

### 1-node clusters

It is not recommended to run a 1-node cluster in any "production" environment.
Having a 1-node cluster does not ensure any amount of data persistence.

WARNING: Scaling your cluster to or from a 1-node configuration may result in data loss.

## Failure Recovery

### TLS Certificate Issues

A common source of failure is TLS certification configuration.  If you have a failed
deploy and see errors related to certificates, authorities, "crypto", etc. it's a good
idea to confirm that:

* all Consul-related certificates and keys in your manifest are correctly PEM-encoded;
* certificates match their corresponding keys;
* certificates have been signed by the appropriate CA certificate; and
* the YAML syntax of your manifest is correct.

### Failed Deploys, Upgrades, Split-Brain Scenarios, etc.

In the event that the consul cluster ends up in a bad state that is difficult
to debug, a simple

```
bosh restart consul_server
```

should fix the cluster.

If the consul cluster does not recover via the above method, you have the option
of stopping the consul agent on each server node, removing its data store, and
then restarting the process:

```
monit stop consul_agent (on all server nodes in consul cluster)
rm -rf /var/vcap/store/consul_agent/* (on all server nodes in consul cluster)
monit start consul_agent (one-by-one on each server node in consul cluster)
```

There are often more graceful ways to solve specific issues, but it is hard
to document all of the possible failure modes and recovery steps. As long as
your Consul cluster does not contain critical data that cannot be repopulated,
this option is safe and will probably get you unstuck.  If you are debugging
a Consul server cluster in the context of a Cloud Foundry deployment, it is
indeed safe to follow the above steps.

Additional information about outage recovery can be found on the consul
[documentation page](https://www.consul.io/docs/guides/outage.html).


#### X/Y nodes reported success

NOTE: The above will not work if a consul server is failing because of an issue
with a client node. In fact if you were to perform the recovery steps outlined
above you may fall below quorum and take down the consul servers. To know if you
are in this situation check the `/var/vcap/sys/log/consul_agent/consul.stdout.log`.
If you see an error similar to the below message then you do not want to follow
the above steps:

```
2017/05/04 12:31:39 [INFO] serf: EventMemberFailed: some-node-0 10.0.0.100 {"timestamp":"1493901099.471777678","source":"confab","message":"confab.agent-client.set-keys.list-keys.request.failed","log_level":2,"data":{"error":"7/8 nodes reported success"}}
```

The important part is:

```
{"error":"214/215 nodes reported success"}
```
This indicates that the consul server attempted to set keys on a client node
and the client node failed to set the key.

An additional variation of this error message can be:
```
{"error":"1/215 nodes reported failure"}
```
This indicates that the consul server attempted to verify the keys on a client
node and the client node failed to reply.

As of consul-release v181 and later, this should only occur when you are
explicitly attempting to rotate keys by changing the value of the manifest
property `consul.encrypt_keys`.

The resolution in this case is to find the client node that is failing to set
keys and resolve whatever issue is occuring on that node. The most common issue
is either the disk is full/not writeable and consul cannot write to its 
datastore or consul cannot communicate with a node on port 8301. If you are unable
to determine the VM that is the root cause of the issue, you may still be able
to work around this issue by waiting longer for the key rotation to complete.
Try setting `confab.timeout_in_seconds` to a much larger value like 600
(10 minutes).

### Frequent Disappearance of Registered Services

Many BOSH jobs that colocate the `consul_agent` process do so in order to
register a service with Consul so that other jobs within the system can 
discover them.  If you observe frequent service discovery failures affecting
many services, this may be due to something affecting Consul's gossip
protocol.  Common causes include:

* network latency;
* network failures such as high packet loss;
* firewalls/ACLs preventing some `consul_agent`s communicating with others
  over TCP on port 8300 and both TCP and UDP on port 8301; and
* having a very large number of VMs driving CPU requirements for the
  `consul_agent`s too high for the current resources allocated to their VMs.

## Contributing

### Contributor License Agreement

Contributors must sign the Contributor License Agreement before their
contributions can be merged. Follow the directions
[here](https://www.cloudfoundry.org/community/contribute/) to complete
that process.

### Developer Workflow

Make sure that you are working against the `develop` branch. PRs submitted
against other branches will need to be resubmitted with the correct branch
targeted.

Before submitting a PR, make sure to run the test suites. Information about
how to run the suites can be seen in the [Confab Tests](#confab-tests) and
[Acceptance Tests](#acceptance-tests) sections.

## Confab Tests

Run the `confab` tests by executing the `src/confab/scripts/test` executable.

## Acceptance Tests

The acceptance tests deploy a new consul cluster and exercise a variety of features, including scaling the number of nodes, as well as destructive testing to verify resilience.

### Prerequisites

The following should be installed on the local machine:

- jq
- Golang (>= 1.7)

If using homebrew, these can be installed with:

```
brew install consul go jq
```

### Network setup

#### BOSH-Lite

Make sure you’ve run `bin/add-route` and `bin/enable_container_internet`.
This will setup some routing rules to give the tests access to the consul VMs.

#### AWS

You will want to run CONSATS from a VM within the same subnet specified in your manifest.
This assumes you are using a private subnet within a VPC.

### Environment setup

This repository assumes that it is the root of your `GOPATH`. Move it to `~/go/src/github.com/cloudfoundry-incubator`.

### Running the CONSATS

#### Running locally

Run all the tests with:

```
CONSATS_CONFIG=[/path/to/config_file.json] ./scripts/test
```

Run a specific set of tests with:

```
CONSATS_CONFIG=[/path/to/config_file.json] ./scripts/test <some test packages>
```

The `CONSATS_CONFIG` environment variable points to a configuration file which specifies the endpoint of the BOSH director.
When specifying location of the CONSATS_CONFIG, it must be an *absolute path on the filesystem*.

See below for more information on the contents of this configuration file.

### CONSATS config

An example config json for BOSH-lite would look like:

Note: If you're running the tests as BOSH errand (see below), then the integration
config json will be generated from the values you specify in the `consats.yml` manifest.

```json
cat > integration_config.json << EOF
{
  "bosh":{
    "target": "https://192.168.50.6:25555",
    "username": "admin",
    "password": "admin",
	"director_ca_cert": "some-cert"
  }
}
EOF
export CONSATS_CONFIG=$PWD/integration_config.json
```

The full set of config parameters is explained below:
* `bosh.target` (required) Public BOSH IP address that will be used to host test environment
* `bosh.username` (required) Username for the BOSH director login
* `bosh.password` (required) Password for the BOSH director login
* `bosh.director_ca_cert` (required for non-BOSH-lite environments) BOSH Director CA Cert
* `aws.subnet` Subnet ID for AWS deployments
* `aws.access_key_id` Key ID for AWS deployments
* `aws.secret_access_key` Secret Access Key for AWS deployments
* `aws.default_key_name` Default Key Name for AWS deployments
* `aws.default_security_groups` Security groups for AWS deployments
* `aws.region` Region for AWS deployments
* `registry.host` Host for the BOSH registry
* `registry.port` Port for the BOSH registry
* `registry.username` Username for the BOSH registry
* `registry.password` Password for the BOSH registry

#### Running as BOSH errand

##### Dependencies

The `acceptance-tests` BOSH errand assumes that the BOSH director has already uploaded the correct versions of the dependent releases.
The required releases are:
* [consul-release](http://bosh.io/releases/github.com/cloudfoundry-incubator/consul-release) or `bosh create-release && bosh upload release`
* [turbulence-release](http://bosh.io/releases/github.com/cppforlife/turbulence-release) (if running with turbulence)

##### Creating a consats deployment manifest

We provide an example deployment manifest for running the errand.
The manifest can be used by replacing all of the placeholder values in the file `manifests/consats.yml`.

##### Deploying the errand

Run `bosh deploy manifests/consats.yml -d consats`.

##### Running the errand

Run `bosh run-errand acceptance-tests -d consats`.

