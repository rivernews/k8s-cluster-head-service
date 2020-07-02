package queue

import (
	"errors"
	"log"
	"time"

	"github.com/gocraft/work"

	"github.com/rivernews/k8s-cluster-head-service/v2/src/types"
	"github.com/rivernews/k8s-cluster-head-service/v2/src/utilities"
)

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	log.Println("Processing job " + job.ID + "...")

	time.Sleep(1 * time.Second)

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)
	log.Println("address is " + addr)
	log.Println("subject is " + subject)

	return errors.New("This is a testing failure job")
}

func (c *Context) Export(job *work.Job) error {
	return nil
}

func (c *Context) GuidedSLKS3JobElasticScalingSession(job *work.Job) error {
	utilities.Logger("INFO", "Starting GuidedSLKS3JobElasticScalingSession...")

	// request k8s provision
	simulatedSlackRequest := types.SlackRequestType{
		Token:       utilities.RequestFromSlackTokenCredential,
		TriggerWord: "kkk",
		Text:        "kkk:m",
	}
	k8sProvisioningRequestedPipeline, triggerK8sError := utilities.CircleCITriggerK8sClusterHelper(simulatedSlackRequest)
	if triggerK8sError != nil {
		utilities.Logger("ERROR", "Failed to trigger K8s provisioning: ", triggerK8sError.Error(), "... job aborted")
		return triggerK8sError
	}
	// poll till k8s finish
	k8sProvisioningFinalStatus, waitK8sProvisioningError := utilities.CircleCIWaitTillWorkflowFinish(k8sProvisioningRequestedPipeline.ID)
	// error handling
	if waitK8sProvisioningError != nil {
		utilities.Logger("ERROR", "Failed to poll k8s provisioning: ", waitK8sProvisioningError.Error(), "... job aborted")
		return waitK8sProvisioningError
	}
	if k8sProvisioningFinalStatus != "success" {
		utilities.Logger("ERROR", "K8s provision completed but wasn't successful, final status: ", k8sProvisioningFinalStatus)
	}

	log.Print("K8s provisioning finished: " + k8sProvisioningFinalStatus)

	// TODO: request SLK deployment provision
	// slkDeploymentRequestID := "123"
	// TODO: polling till build provisioned
	// TODO: polling till SLK finish
	// slkDeploymentFinalStatus, slkDeploymentError := utilities.TravisCIWaitUntilBuildProvisioned(slkDeploymentRequestID)
	// TODO: error handing
	// if slkDeploymentError != nil {
	// 	return slkDeploymentError
	// }
	// log.Print("SLK deploymeny finished: " + slkDeploymentFinalStatus)

	// TODO: request s3 job
	// TODO: polling till s3 finish
	// TODO: error handling

	// TODO: scale dowm k8s cluster
	// TODO: polling
	// TODO: error handling

	utilities.Logger("DEBUG", "Job finished w/o error")

	return nil
}
