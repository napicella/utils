### API Gateway and Lambda
- sam logs --name Lambda-name --tail
  aws lambda list-functions --region us-east-1 --query 'Functions[].FunctionName' --output text
- to remove the log stream from the sam logs  
  sam logs --name Lambda-name --tail | awk {'first = $1; $1=""; print $0'}
- https://stackoverflow.com/questions/50331588/aws-api-gateway-custom-authorizer-strange-showing-error
- Use swagger for everything in the CFN template instead of the Events property of AWS::Serverless::Function, much better documented and friendly!
- Explicitely add permission to API gateway to call your lambda (this should ve been done by SAM behind the scene for you)
  Current workaround needed for error: "Execution failed due to configuration error: Invalid permissions on Lambda function". 
  If you get this error (or a similar permission denied message) even when you try to test the API Gateway from the AWS     console than odds are that you are impacted by a bug in SAM that does not appear to have been fixed (note that it also     happens with CDK, which might imply that the problem is in Cloudformation or API Gateway). See https://github.com/awslabs/serverless-application-model/issues/59  
  As mentioned, the solution is to explicitely define permissions to invoke the Lambda Function.
  ```
  ConfigLambdaPermission:
    Type: "AWS::Lambda::Permission"
    DependsOn:
      - CustomAuthorizerFunction
    Properties:
      Action: lambda:InvokeFunction
      FunctionName: !Ref CustomAuthorizerFunction
      Principal: apigateway.amazonaws.com
  ``` 
  This is the good way to do it. Another option is to open the API Gateway in the AWS console, click on Resources, select   the endpoint, then click on Integration Request. Remove the role and click the save button. Then put back the role arn and   click save. Finally deploy API Gateway. As strange as it sounds, this actually solves the issue.
- API gateway timeout hard limit is 30 seconds. A lambda with timeout grater than 30 seconds does
  not make sense in combination with API Gateway

### ECS/Fargate with SSM 
```
aws --region eu-central-1 ssm start-session --target <cluster Name>_<task ID>_<container runtime ID>
```

if containers log contain:
2020-03-10 10:15:03 ERROR [MessageGatewayService] Failed to get controlchannel token, error: CreateControlChannel failed with error: createControlChannel request failed: unexpected response from the service <BadRequest xmlns=""><message>Unauthorized request.</message></BadRequest>

Then u likely forgot to give ssm permissions to the Task role

### FARGATE with CloudFormation
After the cluster has been created, you need to associate  a cluster capacity provider.
You need to use the AWS Cli, because CloudFormation does not support this yet:
```bash
aws ecs put-cluster-capacity-providers --cluster WebServiceExampleCluster --capacity-providers FARGATE --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1
```

### XRAYS

__Takeaway 1__
XRay allows to build a latency map for a request.
The number in the circle is the time between the moment the request entered the node and the one in which it exited.
By clicking on the edge, you get the time for the request to reach the next node. Notice that this measure is only accurate if the next node is instrumented with XRay, otherwise it's an estimate computed on the caller side.

![AWS XRay service map](https://dev-to-uploads.s3.amazonaws.com/i/8oe7d3lkju833shzuoqt.png)

__Takeaway 2__
Dropping in XRay in an existing application is easy, just wrap the client.
	
Instrumenting http calls made via an http client:
```golang
client := xray.Client(&http.Client{})
req, _ := http.NewRequest(http.MethodGet, "...", nil)
res, _ := client.Do(req.WithContext(r.Context()))
```

Instrumenting calls to AWS Services:
```golang
// DynamoDB example
dynamoDbSvc := dynamodb.New(sess)
xray.AWS(dynamoDbSvc.Client)
```

__Takeaway 3__
Using XRay in a Lambda with Golang requires a XRay libray version greater or equal to v1. This is important, otherwise you will be presented with a cryptic error message at runtime!
```
go get –u github.com/aws/aws-xray-sdk-go@v1.0.0-rc.14
```

__Takeaway 4__
You need to pass the context to each call you make. In the following snippet, notice the _PutItemWithContext_ call instead of _PutItem_:
```golang
func handleRequest(
   ctx context.Context, request events.ALBTargetGroupRequest) 
  (events.ALBTargetGroupResponse, error) {
...
_, e = dynamoDbSvc.PutItemWithContext(ctx, &dynamodb.PutItemInput{ 
   Item:  av,  
   TableName: aws.String(tableName),
   ConditionExpression: aws.String("attribute_not_exists(wid)"),
})
```
If you are running in a Lambda, the context is what you receive from the Lambda invocation. Otherwise you need to create a new Context.

If you get an error like:
```
failed to begin subsegment
named 'dynamo': segment cannot be found.
```
Then you likely forgot to pass a context to the request.

__Takeaway 5__
Many services integrate with XRay, for example API Gateway.
If you enable XRay for API Gateway via console or CloudFormation you need to manually trigger a deployment. Otherwise you'll get UNAUTHORIZED from API Gateway!

