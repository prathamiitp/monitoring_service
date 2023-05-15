<div align="center">  
    
# Request Monitoring Microservice
[<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white"/>](https://go.dev/doc/)
[<img src="https://img.shields.io/badge/MySQL-4479A1.svg?style=for-the-badge&logo=MySQL&logoColor=white"/>](https://docs.oracle.com/en-us/iaas/mysql-database/index.html)
[<img src="https://img.shields.io/badge/Redis-DC382D.svg?style=for-the-badge&logo=Redis&logoColor=white"/>](https://redis.io/docs/getting-started/)

  ---
  
</div>  

## Monitoring service
This request monitoring service logs the time the request was made and other request related details to a MySQL database.  
Listener and Worker are seperate code and share data using Redis **in-memory** database (Redis Lists).  
The Listener handles request concurrently and save logs to Redis Lists acting as Queue, and Worker proccess upto 20,000req/min and log them in bulk on MySQL database clearing Redis Queue.  
The Worker code can handle Graceful shutdown without loosing any relavent data.  

## Steps to setup the monitoring service
#### 1. Clone repo monitoring_service : 
```
git clone https://github.com/prathamiitp/monitoring_service
```
#### 2. Setup Golang : 
[Installation and documentation](https://go.dev/doc/install)
#### 3. Set GOPATH and GOROOT variables by adding following lines in `file: .bashrc` for ubuntu inside HOME folder : 
```
export GOPATH=$HOME/go
export GOROOT=/usr/local/go/
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
```
#### 4. Setup Database servers : 
  - First Setup MySQL and Redis system servers from database-server-setup-for-system mentioned below.
#### 5. Run Both database servers on their default port setting : 
```
sudo systemctl start mysql
```
```
sudo systemctl start redis-server
```
To check if the server is up and running, just replace 'start' with 'status'  
To stop once work is done, replace 'start' with 'stop'
## Database Server setup for system
- [MySQL](https://www.mysqltutorial.org/install-mysql-ubuntu/)
- [Redis](https://redis.io/docs/getting-started/installation/install-redis-on-linux/) (use the APT repository for it includes redis-cli)
Remeber to add the database for both Redis and MySQL with the schema and name used in code.

## Database Driver to be used in Go
- MySQL: [[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)] , [[database/sql](https://pkg.go.dev/database/sql)]
- Redis: [[github.com/go-redis/redis/v9](https://pkg.go.dev/github.com/go-redis/redis/v9@v9.0.0-beta.1#section-readme)] [(Github repo)](https://github.com/go-redis/redis)
  - Also check [Uptrace's documentation](https://redis.uptrace.dev/guide/)

## Work to be done
- optimize cpu utilization (utilization is high because of infinite loop in worker which make connection to redis to check for any new datqa in queue)
- cleaner code implimentation (break logging function from worker code into smaller function, also improve the implimentation of graceful shutdown)
- upload go mod file and create a shell script for all the setup
- update readme with testing method for the monitoring service.
- reduce redis connection in worker, try to take up data in bulk from redis(just like done in case of mysql) which is taken unit by unit for now 
