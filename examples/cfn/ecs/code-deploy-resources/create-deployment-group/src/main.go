package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
)

func deleteDeploymentGroup(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println

	pL("Inside Create")
	svc := codedeploy.New(session.New())
	pL(svc)

	var input codedeploy.DeleteDeploymentGroupInput
	props := event.ResourceProperties

	if a, ok := props["ApplicationName"].(string); ok {
		input.ApplicationName = &a

	} else {
		err = errors.New("CodeDeploy Application name should be provided")
		return
	}

	if d, ok := props["DeploymentGroupName"].(string); ok {
		input.DeploymentGroupName = &d
	} else {
		err = errors.New("CodeDeploy DeploymentGroupName name should be provided")
		return
	}

	//Deleting DeploymentGroup
	result, error := svc.DeleteDeploymentGroup(&input)
	if error != nil {
		err = error
		return
	}

	pL("Delete DeploymentGroup API Response:", result)

	return
}

//It is needed inorder to Unmarshal string type integer values to actual integer values
type BlueGreenDeploymentConfiguration struct {
	DeploymentReadyOption struct {
		ActionOnTimeout   string
		WaitTimeInMinutes int64 `json:",string"`
	}
	TerminateBlueInstancesOnDeploymentSuccess struct {
		Action                       string
		TerminationWaitTimeInMinutes int64 `json:",string"`
	}
}

func createDeploymentGroup(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println

	pL("Inside CreateDeploymentGroup")
	svc := codedeploy.New(session.New())
	pL(svc)

	//Preparing Input for the CreateService API
	var input codedeploy.CreateDeploymentGroupInput
	var myInput BlueGreenDeploymentConfiguration
	props := event.ResourceProperties
	if jData, error := json.Marshal(props); error != nil {
		err = error
		return
	} else {
		json.Unmarshal(jData, &input)
	}

	if _, ok := props["BlueGreenDeploymentConfiguration"]; ok {
		if jData, error := json.Marshal(props["BlueGreenDeploymentConfiguration"]); error != nil {
			err = error
			return
		} else {
			json.Unmarshal(jData, &myInput)

			//Assigning the WaitTimeInMinutes and TerminationWaitTimeMinutes from custom input to codeploy input
			input.BlueGreenDeploymentConfiguration.DeploymentReadyOption.WaitTimeInMinutes = &myInput.DeploymentReadyOption.WaitTimeInMinutes
			input.BlueGreenDeploymentConfiguration.TerminateBlueInstancesOnDeploymentSuccess.TerminationWaitTimeInMinutes = &myInput.TerminateBlueInstancesOnDeploymentSuccess.TerminationWaitTimeInMinutes
		}
	}

	//Settting physicalResourceID as the deploymentgroup name
	if dN, ok := props["DeploymentGroupName"].(string); ok {
		physicalResourceID = dN
	} else {
		err = errors.New("DeploymentGroupName not provided!!")
		return
	}

	//Creating DeploymentGroup
	result, error := svc.CreateDeploymentGroup(&input)
	if error != nil {
		err = error
		return
	}

	pL("Create DeploymentGroup API Response:", result)

	return
}

func updateDeploymentGroup(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println

	pL("Inside UpdateDeploymentGroup")

	// Setting the physicalresourceId
	if event.PhysicalResourceID != "" {
		physicalResourceID = event.PhysicalResourceID
	} else {
		err = errors.New("PhysicalResource ID not found")
		return
	}

	svc := codedeploy.New(session.New())
	pL(svc)

	//Preparing Input for the CreateService API
	var input codedeploy.UpdateDeploymentGroupInput
	var myInput BlueGreenDeploymentConfiguration
	props := event.ResourceProperties

	if jData, error := json.Marshal(props); error != nil {
		err = error
		return
	} else {
		json.Unmarshal(jData, &input)
	}

	if _, ok := props["BlueGreenDeploymentConfiguration"]; ok {
		if jData, error := json.Marshal(props["BlueGreenDeploymentConfiguration"]); error != nil {
			err = error
			return
		} else {
			json.Unmarshal(jData, &myInput)

			//Assigning the WaitTimeInMinutes and TerminationWaitTimeMinutes from custom input to codeploy input
			input.BlueGreenDeploymentConfiguration.DeploymentReadyOption.WaitTimeInMinutes = &myInput.DeploymentReadyOption.WaitTimeInMinutes
			input.BlueGreenDeploymentConfiguration.TerminateBlueInstancesOnDeploymentSuccess.TerminationWaitTimeInMinutes = &myInput.TerminateBlueInstancesOnDeploymentSuccess.TerminationWaitTimeInMinutes
		}
	}

	//DeploymentGroupName is required by the update
	if d, ok := props["DeploymentGroupName"].(string); ok {
		input.CurrentDeploymentGroupName = &d
	}

	//Updating DeploymentGroup
	result, error := svc.UpdateDeploymentGroup(&input)
	if error != nil {
		err = error
		return
	}

	pL("Update DeploymentGroup API Response:", result)

	return
}

func handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println

	// Details about the cloudformation event https://godoc.org/github.com/aws/aws-lambda-go/cfn#Event
	pL("Event Received:", event)

	reqType := string(event.RequestType)
	fmt.Printf("Recieved Request Type: %s\n", reqType)
	resType := string(event.ResourceType)
	fmt.Printf("Resource Type: %s\n", resType)

	physicalResourceID = event.PhysicalResourceID
	switch reqType {

	case "Create":
		pL("Create")
		return createDeploymentGroup(event)

	case "Update":
		if physicalResourceID == "" {
			// Invalid Request, Update must have a valid Physical ID of the resource to be deleted
			err = errors.New("Invalid Request, Update must have a valid Physical ID of the resource to be deleted")
			return
		}

		return updateDeploymentGroup(event)

	case "Delete":
		if physicalResourceID == "" {
			// Invalid Request, Update must have a valid Physical ID of the resource to be deleted
			err = errors.New("Invalid Request, Update must have a valid Physical ID of the resource to be deleted")
			return
		}

		return deleteDeploymentGroup(event)

	default:

		log := "Unkown Request Type, should be one of Create/Update/Delete"
		pL(log)
		return "", nil, errors.New(log)
	}
	return
}

func main() {
	lambda.Start(cfn.LambdaWrap(handler))
}
