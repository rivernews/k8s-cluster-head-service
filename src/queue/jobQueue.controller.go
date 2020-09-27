package queue

import (
	"github.com/gocraft/work"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"
)

func HandleJobQueueRequest() {
	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer("my_app_namespace", utilities.RedisPool)

	enqueuer.Enqueue("guided_k8s_s3_elastic_session", work.Q{})
	utilities.SendSlackMessage("`OK`")
}
