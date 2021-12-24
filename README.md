# Elastic stack on Docker

[![Elastic Stack version](https://img.shields.io/badge/Elastic%20Stack-7.16.2-00bfb3?style=flat&logo=elastic-stack)](https://www.elastic.co/blog/category/releases)

Run the [Elastic stack](https://www.elastic.co/what-is/elk-stack) with Docker Compose.

It gives you the ability to analyze any data set by using the searching/aggregation capabilities of Elasticsearch and
the visualization power of Kibana.

![Animated demo](https://user-images.githubusercontent.com/3299086/140641708-cea70d17-cc04-459f-89d9-3fcb5c58bc35.gif)

Uses the official Docker images from Elastic:

* [Elasticsearch](https://github.com/elastic/elasticsearch/tree/master/distribution/docker)
* [Logstash](https://github.com/elastic/logstash/tree/main/docker)
* [Kibana](https://github.com/elastic/kibana/tree/master/src/dev/build/tasks/os_packages/docker_generator)
---


## Contents

1. [Features](#features)
1. [Requirements](#requirements)
   * [Host setup](#host-setup)
   * [Docker Desktop](#docker-desktop)
      * [Windows](#windows)
      * [macOS](#macos)
1. [Usage](#usage)
    * [Prepare docker host](#prepare-docker-host)
    * [Initial setup](#initial-setup)
    * [Docker network driver](#docker-network-driver)
    * [Cleanup](#cleanup)
    * [Access Kibana](#access-kibana)
    * [Default Kibana index pattern creation](#default-kibana-index-pattern-creation)
      * [Via the Kibana web UI](#via-the-kibana-web-ui)
      * [On the command line](#on-the-command-line)
1. [Configuration](#configuration)
    * [How to configure Elasticsearch](#how-to-configure-elasticsearch)
    * [How to configure Kibana](#how-to-configure-kibana)
      * [Kibana TLS](#kibana-tls)
    * [How to configure Logstash](#how-to-configure-logstash)
    * [How to scale out the Elasticsearch cluster](#how-to-scale-out-the-elasticsearch-cluster)
    * [Healthcheck binary](#healthcheck-binary)
2. [Extensibility](#extensibility)
    * [How to add plugins](#how-to-add-plugins)
3. [JVM tuning](#jvm-tuning)
    * [How to specify the amount of memory used by a service](#how-to-specify-the-amount-of-memory-used-by-a-service)
    * [How to enable a remote JMX connection to a service](#how-to-enable-a-remote-jmx-connection-to-a-service)
7. [Going further](#going-further)
    * [Swarm mode](#swarm-mode)


## Features

This repository based at [deviantony/docker-elk](https://github.com/deviantony/docker-elk/), but adapted to my own requirements. The main goal of this project is running production-ready single node Elasticsearch instance.

**Comparing to original repo:**

1. Using original container images. This time I don't use plugins and see no point to build custom images.
2. Using `basic` license by default.  
3. Enabled bootstrap checks.
4. Enabled TLS and X-Pack security features.
5. Configured container memory ulimits according to Elasticsearch documentation.
6. Added healthcheck scripts.
7. Added Logstash pipelines config file binding.


## Requirements
### Host setup

* [Docker Engine](https://docs.docker.com/install/) version **17.05** or newer
* [Docker Compose](https://docs.docker.com/compose/install/) version **1.20.0** or newer
* 3 GB of RAM

*:information_source: Especially on Linux, make sure your user has the [required permissions](https://docs.docker.com/install/linux/linux-postinstall/) to
interact with the Docker daemon.*
*:information_source: [Change Java heap](#how-to-specify-the-amount-of-memory-used-by-a-service) with your requirements.*
**:warning: Docker-compose commands below assume using docker-compose v2.**

### Docker Desktop
#### Windows

If you are using the legacy Hyper-V mode of _Docker Desktop for Windows_, ensure [File Sharing](https://docs.docker.com/desktop/windows/#file-sharing) is
enabled for the `C:` drive.

#### macOS

The default configuration of _Docker Desktop for Mac_ allows mounting files from `/Users/`, `/Volume/`, `/private/`,
`/tmp` and `/var/folders` exclusively. Make sure the repository is cloned in one of those locations, or follow the
instructions from the [documentation](https://docs.docker.com/desktop/mac/#file-sharing) to add more locations.

## Usage
### Prepare docker host

Increase virtual memory map:
```console
$ sudo sysctl -w vm.max_map_count=262144
```
**:warning: This is not persisted change. Make a file in `/etc/sysctl.d/` dir with this setting.**

Please, check [Elastic docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/system-config.html) for more information.

### Initial setup
**:warning: This project prepared to use Elasticsearch with enabled TLS for Elasticsearch. You can disable TLS in services configs and healthcheck scripts if you don't need TLS.**

1. Clone this repository onto the Docker host.
2. Follow [TLS setup settings](./tls/)
*:information_source: Instead built-in users activation, you can use Elasticsearch root user `elastic`, activated by Elasticsearch environment variable `ELASTIC_PASSWORD`.*
3. Enable built-in system accounts:
   1. Start Elasticsearch with `docker compose up elastic`
   2. After a few seconds run
         `docker compose exec elasticsearch bin/elasticsearch-setup-passwords auto --batch -u https://localhost:9200`
         That will generate passwords for system accounts.
   3. Fill passwords with generated ones in following files:
         `elastic/elastic_healthcheck.sh` `kibana/kibana.yml`, `kibana/kibana_healthcheck.sh`, `logstash/logstash.yml`, `logstash/pipeline/main.conf`
4. Start services locally using Docker Compose:`docker compose up --force-recreate`
    You can also run all services in the background (detached mode) by adding the `-d` flag to the above command.

### Docker network driver

There are two network drivers that can be used with docker-compose: `bridge` and `host`.

**bridge:** Add virtual network and pass-through selected ports. Also provide ability to use internal domain names (`elastic`, `kibana`, etc). Unfortunately, brings some routing overhead.

**host:** Just use host network.  No network isolation, no internal domains, no overhead.

According to [Rally](https://github.com/elastic/rally) testing with `metricbeat` race, there is no significant difference.

**Using host network:**
To use host network for Elastic stack, remove `network` and `ports` sections from `docker-compose.yml` file and add `network_mode: host` key to services you want to use host network driver. You can use all services with host network mode.
When Elasticsearch set to use host network, change `elasticsearch.hosts` to `localhost` both in Kibana and Logstash configs.

Check [docker compose reference](https://docs.docker.com/compose/compose-file/compose-file-v3/#network-configuration-reference) for more information.

### Cleanup

Elasticsearch data is persisted inside a volume by default.

In order to entirely shutdown the stack and remove all persisted data, use the following Docker Compose command:

```console
$ docker compose down -v
```

### Access Kibana

Give Kibana about a minute to initialize, then access the Kibana web UI by opening <http://localhost:5601> in a web
browser and use the following credentials to log in:

* user: *elastic*
* password: *\<your generated elastic password>*

### Default Kibana index pattern creation

When Kibana launches for the first time, it is not configured with any index pattern.

#### Via the Kibana web UI

*:information_source: You need to inject data into Logstash before being able to configure a Logstash index pattern via
the Kibana web UI.*

Navigate to the _Discover_ view of Kibana from the left sidebar. You will be prompted to create an index pattern. Enter
`logstash-*` to match Logstash indices then, on the next page, select `@timestamp` as the time filter field. Finally,
click _Create index pattern_ and return to the _Discover_ view to inspect your log entries.

Refer to [Connect Kibana with Elasticsearch](https://www.elastic.co/guide/en/kibana/current/connect-to-elasticsearch.html) and [Creating an index pattern](https://www.elastic.co/guide/en/kibana/current/index-patterns.html) for detailed
instructions about the index pattern configuration.

#### On the command line

Create an index pattern via the Kibana API:

```console
$ curl -XPOST -D- 'http://localhost:5601/api/saved_objects/index-pattern' \
    -H 'Content-Type: application/json' \
    -H 'kbn-version: 7.16.2' \
    -u elastic:<your generated elastic password> \
    -d '{"attributes":{"title":"logstash-*","timeFieldName":"@timestamp"}}'
```

The created pattern will automatically be marked as the default index pattern as soon as the Kibana UI is opened for the
first time.

## Configuration

*:information_source: Configuration is not dynamically reloaded, you will need to restart individual components after any configuration change.*

### How to configure Elasticsearch

Learn more about the security of the Elastic stack at [Secure the Elastic Stack](https://www.elastic.co/guide/en/elasticsearch/reference/current/secure-cluster.html).

The Elasticsearch configuration is stored in [`elastic/elasticsearch.yml`](./elastic/elasticsearch.yml).

You can also specify the options you want to override by setting environment variables inside the Compose file:

```yml
elastic:

  environment:
    network.host: _non_loopback_
    cluster.name: my-cluster
```

Please refer to the following documentation page for more details about how to configure Elasticsearch inside Docker
containers: [Install Elasticsearch with Docker](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html).

### How to configure Kibana

The Kibana default configuration is stored in [`kibana/config/kibana.yml`](./kibana/kibana.yml).

#### Kibana TLS

It's highly recommended to use Kibana with secure TLS connection. There is two ways to achieve that:

* Setup reverse proxy (like NGiNX).
* Setup Kibana using TLS itself.

You can find Kibana TLS setup instructions in [`tls/README.md`](./tls/)

Please refer to the following documentation page for more details about how to configure Kibana inside Docker
containers: [Install Kibana with Docker](https://www.elastic.co/guide/en/kibana/current/docker.html).

### How to configure Logstash

*:information_source: Do not use the `logstash_system` user inside the Logstash **pipeline** file, it does not have sufficient permissions to create indices. Follow the instructions at [Configuring Security in Logstash](https://www.elastic.co/guide/en/logstash/current/ls-security.html) to create a user with suitable roles.*

The Logstash configuration is stored in [`logstash/logstash.yml`](./logstash/logstash.yml), Logstash pipelines configuration is in [`logstash/pipelines.yml`](./logstash/pipelines.yml)

Please refer to the following documentation page for more details about how to configure Logstash inside Docker
containers: [Configuring Logstash for Docker](https://www.elastic.co/guide/en/logstash/current/docker-config.html).

### How to scale out the Elasticsearch cluster

Follow the instructions from the Wiki: [Scaling out Elasticsearch](https://github.com/deviantony/docker-elk/wiki/Elasticsearch-cluster)

### Healthcheck

Repo contains healthcheck bash scripts and utility buit with Go. You can choose one oh them or don't use service healthcheck.

#### Healthcheck Go utility

**Usage:** healthcheck [options] [elastic | kibana | logstash] [host]

By default tool configurated for default repo settings (https for elastic, default ports, ignoring invalid certs).

*:warning: Flags should be defore service type and host!*
* To use basic auth, add `-u <username`(Default remote_monitoring_user) and `-p <password>` flags.
* Trigger status can be setted with RegExp by `-s` flag, e.g: `healthcheck -f "-f 'green|yellow' elastic`
* Accept non default hostname/scheme, e.g: `healthcheck elastic http://elastic`

#### Healthcheck scripts

1. Add mount point for each script to corresponding service.
2. Change **`healthcheck: test: "CMD"`** to service healthcheck script.
3. Change checking endpoint and username/password.

## Extensibility
### How to add plugins

To add plugins to any Elastic stack component, you have to:

1. Create Dockerfile for service you want to apply plugin.
2. Add a `RUN` statement to the corresponding `Dockerfile` (e.g. `RUN logstash-plugin install logstash-filter-json`)
```dockerfile
# https://www.docker.elastic.co/
FROM docker.elastic.co/logstash/logstash:${LOGSTASH_VERSION}

# Add your logstash plugins setup here
RUN logstash-plugin install logstash-filter-json
```
3. Add the associated plugin code configuration to the service configuration (eg. Logstash input/output)
4. Add following to docker compose service section you want to apply plugin (e.g. Logstash):
```yaml
build:
      context: logstash/
```
5. (Re)Build the images using the `docker compose build` command

## JVM tuning
### How to specify the amount of memory used by a service

By default, both Elasticsearch and Logstash start with [1/4 of the total host memory](https://docs.oracle.com/javase/8/docs/technotes/guides/vm/gctuning/parallel.html#default_heap_size) allocated to the JVM Heap Size.

The startup scripts for Elasticsearch and Logstash can append extra JVM options from the value of an environment
variable, allowing the user to adjust the amount of memory that can be used by each component:

| Service       | Environment variable |
|---------------|----------------------|
| Elasticsearch | ES_JAVA_OPTS         |
| Logstash      | LS_JAVA_OPTS         |


For example, to increase the maximum JVM Heap Size for Logstash:

```yml
logstash:

  environment:
    LS_JAVA_OPTS: -Xmx1g -Xms1g
```

### How to enable a remote JMX connection to a service

As for the Java Heap memory (see above), you can specify JVM options to enable JMX and map the JMX port on the Docker host.

Update the `{ES,LS}_JAVA_OPTS` environment variable with the following content (I've mapped the JMX service on the port 18080, you can change that). Do not forget to update the `-Djava.rmi.server.hostname` option with the IP address of your Docker host (replace **DOCKER_HOST_IP**):

```yml
logstash:

  environment:
    LS_JAVA_OPTS: -Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.port=18080 -Dcom.sun.management.jmxremote.rmi.port=18080 -Djava.rmi.server.hostname=DOCKER_HOST_IP -Dcom.sun.management.jmxremote.local.only=false
```

## Going further
### Swarm mode

This time, there are no plans on support for Docker [Swarm mode](https://docs.docker.com/engine/swarm/).
