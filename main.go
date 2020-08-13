package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func main() {
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	err := setStruct(conn)
	if err != nil {
		fmt.Println(err)
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

// set executes the redis SET command
func set(c redis.Conn) error {
	_, err := c.Do("SET", "Favorite Movie", "Repo Man")
	if err != nil {
		return err
	}
	_, err = c.Do("SET", "Release Year", 1984)
	if err != nil {
		return err
	}
	return nil
}

// get executes the redis GET command
func get(c redis.Conn) error {

	// Simple GET example with String helper
	key := "Favorite Movie"
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	// Simple GET example with Int helper
	key = "Release Year"
	i, err := redis.Int(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %d\n", key, i)

	// Example where GET returns no results
	key = "Nonexistent Key"
	s, err = redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		fmt.Printf("%s does not exist\n", key)
	} else if err != nil {
		return err
	} else {
		fmt.Printf("%s = %s\n", key, s)
	}

	return nil
}

func setStruct(c redis.Conn) error {

	const objectPrefix string = "12345"
	// {"ispro":true,"userid":"4","username":"Andi","playerid":""}"
	usr := User{
		Ispro:    true,
		UserID:   "4",
		Username: "agnynureza",
		PlayerID: "",
	}

	// serialize User object to JSON
	json, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", objectPrefix, json)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Ispro    bool   `json:"ispro"`
	UserID   string `json:"userid"`
	Username string `json:"username"`
	PlayerID string `json:"playerid"`
}

// {"loggedInUser":"4",
// 	"chname":{
// 		"C":["BBRI","TLKM","BBCA","BBNI","MNCN","PTBA","BRIS","ICBP","ASII","TKIM","IHSG"],
// 		"S":["IHSG"],"O":["IHSG"],
// 		"T":[],"H2":false}
// 	}

// 	publish LIVEUPDATE3-IHSG-3 "#C|IHSG|1596507790|2020-08-04 09:23:10|5027.771|5006.223|5054.663|5006.499|5027.771|1671291494|5027.771|21.548|0.43||326340607179|26085843891213|336171294579|127033|0|1330012382479|0|0|0|0|5006.223|0|1596507790"
