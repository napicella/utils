### API Gateway and Lambda
- sam logs --name Lambda-name --tail
  aws lambda list-functions --region us-east-1 --query 'Functions[].FunctionName' --output text
- https://stackoverflow.com/questions/50331588/aws-api-gateway-custom-authorizer-strange-showing-error
- Use swagger for everything in the CFN template instead of the Events property of AWS::Serverless::Function, much better documented and friendly!
- Explicitely add permission to API gateway to call your lambda (this should ve been done by SAM behind the scene for you)
  # Current workaround needed for error:
  # "Execution failed due to configuration error: Invalid permissions on Lambda function."
  # See https://github.com/awslabs/serverless-application-model/issues/59
  # It's possibly a bug in SAM that does not appear to have been fixed
  ConfigLambdaPermission:
    Type: "AWS::Lambda::Permission"
    DependsOn:
      - CustomAuthorizerFunction
    Properties:
      Action: lambda:InvokeFunction
      FunctionName: !Ref CustomAuthorizerFunction
      Principal: apigateway.amazonaws.com
- API gateway timeout hard limit is 30 seconds. A lambda with timeout grater than 30 seconds does
  not make sense in combination with API Gateway

### ECS/Fargate with SSM 
aws --region eu-central-1 ssm start-session --target <cluster Name>_<task ID>_<container runtime ID>

if containers log contain:
2020-03-10 10:15:03 ERROR [MessageGatewayService] Failed to get controlchannel token, error: CreateControlChannel failed with error: createControlChannel request failed: unexpected response from the service <BadRequest xmlns=""><message>Unauthorized request.</message></BadRequest>

Then u likely forgot to give ssm permissions to the Task role

### FARGATE with CloudFormation
After the cluster has been created, you need to associate  a cluster capacity provider.
You need to use the AWS Cli, because CloudFormation does not support this yet:
```bash
aws ecs put-cluster-capacity-providers --cluster WebServiceExampleCluster --capacity-providers FARGATE --default-capacity-provider-strategy capacityProvider=FARGATE,weight=1
```

