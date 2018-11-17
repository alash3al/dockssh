package main

import (
	"flag"
	"os"
	"path"

	"github.com/go-redis/redis"
)

var (
	flagListenAddr  = flag.String("listen-addr", "0.0.0.0:22022", "the ssh listening address")
	flagHostKeyFile = flag.String("host-key", path.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "the host key path for persistence")
	flagEntryPoint  = flag.String("entrypoint", "/bin/sh", "the default command to execute")
	flagRedis       = flag.String("redis", "redis://localhost:6379/0", "the redis dsn (redis://user:password@host:port/db)")
)

var (
	redisConn          *redis.Client
	redisDocksshPrefix = "dockssh"
)

func init() {
	flag.Parse()
}
