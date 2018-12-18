Title: AWS CodeDeploy Deployment Group and Initial Auto Scaling lifecycle hook
Date: 2018-12-19
Category: infrastructure

When we create an AWS Code Deploy [deployment group](https://docs.aws.amazon.com/codedeploy/latest/userguide/deployment-groups.html) via
[Terraform](https://www.terraform.io/) or [CloudFormation](https://aws.amazon.com/cloudformation/) and integrate with an Auto Scaling Group, 
it also by default creates an initial lifecycle hook which ensuresthat a new code deployment gets triggered when a scale-out event occurs. 

It is all very "magical" and it is one of those cases where you have [troublesome](https://github.com/terraform-providers/terraform-provider-aws/issues/2993) behavior especially
when you are managing your infrastructure as code. The troublesome behavior happens as a result of the lifecycle hook creation being
a side-effect of creating a deployment group rather than an explicit operation that the user performs. 

Let's consider this example terraform snippet:

```
resource "aws_codedeploy_app" "code_deploy_app" {
  compute_platform = "Server"
  name             = "${var.service_name}"
}

resource "aws_codedeploy_deployment_group" "deploy_group" {  
  app_name              = "${aws_codedeploy_app.code_deploy_app.name}"
  deployment_group_name = "${var.service_name}-DeploymentGroup${var.environment}"
  service_role_arn      = ".."
  autoscaling_groups = ["${aws_autoscaling_group.autoscaling_group.name}"]
}

resource "aws_launch_configuration" "launch_configuration" {

  lifecycle {
    create_before_destroy = true
  }
  ..
}


resource "aws_autoscaling_group" "autoscaling_group" {

  name_prefix = "${var.service_name}-AutoscalingGroup"
  launch_configuration = "${aws_launch_configuration.launch_configuration.name}"
  ..  
}
```

When we apply the above change to our AWS infrastructure it will:

- Create a code deploy application
- Create a deployment group for this application
- Create an autoscaling group
- Associate the deployment group with the autoscaling group

The above infrastructure changes are all explicit and they map to what we have in code. Let's use the AWS CLI to describe
the deployment group we created above:

```
$ aws deploy get-deployment-group --application-name FinancialReportingService --deployment-group-name FinancialReportingService-DeploymentGroupProduction
{
    "deploymentGroupInfo": {
        "applicationName": "MyService",
        "deploymentGroupId": "b7a6653a-407d-47d8-b9ff-3e0a10b028b3",
        "deploymentGroupName": "MyServiceDeploymentGroup",
        "deploymentConfigName": "CodeDeployDefault.OneAtATime",
        "ec2TagFilters": [],
        "onPremisesInstanceTagFilters": [],
        "autoScalingGroups": [
            {
                "name": "MyServiceAutoscalingGroup20181218054948246300000002",
                "hook": "CodeDeploy-managed-automatic-launch-deployment-hook-myservice-a2d358c8-3525-452c-b76e-978f1746ae74"
            }
        ],
   ...
   }
```

We see that our deployment group has been created, has been associated with the autoscaling group, and we have a hook associated with it.
