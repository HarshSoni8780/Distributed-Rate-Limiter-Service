package limiter

import(
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//create a struct named FixedWindow for database redis and limit of client
type FixedWindow struct{
	Rdb *redis.Client
	Limit int
}
// it calculates counter of every window
var ctx = context.Background()

//in Go, it doesnt have class or "new" like in cpp, so use struct return in function like this:
func NewFixedWindow(rdb *redis.Client, limit int) *FixedWindow{
	return &FixedWindow{
		Rdb: rdb,
		Limit: limit,
	}
}

//every time user request Allow func is called to check if user is within cuurent window or not
func (fw * FixedWindow) Allow(userID string) bool{
	window := time.Now().Unix()/60

	//eg: rate:Harsh:29623400.
	key := fmt.Sprintf("rate:%s%d", userID, window)

	//Incr cmd is Atomid which allows redis to process every increment one by one and not simantaneously
	// key is passed to redis db to check its value, if exist add 1 else if doest not exist then value = 1.
	//count = current number of request 
	count,err := fw.Rdb.Incr(ctx,key).Result()
	if err != nil{
		return false
	}

	//adding count == 1 ensures that we only check new request on every new 60 sec or 1 min and not for every 
	//request which reduces cpu cycles
	if count == 1{
		fw.Rdb.Expire(ctx, key, time.Minute)
	}
	
	return count <= int64(fw.Limit)

	//eg limit = 2
	//req A -> count = 1 using Incr set timer 60 sec
	//1<=2 accept req A

	//req B -> count = 2 using Incr 
	//2<=2 accept req B

	//req C ->count = 3 using Incr
	//3<=2 wrong so req failed, User is now rate-limited.
}