package queue

import (
	"errors"
	"log"
	"time"

	"github.com/gocraft/work"

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
	log.Println("Starting GuidedSLKS3JobElasticScalingSession...")

	// TODO: request k8s provision
	k8sProvisioningRequestedPipelineID := "123"
	// TODO: poll till k8s finish
	k8sProvisioningFinalStatus, waitK8sProvisioningError := utilities.CircleCIWaitTillWorkflowFinish(k8sProvisioningRequestedPipelineID)
	// TODO: error handling
	if waitK8sProvisioningError != nil {
		return waitK8sProvisioningError
	}
	log.Print("K8s provisioning finished: " + k8sProvisioningFinalStatus)

	// TODO: request SLK deployment provision
	slkDeploymentRequestID := "123"
	// TODO: polling till SLK finish
	slkDeploymentFinalStatus, slkDeploymentError := utilities.TravisCIWaitUntilBuildProvisioned(slkDeploymentRequestID)
	// TODO: error handing
	if slkDeploymentError != nil {
		return slkDeploymentError
	}
	log.Print("SLK deploymeny finished: " + slkDeploymentFinalStatus)

	// TODO: request s3 job
	// TODO: polling till s3 finish
	// TODO: error handling

	// TODO: scale dowm k8s cluster
	// TODO: polling
	// TODO: error handling

	return nil
}
