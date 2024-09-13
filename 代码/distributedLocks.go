/* 
在 MongoDB 中，_id 字段并不是自增的，也没有固定的生成方式。实际上，_id 是文档的唯一标识符，默认情况下 MongoDB 会为其分配一个唯一的 ObjectId，但你也可以使用自定义的值（如字符串、数字、UUID 等）来代替默认的 ObjectId。

因此，将 lockKey 设置为 _id 字段是合理的做法。如果锁的键是唯一的，使用 lockKey 作为 _id 可以确保锁的唯一性，MongoDB 会自动利用 _id 的唯一性约束，确保同一个锁不能被重复创建。


在 Go 语言中，<-ctx.Done() 是一种常见的模式，用于监听上下文 (context.Context) 的取消信号。

解释 <-ctx.Done()
ctx.Done() 是一个返回只读 channel 的方法，该 channel 会在以下情况下关闭：

上下文被取消：通过调用 cancel() 函数（例如使用 context.WithCancel 创建的上下文）。
上下文的超时：如果上下文是通过 context.WithTimeout 或 context.WithDeadline 创建的，当超时或到达指定的时间时，该 channel 会关闭。
一旦 ctx.Done() 返回的 channel 关闭，监听它的 goroutine 就会从该 channel 读取到一个值，这意味着该操作应该终止或做相应的清理。

*/
package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 锁
type Lock struct {
	Key        string    `bson:"_id"`        //唯一键
	Holder     string    `bson:"Holder"`     //持有者
	Expiration time.Time `bson:"Expiration"` //锁的过期时间
}

func ConnectToMongodb() (*mongo.Client, error) {
	options := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	//test connection
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func AcquireLock(client *mongo.Client, key, holder string, ttl time.Duration) (bool, error) {
	collection := client.Database("locks_db").Collection("locks")

	now := time.Now()
	expiration := now.Add(ttl)

	filter := bson.M{
		"_id": key,
		"$or": bson.A{
			bson.M{"expiration": bson.M{"$lte": now}},
			bson.M{"holder": holder},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"holder":     holder,
			"expiration": expiration,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)

	var result Lock
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReleaseLock(client *mongo.Client, key, holder string) error {
	collection := client.Database("locks_db").Collection("locks")

	// 仅当持有者是当前holder时, 才能解锁
	filter := bson.M{
		"_id":    key,
		"holder": holder,
	}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func StartAutoRenewal(ctx context.Context, cancel context.CancelFunc, client *mongo.Client, key, holder string, ttl time.Duration, renewalInterval time.Duration,
	renewalFailed chan<- bool,
) {
	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := AcquireLock(client, key, holder, ttl)
			if err != nil {
				renewalFailed <- true
				cancel()
				return
			}

			if ok {
				//renewal successful
			} else {
				renewalFailed <- true
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func Test(client *mongo.Client, name string) {

}

func main() {
	client, err := ConnectToMongodb()
	if err != nil {
		fmt.Println("Failed to connect to mongodb: ", err)
		return
	}
	defer client.Disconnect(context.TODO())

	key := "Test_distributed_lock"
	holder := "test1"
	ttl := 1 * time.Minute
	renewalInterval := 1 * time.Second

	ok, err := AcquireLock(client, key, holder, ttl)
	if err != nil {
		fmt.Println("Failed to acquire Lock:", err)
		return
	}

	if ok {
		// 管理取消操作的context
		ctx, cancel := context.WithCancel(context.Background())

		renewalFailed := make(chan bool)

		go StartAutoRenewal(ctx, cancel, client, key, holder, ttl, renewalInterval, renewalFailed)

		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Main operation canceled due to lock renewal failure")
					return
				default:
					fmt.Println("Main operation is running....")
					time.Sleep(2 * time.Second)
				}
			}
		}()

		<-ctx.Done()

		if err = ReleaseLock(client, key, holder); err != nil {
			fmt.Println("Failed to release lock")
		}
	} else {
		fmt.Println("Failed to acquire lock. Someone hold it")
	}

	time.Sleep(5 * time.Minute)
}
