package queue

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"
)

// https://github.com/gocraft/work
func TestJobQueue() {
	defer utilities.RedisPool.Close()

	// flushdb

	conn := utilities.RedisPool.Get()
	defer conn.Close()
	reply, flushdbErr := conn.Do("FLUSHDB")
	if flushdbErr == nil {
		log.Printf("Successfully flushed! %s", reply)
	} else {
		log.Fatalf("Flush failed: %s", flushdbErr)
	}

	// client connection stats
	// redis doc:
	// https://redis.io/commands/client-list

	statsReply, statesReplyErr := conn.Do("CLIENT", "LIST")
	if statesReplyErr != nil {
		log.Panicf("Failed to read redis client stats: %s", statesReplyErr)
	} else {
		log.Printf("Redis client stats:\n%s\n", statsReply)
	}
	statValues, statValuesErr := redis.String(statsReply, statesReplyErr)
	if statValuesErr != nil {
		log.Fatalf("Failed to parse redis client stat value: %s", statValuesErr)
	}
	clientStatLines := strings.Split(statValues, "\n")
	log.Printf("%d redis client connections", len(clientStatLines)-1)

	// health check

	if !(flushdbErr == nil && statesReplyErr == nil) {
		return
	}

	// enqueue

	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer("my_app_namespace", utilities.RedisPool)

	// worker pool

	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, "my_app_namespace", utilities.RedisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)

	// Map the name of jobs to handler functions
	pool.JobWithOptions("send_email", work.JobOptions{
		MaxFails: 1,
	}, (*Context).SendEmail)

	// Customize options:
	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)
	pool.JobWithOptions("guided_k8s_s3_elastic_session", work.JobOptions{
		MaxFails: 1,
	}, (*Context).GuidedSLKS3JobElasticScalingSession)

	// Start processing jobs
	pool.Start()

	// enqueue jobs
	// Enqueue a job named "send_email" with the specified parameters.
	// for i := 1; i < 5; i++ {
	// 	enqueuer.Enqueue("send_email", work.Q{"address": "test@example.com", "subject": "hello world", "customer_id": 4})
	// }
	enqueuer.Enqueue("guided_k8s_s3_elastic_session", work.Q{})

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	utilities.Logger("INFO", "Worker pool stopped")
	pool.Stop()
}
