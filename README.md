# SideShell
 
The Side Shell is a control channel for **debug purposes only** in the environment where traditional console access is disabled. It comes as cli utility and Golang embedded server.


## Inspiration

Rarely while doing a system development you need to knock inside container(s), embedded system(s) or other platform(s) that runs your application code for audit purposes. Typically, the access becomes an issue either restricted by network topology, configuration of networking gears or other factors. The ability to establish inside-out shell session is required.

Let's imaging the following topologies
```
      10.x.x.x                54.x.x.x         192.168.x.x
[ { bash }-{ sshd } ] ---> [  tcp relay  ] <---- [ ssh ]

[ { bash }-{ sshd } ] <------------------------- [ ssh ]
      54.x.x.x                                 192.168.x.x
```


## Getting started

Embed the a side shell to your Golang application either using `sshd.Listen` for direct listening if your workload is hosted on public domain or `sshd.ViaProxy` if it is private. 

```go
import "github.com/fogfish/sideshell/sshd"

sshd.ViaProxy(
  os.Getenv("CONFIG_SSHD_PROXY"),
  sshd.Credentials(
    os.Getenv("CONFIG_SSHD_ACCESS"),
    os.Getenv("CONFIG_SSHD_SECRET"),
  ),
  sshd.PrivateKey(4096),
)
```

You can get the cli utility with

```bash
go get github.com/fogfish/sideshell
```


## Let's build the topology

1. Let's start the TCP relay, a component that accepts incoming connection both from ssh client and ssh daemon and routes a traffic between them. We do this on development laptop with following command:

```bash
sideshell relay --daemon :8080 --client :8022
```

2. The proxy must be reachable from outside. Usage of ngrok is an easiest way to do this. 

```bash
./ngrok tcp 8080
```

3. Build the docker with example workload

```bash
docker build --no-cache -t sideshell - < example/Dockerfile
```

4. Run the container, it is mandatory to provide a few parameters: ngrok endpoint and desired access/secret to access the shell 

```bash
docker run \
  -e CONFIG_SSHD_PROXY=2.tcp.ngrok.io:10768 \
  -e CONFIG_SSHD_ACCESS=foo \
  -e CONFIG_SSHD_SECRET=bar \
  sideshell
```

5. Connect to your shell, use credentials defined at first step

```bash
ssh -p 8022 \
  -o "UserKnownHostsFile=/dev/null" \
  -o "StrictHostKeyChecking=no" \
  foo@localhost 
```

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/sideshell.svg?style=for-the-badge)](LICENSE)