// Dockssh is a tiny tool that exposes your docker container to the world
// behind a ssh server, so you can access any of them using your ssh client.
// It uses redis as password storage for the containers.
// Each container has a name, and each name has a pass,
// it uses this format for redis keys to store passwords (dockssh:$CONTAINER_NAME:pass).
package main

import (
	"log"
)

import (
	"github.com/gliderlabs/ssh"
	"github.com/go-redis/redis"
)

func main() {
	opts, err := redis.ParseURL(*flagRedis)
	if err != nil {
		log.Fatal(err)
	}

	redisConn = redis.NewClient(opts)
	if _, err := redisConn.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	ssh.Handle(handler)

	sshOpts := []ssh.Option{
		ssh.PasswordAuth(checkPassword),
	}

	if *flagHostKeyFile != "" {
		sshOpts = append(sshOpts, ssh.HostKeyFile(*flagHostKeyFile))
	}

	log.Fatal(ssh.ListenAndServe(
		*flagListenAddr,
		nil,
		sshOpts...,
	))
}
