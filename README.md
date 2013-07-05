# Minuswell

[![Build Status](https://travis-ci.org/ohlol/minuswell.png)](https://travis-ci.org/ohlol/minuswell)

[minus well](http://eggcorns.lascribe.net/english/129/minus/)

Minuswell is a log file shipper for Logstash, written in Go. It supports a number of outputs (and can multiplex them).

## Building

```
$ go get -d -v ./... && go build -v ./...
```

The `go build` command creates the `minuswell` binary for you to install/run.

### Notes

* The fsnotify library is currently slightly broken with Linux, so you must use commit  `d66819e17205a446430d34a9ba41625cec35be19`.
* If you get errors with gozmq, check [these instructions](https://github.com/alecthomas/gozmq#installing).

## Usage

```
Usage:
  minuswell [OPTIONS]

Help Options:
  -h, --help              Show this help message

Application Options:
  -c FILE                 Path to config file
  -d, --config-dir DIR    Parse config files in dir
  -o OUTPUT               Which output to use (can specify multiple)
```

You must specify both config and output. To send to multiple outputs, simply add more `-o` options to the args.

If you'd like to have minuswell parse multiple config files, it can do that too. Just use `--config-dir=path/to/dir`.

It uses the [shoenice](https://github.com/ohlol/shoenice) library to provide statistics about the app and send them to Graphite.

## Supported outputs

* `pipe` (stdout) no configuration needed
* `tcp` - requires specifying `address` and `port`
* `zmq` - set `addresses` which is an array of ZeroMQ addresses

## Configuration

```
{
    "graphiteHost": "graphite.foo.com",
    "outputs": {
        "tcp": {
            "address": "127.0.0.1",
            "port": 1234
        },
        "zmq": {
            "addresses": ["tcp://127.0.0.1:2120"]
        }
    },
    "files": {
        "/var/log/*.log": {
            "type": "syslog",
            "tags": ["system", "logs"],
            "fields": {
                "field1": "val1"
            }
        }
    }
}
```
