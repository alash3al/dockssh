# Dockssh
Dockssh, ssh into any container from anywhere

# Why
- For learning and fun
- Wasting some free time in my weekend :D
- For testing/staging/development environments

# How it works
- `Dockssh` running on port `22022` on host `example.com`
- A user connects to `container1` using `dockssh` from remote computer i.e `ssh -p 22022 container1@example.com`
- `Dockssh` checks if the user provided password is the same as the one stored in redis key `dockssh:container1:pass`
- On success, `Dockssh` will open a `PTY` (pseudotty) to `docker exec -it container1 /bin/sh`

# Why redis for configurations
- No configurations files
- Simple & tiny
- Makes `Dockssh` loads configurations in realtime, no need to restart

# Requirements
- Linux
- Docker
- Redis

# Downloads
Download the binary from [here](https://github.com/alash3al/dockssh/releases)

# Building from source
You need to get the dependencies using the command:
`go get github.com/alash3al/dockssh`

# Usage
<strong>On the host machine:</strong>
- Install [Redis](https://redis.io/) using the commands:<br/>
    Debian: `sudo apt install redis`<br/>
    RHEL: `sudo yum install redis`
- Create a container for testing, I will name it `TestCont`:<br/>
    `sudo docker create --name TestCont -it ubuntu:latest bash`
- Start the container:<br/>
    `sudo docker start TestCont`
- Set a password for the container over SSH:<br/>
    `redis-cli set dockerssh:TestCont:pass "mypass"`
- Download the latest `Dockssh` binary from [here](https://github.com/alash3al/dockssh/releases).
- Rename the file to `dockssh`.
- Make it executable:<br/>
    `chmod 775 dockssh`
- Make sure to open the port in the firewall:<br/>
    `sudo ufw allow 22022`
- Run the server:<br/>
    `./dockssh`
- You should see a message:<br/>
    `Now listening on port: 22022`

<strong>On the remote machine:</strong>
- Connect to your container:<br/>
    `ssh TestCont@host_ip_address -p 22022`
- Enter `yes`.
- Enter your password and press Enter.
