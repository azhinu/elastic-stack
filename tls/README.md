# TLS certificates

**:warning: All commands assume you are inside the `tls` directory of the repository.**

##TL;DR

Modify `instances.yml` if needed and run inside `tls` directory.
```shell
docker run --rm -it \
  -v "$PWD":/usr/share/elasticsearch/tls \
  docker.elastic.co/elasticsearch/elasticsearch:8.1.0 \
  elasticsearch-certutil ca \
    --silent \
    --pem \
    --out tls/ca.zip \
  && unzip tls/ca.zip \
  && echo "" | elasticsearch-certutil cert \
    --silent \
    --ca-cert tls/ca/ca.crt \
    --ca-key tls/ca/ca.key \
    --in tls/instances.yml \
    --out tls/certificate-bundle.zip \
  && unzip tls/certificate-bundle.zip \
  && rm tls/ca.zip tls/certificate-bundle.zip
```

## Certificate Authority and Elasticsearch node TCP certificate

Generate a bundle containing the Certificate Authority (CA) in PEM format and a certificate/key for Elasticsearch (PKCS#12) nodes communication, using the `elasticsearch-certutil` tool that ships with Elasticsearch:

```shell
./certutil.sh ca \
    --pem \
    --out tls/ca.zip \
&& unzip ca.zip && rm -f ca.zip

./certutil.sh cert \
    --ca-cert tls/ca/ca.crt \
    --ca-key tls/ca/ca.key \
    --in tls/instances.yml \
    --out tls/certificate-bundle.zip \
&& unzip certificate-bundle.zip && rm -f certificate-bundle.zip
```

You will be prompted to enter an optional password to protect Elasticsearch keys. Please be aware that
the password you enter here, if not empty, will have to be **added to the Elasticsearch keystore on every node** using
the [`elasticsearch-keystore`][es-keystore] CLI tool (see below)

**:warning: `<none>` indicates that the value can be left empty.**

```none
Enter password for elasticsearch/elasticsearch.p12: <none>
```

> :information_source: In case you entered a password above, execute the following 2 additional commands (Elasticsearch
> *must* be running):
>
> ```console
> $ docker compose run -T elastic bin/elasticsearch-keystore add xpack.security.transport.ssl.keystore.secure_password
> Enter value for xpack.security.transport.ssl.keystore.secure_password: <password>
> ```
>
> ```console
> $ docker compose run -T elastic bin/elasticsearch-keystore add xpack.security.transport.ssl.truststore.secure_password
> Enter value for xpack.security.transport.ssl.truststore.secure_password: <password>
> ```

At this stage, the structure of the `tls` directory should be:

```tree
├── ca
│   ├── ca.crt
│   └── ca.key
├── certutil.sh
├── elasticsearch
│   └── elasticsearch.p12
├── instances.yml
└── README.md
```

## Elasticsearch HTTP certificate and CA PEM certificate
*:information_source: Now it's optional. You can safely use certificate from previous step.*


</br>
Using the same `elasticsearch-certutil` as above, generate a set of *elasticsearch HTTP* and *CA PEM* certificates that can be used by other
components (Logstash, Kibana, ...) to communicate with Elasticsearch over HTTPS:

```bash
$ ./certutil.sh http
```

**:warning: `<none>` indicates that the value is to be left empty.**

```none
Generate a CSR? [y/N] n
Use an existing CA? [y/N] y
CA Path: /usr/share/elasticsearch/tls/ca/ca.p12
Password for ca.p12: <none>
For how long should your certificate be valid? [5y] <none>
Generate a certificate per node? [y/N] n
(Enter all the hostnames that you need, one per line.)
elastic
localhost
Is this correct [Y/n] y
(Enter all the IP addresses that you need, one per line.)
<none>
Is this correct [Y/n] y
...
Do you wish to change any of these options? [y/N] n
Provide a password for the "http.p12" file: <none>
What filename should be used for the output zip file? tls/elasticsearch-ssl-http.zip
```

Extract the generated certificates:

```console
$ sudo unzip elasticsearch-ssl-http.zip
Archive:  elasticsearch-ssl-http.zip
  inflating: elasticsearch/README.txt
  inflating: elasticsearch/http.p12
  inflating: elasticsearch/sample-elasticsearch.yml
   creating: kibana/
  inflating: kibana/README.txt
  inflating: kibana/elasticsearch-ca.pem
  inflating: kibana/sample-kibana.yml
```

You can then safely remove the Zip file generated by the `elasticsearch-certutil` command in the previous step:

```console
$ sudo rm elasticsearch-ssl-http.zip
```

At this stage, the structure of the `tls` directory should be:

```tree
tls
├── ca
│   ├── ca.crt
│   └── ca.key
├── elasticsearch
│   ├── elasticsearch.p12
│   ├── http.p12
├── instances.yml
└── README.md
```

## Kibana TLS

Optionally, using `elasticsearch-certutil` you can also make certificates to set up TLS for Kibana:

:information_source: Set kibana access hostname with `--dns` argument. You can set multiple domains separated by commas. Also, if Kibana used without domain name, IP address can be added with `--ip` argument.

```bash
$ ./certutil.sh csr \
  --name kibana \
  --dns kibana.external.domain.com,your-alt-domain.com \
  --out tls/kibana-csr.zip
```

Extract the generated certificates:

```console
$ sudo unzip kibana-csr.zip
Archive:  kibana-csr.zip
  inflating: kibana/kibana.csr       
  inflating: kibana/kibana.key
```

You can then safely remove the Zip file generated by the `elasticsearch-certutil` command in the previous step:

```console
$ sudo rm kibana-csr.zip
```

At this stage, the structure of the `tls` directory should be:

```tree
tls
├── ca
│   ├── ca.crt
│   └── ca.key
├── elasticsearch
│   ├── elasticsearch.p12
│   ├── http.p12
├── instances.yml
└── README.md
tls
├── ca
│   └── ca.p12
├── elasticsearch
│   ├── README.txt
│   ├── elasticsearch.p12
│   ├── http.p12
│   └── sample-elasticsearch.yml
├── instances.yml
├── kibana
│  ├── kibana.csr
│  └── kibana.key
└── README.md
```

>**:warning: Uncomment options related to Kibana TLS setup:**
>
>docker compose.yml:
>
>```yaml
>kibana:
>  image: docker.elastic.co/kibana/kibana:7.16.1
>  volumes:
>    - ./kibana/kibana_healthcheck.sh:/usr/local/bin/elastic_healthcheck:ro #Kibana healthcheck script
>    - ./tls/kibana/elasticsearch-ca.pem:/usr/share/kibana/config/ca.crt:ro #Elastic CA
>    - ./tls/kibana/kibana.csr:/usr/share/kibana/config/kibana.csr:ro # Kibana CRT
>    - ./tls/kibana/kibana.key:/usr/share/kibana/config/kibana.key:ro # Kibana CRT Key
>   environment:
>    SERVER.SSL.CERTIFICATE: config/kibana.csr
>    SERVER.SSL.KEY: config/kibana.key
>```
>


## Verification

(Re)start the stack using `docker compose [up|restart]` to load the certificates generated in the steps described
earlier in this document.

To verify that stack components are successfully communicating with Elasticsearch over TLS, ensure the following
messages are logged to the console.

Logstash:

```log
[INFO ][logstash.outputs.elasticsearch][main] Elasticsearch pool URLs updated {:changes=>{:removed=>[], :added=>[https://logstash-es:xxxxxx@elasticsearch:9200/]}}
[WARN ][logstash.outputs.elasticsearch][main] Restored connection to ES instance {:url=>"https://logstash-es:xxxxxx@elasticsearch:9200/"}
[INFO ][logstash.outputs.elasticsearch][main] ES Output version determined {:es_version=>7}
```

Kibana:

```json
{"type":"log",...,"tags":["status","plugin:elasticsearch@7.15.0","info"],...,"message":"Status changed from yellow to green - Ready","prevState":"yellow","prevMsg":"Waiting for Elasticsearch"}
```