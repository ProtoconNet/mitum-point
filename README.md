### mitum-point

*mitum-point* is a [mitum](https://github.com/ProtoconNet/mitum2)-based contract model and is a service that provides mint functions.

#### Installation

```sh
$ git clone https://github.com/ProtoconNet/mitum-point

$ cd mitum-point

$ go build -o ./mp ./main.go
```

#### Run

```sh
$ ./mp init --design=<config file> <genesis config file>

$ ./mp run --design=<config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
[genesis-design.yml](genesis-design.yml) is a sample of `genesis config file`.