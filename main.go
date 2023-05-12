package main

import (
	"fmt"
	"metrics-consumer/consumer"
	"metrics-consumer/mongo_client"
	"metrics-consumer/redis_client"
)

func main() {
	fmt.Println("Starting metrics consumer...")

	redisClient := redis_client.NewRedisClient()
	if redisClient == nil {
		fmt.Println("Failed creating redis client, cron cannot start!")
		return
	}

	mongoClient, mongoErr := mongo_client.NewMongoClient()
	if mongoErr != nil {
		fmt.Println(mongoErr.Error())
		return
	}

	consumer.ConsumeMetrics(redisClient, mongoClient)

	mongoClient.Disconnect()
	fmt.Println("Metrics consumer finished")
}
