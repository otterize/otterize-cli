# Otterize CLI

![build](https://github.com/otterize/network-mapper/actions/workflows/build.yaml/badge.svg)
![go report](https://img.shields.io/static/v1?label=go%20report&message=A%2B&color=success)
[![community](https://img.shields.io/badge/slack-Otterize_Slack-purple.svg?logo=slack)](https://joinslack.otterize.com)

[About](#about) | [Installation & usage with the network mapper](#installation-instructions--usage-with-the-network-mapper) | [Docs](https://docs.otterize.com/cli/) | [Contributing](#contributing) | [Slack](#slack)

## About

The Otterize CLI is a command-line utility used to control and interact with the [Otterize network mapper](https://github.com/otterize/network-mapper), manipulate local intents files, and interact with Otterize Cloud.

See the [CLI command reference](https://docs.otterize.com/cli/) for how to use it.

Example output from running the network mapper on the [Google Cloud microservices demo](https://github.com/GoogleCloudPlatform/microservices-demo):
```bash
$ otterize mapper list
cartservice in namespace ecommerce calls:
  - redis-cart
checkoutservice in namespace ecommerce calls:
  - kafka-secure
frontend in namespace ecommerce calls:
  - adservice
  - cartservice
  - checkoutservice
  - currencyservice
  - productcatalogservice
  - recommendationservice
  - shippingservice
kafka-secure in namespace ecommerce calls:
  - kafka-secure
  - lab-zookeeper
paymentservice in namespace ecommerce calls:
  - kafka-secure
recommendationservice in namespace ecommerce calls:
  - productcatalogservice
```

## Installation instructions & usage with the network mapper
### Install the network mapper using Helm
```bash
helm repo add otterize https://helm.otterize.com
helm repo update
helm install network-mapper otterize/network-mapper -n otterize-system --create-namespace --wait
```
### Install Otterize CLI to query data from the network mapper
Mac
```bash
brew install otterize/otterize/otterize-cli
```
Linux 64-bit
```bash
wget https://get.otterize.com/otterize-cli/v0.1.20/otterize_Linux_x86_64.tar.gz
tar xf otterize_Linux_x86_64.tar.gz
sudo cp otterize /usr/local/bin
```
Windows
```bash
scoop bucket add otterize-cli https://github.com/otterize/scoop-otterize-cli
scoop update
scoop install otterize-cli
```
For more platforms, see [the installation guide](https://docs.otterize.com/installation#install-the-otterize-cli).


## Contributing
1. Feel free to fork and open a pull request! Include tests and document your code in [Godoc style](https://go.dev/blog/godoc)
2. In your pull request, please refer to an existing issue or open a new one.
3. See our [Contributor License Agreement](https://github.com/otterize/cla/).

## Slack
[Join the Otterize Slack!](https://joinslack.otterize.com)
