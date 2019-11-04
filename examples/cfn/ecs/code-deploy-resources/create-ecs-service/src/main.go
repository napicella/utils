package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func waitServiceSteadyState(input *ecs.DescribeServicesInput) error {

	pL := fmt.Println
	pF := fmt.Printf

	pF("Waiting for %s Services in %s Cluster to reach steady state\n", *input.Services[0], *input.Cluster)
	svc := ecs.New(session.New())

	timeOut := time.Now().Add(time.Minute * 10)
	nowIs := time.Now

	//Wait until timeOut
	for nowIs().Before(timeOut) {
		//Waiting before making DescribeRequest
		time.Sleep(time.Second * 15)

		result, err := svc.DescribeServices(input)
		if err != nil {
			pL(err)
			return err
		}

		//Just Printing the events of the service for debugging
		if len(result.Services[0].Events) > 0 {
			msg := *result.Services[0].Events[0].Message
			pF("Latest Event Message: %s\n", msg)
		}

		//result.Services[0].TaskSets[0].StabilityStatus has two value STEADY_STATE and STABILIZING
		if len(result.Services[0].TaskSets) > 0 {
			status := *result.Services[0].TaskSets[0].StabilityStatus
			if status == "STEADY_STATE" {
				pF("Reached Stability State")
				return nil
			}
		}

	}

	return errors.New("Wait Timeout for the ECS service..")
}

func createService(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println
	pL("Inside Create Service Method...")

	svc := ecs.New(session.New())
	pL(svc)

	//Preparing Input for the CreateService API
	var input ecs.CreateServiceInput
	props := event.ResourceProperties
	if jData, error := json.Marshal(props); error != nil {
		err = error
		return
	} else {
		json.Unmarshal(jData, &input)
	}

	//While Unmarshaling string format numeric data is converted to zero. eg. DesiredCount: "1"(is integer value in string form), DesiredCount: 0 while Unmarshaling so, need to do extra job here, for DesiredCount, MaximumPercent, MinimumHealthyPercent, Port

	//DesiredCount
	if dCount, ok := props["DesiredCount"].(string); ok {
		//do something here
		intDCount, error := strconv.ParseInt(dCount, 0, 64)
		if error != nil {
			pL(error)
			err = errors.New("DesiredCount is invalid, Please refer docs")
			return
		}

		input.DesiredCount = &intDCount
	}

	//Working with MaximumPercent, MinimumHealthyPercent
	if dConf, ok := props["DeploymentConfiguration"].(map[string]interface{}); ok {

		var dcon ecs.DeploymentConfiguration
		if maxPer, mOk := dConf["MaximumPercent"]; mOk {
			intMaxPer, maxErr := strconv.ParseInt(maxPer.(string), 0, 64)
			if maxErr != nil {
				pL(maxErr)
				err = errors.New("MaximumPercent value is invalid, Please refer docs")
				return
			}
			dcon.MaximumPercent = &intMaxPer
		}

		if minHeal, mOk := dConf["MinimumHealthyPercent"]; mOk {
			intMinHeal, minErr := strconv.ParseInt(minHeal.(string), 0, 64)
			if minErr != nil {
				pL(minErr)
				err = errors.New("MinimumHealthyPercent value is invalid, Please refer docs")
				return
			}
			dcon.MinimumHealthyPercent = &intMinHeal
		}
		input.DeploymentConfiguration = &dcon
	}

	//Setting port for loadbalancers/ServiceRegistries
	if lBs, ok := props["LoadBalancers"].([]interface{}); ok {
		lB0 := lBs[0].(map[string]interface{})
		if cPort, cOk := lB0["ContainerPort"]; cOk {
			if intCPort, error := strconv.ParseInt(cPort.(string), 0, 64); error != nil {
				err = error
				return
			} else {

				for _, lb := range input.LoadBalancers {
					lb.ContainerPort = &intCPort
				}

				for _, sr := range input.ServiceRegistries {
					sr.Port = &intCPort
				}
			}
		}
	}

	pL("Input data for the createService API: ", input)

	//Creating Serivce
	result, err := svc.CreateService(&input)
	if err != nil {
		return physicalResourceID, data, err
	}
	pL("Create Service API Response:", result)

	//Wait for the service to become stable
	dSI := &ecs.DescribeServicesInput{
		Cluster:  input.Cluster,
		Services: []*string{input.ServiceName},
	}
	if error := waitServiceSteadyState(dSI); error != nil {
		err = error
	}

	//Setting the physicalResourceID with the ARN of the ECS service
	physicalResourceID = *result.Service.ServiceArn

	//Setting the Name attribute, for the GetAtt cloudformation function
	data = make(map[string]interface{})
	data["Name"] = *input.ServiceName
	pL(data)

	pL("Successfully Created the BlueGreen ECS service")

	return
}

func deleteService(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println
	pF := fmt.Printf

	pL("Inside Delete Service Method...")

	//Preparing Input for the CreateService API
	props := event.ResourceProperties
	var input ecs.DeleteServiceInput

	force := true
	input.Force = &force

	if c, ok := props["Cluster"].(string); ok {
		input.Cluster = &c
	}

	if s, ok := props["ServiceName"].(string); ok {
		input.Service = &s
	}

	//Creating Session
	svc := ecs.New(session.New())
	result, error := svc.DeleteService(&input)
	if error != nil {
		pL("Error Delete API Response: ", error)
		err = error
		return
	}
	pL("Delete Service API Response: ", result)

	//Wait Delete Service
	pF("Waiting for %s Service to be deleted in %s Cluster\n", *input.Cluster, *input.Cluster)

	dSI := &ecs.DescribeServicesInput{
		Cluster:  input.Cluster,
		Services: []*string{input.Service},
	}
	if error := svc.WaitUntilServicesInactive(dSI); error != nil {
		pL("Error Delete API Response: ", error)
		err = error
		return
	}

	pL("Successfully Delete...")

	return
}

/*
For services using the blue/green (CODE_DEPLOY) deployment controller, only the
- desired count,
- deployment configuration, and
- health check grace period can be updated.
Rest Other Parameter updates are ignored
TO Update
- network configuration
- platform version
- task definition
a new AWS CodeDeploy deployment should be created. For more information, see https://docs.aws.amazon.com/sdk-for-go/api/service/ecs/#ECS.UpdateService
*/
func updateService(event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	pL := fmt.Println

	pL("Inside update service..")

	// Setting the physicalresourceId
	if event.PhysicalResourceID != "" {
		physicalResourceID = event.PhysicalResourceID
	} else {
		err = errors.New("PhysicalResource ID not found")
		return
	}

	//Check if the update needs wait. Only desiredCount update needs wait
	var needsWait bool
	var input ecs.UpdateServiceInput

	//Preparing Input for the CreateService API
	props := event.ResourceProperties
	oldProps := event.OldResourceProperties

	//Gets the New and Old DesiredCount; Compares old/new DesiredCount if changed sets needsWait=true
	if dCount, ok := props["DesiredCount"].(string); ok {
		//do something here
		intDCount, error := strconv.ParseInt(dCount, 0, 64)
		if error != nil {
			pL(error)
			err = errors.New("DesiredCount is invalid, Please refer docs")
			return
		}

		// There is only need to wait while updating the service, when the desired count has changed
		if oDCount, ok := oldProps["DesiredCount"].(string); ok {
			if intODCount, error := strconv.ParseInt(oDCount, 0, 64); error == nil {
				if intODCount != intDCount {
					needsWait = true
				}
			}
		}

		input.DesiredCount = &intDCount

	}

	//Getting with MaximumPercent, MinimumHealthyPercent
	if dConf, ok := props["DeploymentConfiguration"].(map[string]interface{}); ok {

		var dcon ecs.DeploymentConfiguration
		if maxPer, mOk := dConf["MaximumPercent"].(string); mOk {
			intMaxPer, maxErr := strconv.ParseInt(maxPer, 0, 64)
			if maxErr != nil {
				pL(maxErr)
				err = errors.New("MaximumPercent value is invalid, Please refer docs")
				return
			}
			dcon.MaximumPercent = &intMaxPer
		}

		if minHeal, mOk := dConf["MinimumHealthyPercent"].(string); mOk {
			intMinHeal, minErr := strconv.ParseInt(minHeal, 0, 64)
			if minErr != nil {
				pL(minErr)
				err = errors.New("MinimumHealthyPercent value is invalid, Please refer docs")
				return
			}
			dcon.MinimumHealthyPercent = &intMinHeal
		}
		input.DeploymentConfiguration = &dcon
	}

	if s, ok := props["ServiceName"].(string); ok {
		input.Service = &s
	}
	if c, ok := props["Cluster"].(string); ok {
		input.Cluster = &c
	}

	// Create session and update the service
	svc := ecs.New(session.New())
	result, uErr := svc.UpdateService(&input)
	pL("Input Data for update: ", input)
	if uErr != nil {
		pL(uErr)
		err = uErr
		return
	}

	pL("Success API Response from Update Service", result)

	if needsWait {
		//Wait Section for the service to become stable
		dSI := &ecs.DescribeServicesInput{
			Cluster:  input.Cluster,
			Services: []*string{input.Service},
		}
		if error := waitServiceSteadyState(dSI); error != nil {
			err = error
			return
		}
	}

	//Setting the Name attribute, for the GetAtt cloudformation function
	data = make(map[string]interface{})
	data["Name"] = *input.Service
	pL(data)

	pL("Successfully Updated the Service")

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
		return createService(event)

	case "Update":
		if physicalResourceID == "" {
			// Invalid Request, Update must have a valid Physical ID of the resource to be deleted
			err = errors.New("Invalid Request, Update must have a valid Physical ID of the resource to be deleted")
			return
		}

		return updateService(event)

	case "Delete":
		if physicalResourceID == "" {
			// Invalid Request, Update must have a valid Physical ID of the resource to be deleted
			err = errors.New("Invalid Request, Update must have a valid Physical ID of the resource to be deleted")
			return
		}

		return deleteService(event)

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
