Title: Managing AWS lambda functions from start to finish with Terraform
Date: 2018-08-02 14:00
Category: infrastructure
Status: draft

[AWS lambda](https://aws.amazon.com/lambda/) functions look deceptively simple. The devil is in the details though. Once you
have written the code and have created a `.zip` file, there's a few more steps to go.

For starters, we need an IAM profile to be defined with appropriate policies allowing the function to access the AWS resources. 
To setup the lambda function to be invoked automatically in reaction to another event, we need some more permissions. Then, we 
have to create a lambda function in AWS infrastructure and point it to our `.zip` file that we have created above. Everytime, we 
update this, `.zip`, we have to ask AWS lambda to update the code again. 

One straight forward, no fuss approach is to use the [AWS CLI](https://docs.aws.amazon.com/cli/latest/reference/lambda/index.html).
The main problem I think with this approach and using any of the serverless
tools and frameworks out there like `apex`, `serverless` or `zappa` out there is that they treat the infrastructure of 
your lambda functions as islands, rather than being part of your broader AWS infrastructure. The same S3 bucket's contents 
which you want your lambda function to be triggered in reaction to changes in may be the bucket some other non-lambda
application writes to. You want to run your lambda function in the same VPC as your database RDS instance. Needless to say,
there will cross-application infrastructure references. 

What follows is a non-production tested suggestion for managing your lambda functions and their infrastructure as part of
your global infrastructure as code repository.

## Example: Managing lambda functions using Terraform

Consider a lambda function [ec2_state_change](https://github.com/amitsaha/cloudwatch-event-lambda/tree/master/functions/ec2_state_change).
I wrote this for a recent [article](https://blog.codeship.com/cloudwatch-event-notifications-using-aws-lambda/). The `src` directory
has the source of the lambda function which is written in Python. To create the lambda function (for the first time) and to
deploy new versions of the code, the following BASH script is run:

```
!/usr/bin/env bash
set -ex

# Create a .zip of src
pushd src
zip -r ../src.zip *
popd

aws s3 cp src.zip s3://aws-health-notif-demo-lambda-artifacts/ec2-state-change/src.zip
version=$(aws s3api head-object --bucket aws-health-notif-demo-lambda-artifacts --key ec2-state-change/src.zip)
version=$(echo $version | python -c 'import json,sys; obj=json.load(sys.stdin); print(obj["VersionId"])')

# Deploy to demo environment
pushd ../../terraform/environments/demo
terraform init
terraform apply \
    -var aws_region=ap-southeast-2 \
    -var ec2_state_change_handler_version=$version \
    -target=module.ec2_state_change_handler.module.ec2_state_change_handler.aws_lambda_function.lambda \
    -target=module.ec2_state_change_handler.module.ec2_state_change_handler.aws_cloudwatch_event_rule.rule \
    -target=module.ec2_state_change_handler.module.ec2_state_change_handler.aws_cloudwatch_event_target.target \
    -target=module.ec2_state_change_handler.module.ec2_state_change_handler.aws_iam_role_policy.lambda_cloudwatch_logging \
    -target=module.ec2_state_change_handler.module.ec2_state_change_handler.aws_lambda_permission.cloudwatch_lambda_execution
popd

```

The above script does the following main things:

- First, it creates a ZIP of the `src` sub-directory
- It uploads it to a designated bucket in S3
- It gets the version of the file it just uploaded
- It then changes the current working directory to one with `terraform` configuration in it
- Then runs `terraform apply` supplying specific targets related to the function it is uploading


The first time, this script is run, it will create all the infrastructure that is needed by the lambda function
to be run. On subsequent applications, only the lambda function's version will change.

That's it. For every terraform environment, we can have a script which will do the same. We can run this script as part of a CI/CD
pipeline. The repository pointed to above has the terraform configuration in it as well, but we can always download the terrraform
configuration tarball or git clone it during a CI run. The key idea here is that, your terraform configuration for the lambda function
can co-exist with the rest of your infrastructure.


## 
