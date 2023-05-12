package consumer

import (
	"encoding/json"
	"fmt"
	"metrics-consumer/mongo_client"
	"metrics-consumer/redis_client"
	"time"
)

type aggregatedMetricDocument struct {
	MetricName string    `json:"metric"`
	Count      int       `json:"count"`
	Timestamp  time.Time `json:"timestamp"`
}

const SystemMetricsQueueRedis = "system-metrics"

func ConsumeMetrics(redisClient *redis_client.RedisClient, mongoClient *mongo_client.MongoClient) {
	metrics, err := redisClient.GetAllFromList(SystemMetricsQueueRedis)
	if err != nil {
		fmt.Println("Failed getting metrics from redis queue: " + err.Error())
		return
	}

	fmt.Println("Got metrics from redis queue")
	err = redisClient.RemoveAllFromList(SystemMetricsQueueRedis)
	if err != nil {
		fmt.Println("Failed removing metrics from redis queue: " + err.Error())
		return
	}

	fmt.Println("Deleted all metrics from redis queue")

	metricsAggregated := aggregateMetricsByName(metrics)
	uploadAggregatedMetricsToMongo(mongoClient, metricsAggregated)

	err = redisClient.RemoveAllFromList("system-metrics")
	if err != nil {
		println("Failed deleting all elements from Redis")
	}

}

func aggregateMetricsByName(metrics []string) map[string]int {
	metricsCount := make(map[string]int)
	for _, metric := range metrics {
		count, ok := metricsCount[metric]
		if ok {
			metricsCount[metric] = count + 1
		} else {
			metricsCount[metric] = 1
		}
	}

	return metricsCount
}

func uploadAggregatedMetricsToMongo(mongoClient *mongo_client.MongoClient, metricsAgg map[string]int) {
	for metricName, metricCount := range metricsAgg {
		if metricCount == 0 {
			fmt.Println("Skipping " + metricName + " metric because it has a value of 0")
			continue
		}

		document := aggregatedMetricDocument{
			MetricName: metricName,
			Count:      metricCount,
			Timestamp:  time.Now(),
		}

		metricJsonData, err := json.Marshal(document)
		if err != nil {
			println("Failed serializing metric " + metricName + "!")
			continue
		}

		metricJsonString := string(metricJsonData)
		err = mongoClient.InsertJSONDocument(&metricJsonString, "system-metrics")
		if err != nil {
			println("Failed saving json document " + metricJsonString + " for metric " + metricName + " into mongo")
		}
	}
}
