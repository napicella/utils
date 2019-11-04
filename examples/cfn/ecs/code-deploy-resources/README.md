Code is taken from the public [git hub repo](https://github.com/BalmanRawat/bgECSService).
It creates custom resources to build an ECS service with Code Deploy controller with CloudFormation. 

### Blue Green ECS Service with custom CloudFormation
#### Why custom ECS Service

AWS ECS Service supports three types of Deployment Controller([For More Details](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeploymentController.html)):
- ECS (default)
- CODE_DEPLOY(blue/green)
- EXTERNAL

We can set CODE_DEPLOY as deployment controller and achieve BLUE/GREEN deployment of the ECS Service unfortunately 
there is no support to manipulate DeploymentController properties through Cloudformation. You have to either use API/CLI 
or build your own custom Cloudformation resource.
                                                                                      

#### Why custom deployment group
AWS CloudFormation does not have native support for the ECS type Blue/Green DeploymentGroup.
use API/CLI or build your own custom cloudformation.


### One time onboarding of a new account
This process creates the customer resources used by the CFN templates for our integ tests.
It needs to be performed once for every new account in which we want to run integ tests 
(including the dev account).

- set env variables for Code deploy stack
```
export AWS_ACCOUNT=<account here>
export SAM_SOURCE_BUCKET="do-not-delete-cloud-debug-custom-code-deploy-${AWS_ACCOUNT}"
export STACK_NAME="CodeDeploy-CustomResource"
```

- create the bucket if does not exist
```
aws s3 mb s3://"${SAM_SOURCE_BUCKET}"
```

- build and deploy the Code Deploy custom resource
```
cd create-ecs-service
make package
make deploy
```

- set env variables for Deployment group stack
```
export SAM_SOURCE_BUCKET="do-not-delete-cloud-debug-custom-deployment-group-${AWS_ACCOUNT}"
export STACK_NAME="DeploymentGroup-CustomResource"
```

- create the bucket if does not exist
```
aws s3 mb s3://"${SAM_SOURCE_BUCKET}"
```

- build and deploy the Deployment Group custom resource
```
cd create-deployment-group
make package
make deploy
```

