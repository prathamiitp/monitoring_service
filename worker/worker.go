package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
)

//ctx required in redis
var ctx = context.Background()

//holding the radis conncection
var redisdb *redis.Client

//holding the mysql connection
var mysqldb *sql.DB

//decide the no of request to add to db at a time
//or the time after which any no(<=max) of requests should be added
var maxReq int = 1000 //amount of request
var maxtime int = 1 //no of min

//wait group to wait for db connections
var wg sync.WaitGroup

//log struct to be saved in db and in redis
type Logs struct {
	Time time.Time `json:"time"`
	Url  string    `json:"url"`
}

//connecting to redis server at port 6379 localhost
func connectRedis() {
	defer wg.Done()

	redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("redisConnection: %v", err))
	}
	fmt.Println("Redis connected at localhost:6379")
}

//connect to mysql server on localhost:3306 and open the db 'requestLogs' to work with, print error if any
func connectMysql() {
	defer wg.Done()

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/requestLogs")
	if err != nil {
		panic(fmt.Errorf("mysqlConnection: %v", err))
	}
	mysqldb = db
	fmt.Println("Mysql connected at localhost:3306")
}

//logs data to db after 'max' amount of logs are accumulated
//runs infinitely
func logToDb(c chan os.Signal) {
	for {
		start := time.Now()
		var arg []interface{}
		var Log Logs
		var a []byte
		var err error

		query := "INSERT INTO initialLogs (time, url) VALUES "
		for i := 0; i < maxReq && int(time.Since(start).Minutes())<maxtime; {
			a, err = redisdb.RPop(ctx, "reqlist").Bytes()
			if err != nil {
				continue
			}
			json.Unmarshal(a, &Log)
			query += "(?, ?),"
			arg = append(arg, Log.Time, Log.Url)
			i++
		}
		if len(arg)==0 {
			continue
		}

		dblog, err := mysqldb.Exec(query[:len(query)-1], arg...)
		if err != nil {
			fmt.Println("dbLogging error: ", err)
		}
		id, _ := dblog.LastInsertId()
		fmt.Println("last set of requests logged in db from id: ", id)
		
		// stop the loop if a interrupt is recorded;
		test:=false
		select{
			case<-c:
				test=true;
			default:
				continue;
		}
		if test {
			break;
		}
	}
}

func main() {
	wg.Add(2)
	go connectRedis()
	go connectMysql()

	// to keep from running the rest of the programn before connecting to db
	wg.Wait()

	// to keep code from suddenly exiting, it logs the last processing set of data
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	//running the data logger infinitely
	logToDb()
}
