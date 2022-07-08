package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
)  

//ctx required in redis
var ctx = context.Background()
//holding the radis conncection
var rdb *redis.Client


//log struct to be saved in db and in redis
type Logs struct {
	Time time.Time	`json:"time"`
	Url  string		`json:"url"`
}

//connecting to redis server
func connectRedis(){
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err!=nil{
		panic(fmt.Errorf("redisConnection: %v", err))
	}
	fmt.Println("Redis connected at localhost:6379")
}

//
func headers(w http.ResponseWriter, req *http.Request) {
	reqUrl := req.URL.String()
	value,_ := json.Marshal(Logs{time.Now(), reqUrl})
	rdb.LPush(ctx, "reqlist", value)
	fmt.Fprint(w, "served")
}

func main() {
	//connecting to redis (in memory database)
	connectRedis()

	//to handle request at /main/ endpoint
	http.HandleFunc("/main/", headers)
	fmt.Println("listening http request at localhost:8090")
	//to listen to port 8090
	http.ListenAndServe(":8090", nil)
}
