# consul-join to help join consuls

![Build status](https://travis-ci.org/wakeful-deployment/consul-join.svg?branch=master)

## What is this for?

`consul-join` can help figure out which IPs to use to join consul agents
to servers. It respects two ENV variables related to joining.

## Usage

Simple join with one IP:

```sh
JOINIP=10.0.0.1 consul-join
# will exec: consul agent -join=10.0.0.1
```

Join by looking up registered A records:

```sh
JOINDNS=example.com consul-join
# assuming example.com has 2 A records of 10.0.0.1 and 10.0.0.2, then
# will exec: consul agent -join=10.0.0.1 -join=10.0.0.2
```

Passes all following arguments to `consul`:

```sh
consul-join -config-dir=/config
# will exec: consul agent -config-dir=/config
```

Also supports `-server` and will respect the `BOOTSTRAP_EXPECT` ENV
varaible (which it defaults to 3 if missing).

```sh
JOINDNS=example.com consul-join -server
# will exec: consul agent -join=10.0.0.1 -join=10.0.0.2 -bootstrap-expect=3 -server
```

## Roadmap

We are considering adding a `JOINURL` and a related `consul-join-server`
to provide a simple alternative to DNS for keeping track of consul
server ips.
