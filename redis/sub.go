package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func Sub(ch chan int) {
	fmt.Println("ABC")
	conn := Pool.Get()
	for {
		// Get a connection from a pool
		psc := redis.PubSubConn{Conn: conn}

		// Set up subscriptions
		psc.Subscribe("redisChat")

		// While not a permanent error on the connection.
		for conn.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Println(string(v.Data))
				Set(string(v.Data), string(v.Data))
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println()
			}
		}
		conn.Close()
	}
	ch <- 1
}
