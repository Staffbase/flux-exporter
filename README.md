⚠️ **This repository is archived and is no longer maintained.**

# Flux Exporter

The Flux Exporter can be used to export additional metrics for Flux. It exports all Docker images with there current running version and the last filtered image version. If the current and the new image versions do not match the value for the metric is `1`. When they are equal the value is `0`.

```
flux_exp_image{current="0.1.8",name="fluent-bit-template",namespace="elastic-system",new="0.1.8"} 0
flux_exp_image{current="0.30.0",name="nginx-ingress-controller",namespace="ingress-nginx",new="0.31.1"} 1
flux_exp_image{current="0.5.18-debian-9-r5",name="external-dns",namespace="external-dns",new="0.7.1"} 1
flux_exp_image{current="1.0.0-rc9",name="flux-helm-operator",namespace="flux",new="1.0.1"} 1
flux_exp_image{current="1.18.0",name="flux",namespace="flux",new="1.19.0"} 1
flux_exp_image{current="1.4.2",name="fluent-bit",namespace="fluent-bit",new="arm32v7-1.4.2"} 1
flux_exp_image{current="1.4.5",name="vault-secrets-operator",namespace="vault-secrets-operator",new="1.4.7"} 1
flux_exp_image{current="6.7.3",name="grafana",namespace="grafana",new="master-ubuntu"} 1
flux_exp_image{current="7.6.1",name="elasticsearch",namespace="elastic-system",new="8.0.0-SNAPSHOT"} 1
flux_exp_image{current="7.6.1",name="kibana",namespace="elastic-system",new="8.0.0-SNAPSHOT"} 1
```

## Usage

Clone the repository and build the binary:

```sh
git clone git@github.com:Staffbase/flux-exporter.git
make build
```

Set the environment variables for the Flux API endpoint:

```sh
export ENDPOINT=
```

Then run the exporter:

```sh
./bin/fluxexporter
```

A Docker image is available at `registry.staffbase.com/public/flux-exporter:<TAG>` and can be retrieved via:

```sh
docker pull registry.staffbase.com/public/flux-exporter:<TAG>
```
