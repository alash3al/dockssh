# Dockssh
Dockssh, ssh into any container from anywhere

# Why
- For learning and fun
- For testing/staging/development environments

# How it works
- `Dockssh` runing on port `22022` on host `example.com`
- A user connects to `dockssh` i.e `ssh -p 22022 container1@example.com`
- `Dockssh` checks if the user provided password is the same as the one stored in redis key `dockssh:container1:pass`
- On success, `Dockssh` will open a `PTY` (pseudotty) to `docker exec -it container1 /bin/sh`
- Have fun ^^! (replace `container1` with any of your containers)

# Why redis for configurations
- No configurations files
- Simple & tiny
- Makes `Dockssh` loads configurations in realtime, no need to restart

# Requirements
- Linux
- Docker
- Redis

# Downloads
Download the binary from [here](https://github.com/alash3al/dockssh/releases/tag/v1.0.0)

# Building from source
`go get github.com/alash3al/dockssh`

# Usage
`./dockssh --help`
