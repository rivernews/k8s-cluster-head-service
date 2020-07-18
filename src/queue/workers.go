package queue

import (
	"errors"
	"log"
	"strconv"
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

	job.Checkin("Checking in..!")

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
	var checkInMessage string
	checkInMessage = "Starting guide job..."
	job.Checkin(checkInMessage)

	// request k8s provision

	simulatedK8sProvisionSlackRequest := types.SlackRequestType{
		Token:       utilities.RequestFromSlackTokenCredential,
		TriggerWord: "kkk",
		Text:        "kkk:m",
	}
	k8sProvisioningRequestedPipeline, triggerK8sError := utilities.CircleCITriggerK8sClusterHelper(simulatedK8sProvisionSlackRequest)
	if triggerK8sError != nil {
		utilities.Logger("ERROR", "Failed to trigger K8s provisioning: ", triggerK8sError.Error(), "... job aborted")
		return triggerK8sError
	}
	job.Checkin("K8s provision request done successfully")
	// poll till k8s finish
	k8sProvisioningFinalStatus, waitK8sProvisioningError := utilities.CircleCIWaitTillWorkflowFinish(k8sProvisioningRequestedPipeline.ID)
	// error handling
	if waitK8sProvisioningError != nil {
		utilities.Logger("ERROR", "Failed to poll k8s provisioning: ", waitK8sProvisioningError.Error(), "... job aborted")
		return waitK8sProvisioningError
	}
	if k8sProvisioningFinalStatus != "success" {
		return utilities.Logger("ERROR", "K8s provision completed but wasn't successful, final status: ", k8sProvisioningFinalStatus)
	}

	checkInMessage = utilities.BuildString("K8s provisioning finished: ", k8sProvisioningFinalStatus)
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)

	// deploy SLK

	// request SLK deployment provision
	simulatedSLKDeploymentSlackRequest := types.SlackRequestType{
		Token:       utilities.RequestFromSlackTokenCredential,
		TriggerWord: "slk",
		Text:        "slk",
	}
	travisCIRequestProvision, requestSLKDeployError := utilities.TravisCITriggerSLKHelper(simulatedSLKDeploymentSlackRequest)
	checkInMessage = "Requested build for SLK deployment"
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)
	if requestSLKDeployError != nil {
		utilities.Logger("ERROR", "Failed to request SLK deployment: ", requestSLKDeployError.Error())
		return requestSLKDeployError
	}
	if travisCIRequestProvision.Request.ID < 0 {
		utilities.Logger("ERROR", "SLK deployment request ID is invalid: ", strconv.Itoa(travisCIRequestProvision.Request.ID))
		return errors.New("SLK deployment request ID is invalid")
	}
	slkDeploymentRequestID := travisCIRequestProvision.Request.ID
	checkInMessage = utilities.BuildString("Requested SLK build successfully, request ID: ", strconv.Itoa(slkDeploymentRequestID))
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)
	// polling till build provisioned
	slkDeploymentRequestIDAsString := strconv.Itoa(slkDeploymentRequestID)
	slkDeploymentBuild, slkDeploymentBuildProvisionError := utilities.TravisCIWaitUntilBuildProvisioned(slkDeploymentRequestIDAsString)
	if slkDeploymentBuildProvisionError != nil {
		utilities.Logger("ERROR", "SLK deployment build failed to provision: ", slkDeploymentBuildProvisionError.Error())
		return slkDeploymentBuildProvisionError
	}
	if slkDeploymentBuild.ID < 0 || slkDeploymentBuild.State == "" {
		return utilities.Logger("ERROR", "SLK deployment build data is empty")
	}
	slkDeploymentBuildIDAsString := strconv.Itoa(slkDeploymentBuild.ID)
	checkInMessage = utilities.BuildString("SLK deployment build provisioned by ID: ", slkDeploymentBuildIDAsString)
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)
	// polling till SLK finish
	slkDeploymentFinalStatus, slkDeploymentError := utilities.TravisCIWaitTillBuildFinish(slkDeploymentBuildIDAsString)
	// error handing
	if slkDeploymentError != nil {
		utilities.Logger("ERROR", "Build failed to deploy SLK: ", slkDeploymentError.Error())
		return slkDeploymentError
	}
	if slkDeploymentFinalStatus != "passed" {
		utilities.Logger("ERROR", "Build for SLK finalized but wasn't successful: ", slkDeploymentFinalStatus)
		return errors.New("Build for SLK finalized but wasn't successful")
	}
	checkInMessage = utilities.BuildString("SLK deploymeny finished: ", slkDeploymentFinalStatus)
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)

	// request s3 job

	// polling till s3 finish
	s3JobMeta, s3JobWaitError := utilities.SLKWaitTillS3JobFinish()
	// error handling
	if s3JobWaitError != nil {
		return utilities.Logger("ERROR", "S3 job failed: ", s3JobWaitError.Error())
	}
	checkInMessage = utilities.BuildString("S3 job finalized with status: ", s3JobMeta.Status)
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)

	// scale dowm k8s cluster

	simulatedK8sDestroySlackRequest := types.SlackRequestType{
		Token:       utilities.RequestFromSlackTokenCredential,
		TriggerWord: "ddd",
		Text:        "ddd",
	}
	// polling wait k8s cluster destroy
	k8sDestroyRequestedPipeline, triggerK8sDestroyError := utilities.CircleCITriggerK8sClusterHelper(simulatedK8sDestroySlackRequest)
	if triggerK8sDestroyError != nil {
		utilities.Logger("ERROR", "Failed to trigger K8s destroy: ", triggerK8sDestroyError.Error(), "... job aborted")
		return triggerK8sDestroyError
	}
	// poll till destroy finish
	k8sDestroyFinalStatus, waitK8sDestroyError := utilities.CircleCIWaitTillWorkflowFinish(k8sDestroyRequestedPipeline.ID)
	// error handling
	if waitK8sDestroyError != nil {
		utilities.Logger("ERROR", "Failed to poll k8s destroy: ", waitK8sDestroyError.Error(), "... job aborted")
		return waitK8sDestroyError
	}
	if k8sDestroyFinalStatus != "success" {
		return utilities.Logger("ERROR", "K8s destroy completed but wasn't successful, final status: ", k8sDestroyFinalStatus)
	}
	checkInMessage = utilities.BuildString("K8s destroy finished: ", k8sDestroyFinalStatus)
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)

	// mark guide job as success

	checkInMessage = "🎉 K8s Guided Job completed successfully."
	utilities.Logger("INFO", checkInMessage)
	job.Checkin(checkInMessage)
	utilities.SendSlackMessage(checkInMessage)

	return nil
}
