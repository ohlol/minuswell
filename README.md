# Minuswell

Minuswell is a log file shipper for Logstash, written in Go. It supports a number of outputs (and can multiplex them).

## Building

```
$ go get github.com/jessevdk/go-flags
$ go get github.com/ActiveState/tail
$ go get github.com/howeyc/fsnotify
$ go get github.com/alecthomas/gozmq
$ go build
```

## Usage

```
Usage:
  minuswell [OPTIONS]

Help Options:
  -h, --help= Show this help message

Application Options:
  -c=         path to config file
  -o=         which output to use
```

You must specify both config and output. To send to multiple outputs, simply add more `-o` options to the args.

## Configuration

```
{
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
        "/varlog/*.log": {
            "type": "syslog",
            "tags": ["system", "logs"],
            "fields": {
                "field1": "val1"
            }
        }
    }
}
```
