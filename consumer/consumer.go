package consumer

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"metrics-consumer/mongo_client"
	"metrics-consumer/redis_client"
	"time"
)

var (
	SystemMetricsCollection = "system-metrics"
)

type queuedMetric struct {
	MetricName string `json:"metric_name"`
}

const SystemMetricsQueueRedis = "system-metrics"

func ConsumeMetrics(redisClient *redis_client.RedisClient, mongoClient *mongo_client.MongoClient) {
	metricsJsons, err := redisClient.GetAllFromList(SystemMetricsQueueRedis)
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

	metricsFromQueue := getMetricNamesFromQueuedJsons(metricsJsons)
	metricsAggregated := aggregateMetricsByName(metricsFromQueue)
	uploadAggregatedMetricsToMongo(mongoClient, metricsAggregated)

	err = redisClient.RemoveAllFromList("system-metrics")
	if err != nil {
		println("Failed deleting all elements from Redis")
	}
}

func getMetricNamesFromQueuedJsons(metricsJsons []string) []string {
	metricsNames := make([]string, 0)
	for _, metricJson := range metricsJsons {
		var metric queuedMetric
		err := json.Unmarshal([]byte(metricJson), &metric)
		if err != nil {
			fmt.Println("Error decoding metric with json " + metricJson + ": " + err.Error())
			continue
		}

		metricsNames = append(metricsNames, metric.MetricName)
	}

	return metricsNames
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

		document := bson.M{
			"metric": metricName,
			"count":  metricCount,
			"timestamp": primitive.Timestamp{
				T: uint32(time.Now().Unix()),
			},
		}

		err := mongoClient.InsertJSONDocument(document, SystemMetricsCollection)
		if err != nil {
			fmt.Println("Failed saving json document for metric " + metricName + " into MongoDB: " + err.Error())
		} else {
			fmt.Printf("Saved document %+v for metric %s in MongoDB", document, metricName)
		}
	}
}
