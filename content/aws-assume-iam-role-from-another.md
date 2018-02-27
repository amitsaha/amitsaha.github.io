Title: Setting up AWS EC2 Assume Role with Terraform
Date: 2018-02-26 
Category: infrastructure
Status: Draft

In this post, we will see how we can implement the AWS `assume role` functionality which allows
an IAM role to be able to obtain temporary credentials to access a resource otherwise only accessible
by another IAM role. We will implement the infrastructure changes using [Terraform](terraform.io)
and see how to obtain temporary credentials and access an AWS resource (a S3 bucket) that the corresponding
IAM role doesn't have access to otherwise via the [AWS CLI](https://aws.amazon.com/documentation/cli/).

If you want to follow along, please install `terraform` and setup AWS config so that it has a profile named
`dev`. If you have a profile or want to use a different AWS profile, you can change it in the `aws.tf` file
of the configuration you are applying. If you are like me, you may have trouble setting up the profile
correctly, so here's my two config files:

```
# ~/.aws/config

[profile dev]
region=ap-southeast-2


# ~/.aws/credentials

[dev]
aws_access_key_id=<Your access key>
aws_secret_access_key=<Your secret key>
```

In each of the configuration directory, you will have to run `terraform init` 
before you can run `terraform apply`. 

## Problem Scenario

Consider the following scenario for 3 services running on their own AWS EC2 instances in
a production setup:

```
                       ┌───────────────────────────┐
                       │   Production AWS Setup    │
                       └───────────────────────────┘




 .───────────────────.      .───────────────────.    .───────────────────.
(      S3Bucket1      )    (      S3Bucket2      )  (      S3Bucket3      )
 `───────────────────'      `───────────────────'    `───────────────────'
            ▲                         ▲                        ▲
      ┌ ─ ─ ┴ ─ ─ ┐             ┌ ─ ─ ┴ ─ ─ ┐            ┌ ─ ─ ┴ ─ ─ ┐
       IAM Role 1                IAM Role 2               IAM Role 3
      └ ─ ─ ┬ ─ ─ ┘             └ ─ ─ ┬ ─ ─ ┘            └ ─ ─ ┬ ─ ─ ┘
            │                         │                        │

     ┌────────────┐            ┌────────────┐            ┌────────────┐
     │ Service A  │            │ Service B  │            │ Service C  │
     └────────────┘            └────────────┘            └────────────┘

       ◀─ ─ ─ ─ ─ ─ ─ ─ ─    AWS EC2 Instances   ─ ─ ─ ─ ─ ─ ─ ─▶
```

Each service running on their own EC2 instance has their own AWS IAM profile which via 
their role and role policy gives them access to the corresponding S3 bucket.

Now, consider the setup below for a developer environment for the above services:

```

                       ┌───────────────────────────┐
                       │   Development AWS Setup   │
                       └───────────────────────────┘




 .───────────────────.      .───────────────────.    .───────────────────.
(      S3Bucket1      )    (      S3Bucket2      )  (      S3Bucket3      )
 `───────────────────'      `───────────────────'    `───────────────────'
              ▲                      ▲                       ▲
              │                      │                       │
      ┌ ─ ─ ─ ┴ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─│─
                            Access Denied                      │
      └ ─ ─ ─ ┬ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─│─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─│─
              │                 ┌ ─ ─ ─ ─ ─ ┐                │
              └────────────────   IAM Role   ────────────────┘
                                └ ─ ─ ─ ─ ─ ┘
                                     ▲
                                     │
                               ┌────────────┐
                               │ Service B  │
                               └────────────┘
                               ┌────────────┐
                               │ Service A  │
                               └────────────┘
                               ┌────────────┐
                               │ Service C  │
                               └────────────┘


                              AWS EC2 Instance

```

Instead of each service running on their own development EC2 instance, we run all the services
on a single EC2 instance. This will lead to access denied errors when `service A` tries to access
`S3Bucket1` and similarly for the other services when they try to access their corresponding
buckets.


## Setting up test infrastructure

As I note in the diagram above, the individual services will get an access denied error since
the EC2 instance above is using a different IAM profile than the ones used when the services
are running in the production setup. Let's see that for ourselves first by created a test 
infrastructure as follows:


- Create a S3 bucket
- Create two IAM profiles, `role1` and `role2`
- Add a policy to `role2` to be able to perform all operations on the S3 bucket
- Spin up an EC2 instance using `role1`

The [terraform](https://terraform.io) configuration for setting up the above infrastructure can be found 
[here](https://github.com/amitsaha/aws-assume-role-demo/tree/master/terraform_configuration/problem_demo). 
If we now try to access the S3 bucket from the EC2 instance via the AWS CLI, we will get:

```
$ ssh ec2-user@<Public-IP>
..

$ aws s3 ls s3://github-amitsaha-bucket/*
An error occurred (AccessDenied) when calling the ListObjects operation: Access Denied
```

There are two solutions to this problem:

- The first is to create an IAM profile which will have all the IAM policies of the constituent services
- The second is to use [AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)

The first solution, although simple has the main problem of duplicating your IAM policies and it doesn't
feel clean. The second approach, although requires some work is a much cleaner approach.

There are two stages to implement this solution. The first stage is to setup the infrastructure to allow the
assume role operation to succeed. If an IAM role, `role1` wants to assume another

## Solution: Infrastructure setup

If an IAM role, `role1` wants to assume another role, `role2`, then:

- `role1` should be allowed to perform the `sts:AssumeRole` action on `role2`
- `role2` should allow `role1` to assume itself

The corresponding IAM configuration earlier will be updated as follows:


```
data "aws_iam_policy_document" "assume_role2_policy" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]
    resources = [
      "${aws_iam_role.role2.arn}",
    ]
  }
}

resource "aws_iam_role_policy" "role1_assume_role2" {
  name   = "AssumeRole2"
  role = "${aws_iam_role.role1.name}"
  policy = "${data.aws_iam_policy_document.assume_role2_policy.json}"
}

resource "aws_iam_role" "role2" {
  name = "test_profile2_role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com",
               "AWS": "${aws_iam_role.role1.arn}"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}
```

The updated terraform configuration can be found [here](https://github.com/amitsaha/aws-assume-role-demo/tree/master/terraform_configuration/solution_demo).
Let's apply the changes:






```
[ec2-user@ip-172-31-6-239 ~]$ aws s3 ls s3://github-amitsaha-bucket/*
An error occurred (AccessDenied) when calling the ListObjects operation: Access Denied
```

## Solution: Perform AssumeRole operation


```
[ec2-user@ip-172-31-6-239 ~]$ aws sts assume-role --role-arn arn:aws:iam::033145145979:role/test_profile2_role --role-session-name s3-example

An error occurred (AccessDenied) when calling the AssumeRole operation: User: arn:aws:sts::033145145979:assumed-role/test_profile1_role/i-0b3f9c44566fbc260 is not authorized to perform: sts:AssumeRole on resource: arn:aws:iam::033145145979:role/test_profile2_role

[ec2-user@ip-172-31-6-239 ~]$ aws sts assume-role --role-arn arn:aws:iam::033145145979:role/test_profile2_role --role-session-name s3-example
{
    "AssumedRoleUser": {
        "AssumedRoleId": "AROAJ3CMHLQFMYPPWQLSQ:s3-example", 
        "Arn": "arn:aws:sts::033145145979:assumed-role/test_profile2_role/s3-example"
    }, 
    "Credentials": {
        "SecretAccessKey": "PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo", 
        "SessionToken": "FQoDYXdzEI7//////////wEaDDqRJAWz11tovnatwSLuAUf1CIjLW0OI5dTCAh610HW7f3fBxglofbntqxCSJVyei1DafEjriLIskDzKoCdz6Y7F5Z/uyv/Ue7dCCCvXFpVYExwt82hE7yTGrYJB/oQl+bkMIzPhlHyegDa3/+vxdFu2kbcve8a1VlNhZE8fnpaRLGMoEr9/Ll+NQLjtRyysQ7DuN0GuMVIDiUzqOZHVDFDt4/c5LBHd2VZNfZ2t/rfPTkIwfkI9JQUVON+lcrk5W+FH16Onp1vuZXX4cmraMWQ1ROGf2x4fHGPIcMqaw674sgOnMSllyCUONLIaSPOeJLfOSDIrM/Xfv0PvslgotNrK1AU=", 
        "Expiration": "2018-02-25T13:33:56Z", 
        "AccessKeyId": "ASIAI7JVCNUGFT6XGMAQ"
    }
}
[ec2-user@ip-172-31-6-239 ~]$ date
Sun Feb 25 12:34:03 UTC 2018
[ec2-user@ip-172-31-6-239 ~]$ aws s3 ls s3://github-amitsaha-bucket/

An error occurred (AccessDenied) when calling the ListObjects operation: Access Denied
[ec2-user@ip-172-31-6-239 ~]$ AWS_ACCESS_KEY_ID=ASIAI7JVCNUGFT6XGMAQ AWS_SECRET_ACCESS_KEY=PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo aws s3 ls s3://github-amitsaha-bucket/

An error occurred (InvalidAccessKeyId) when calling the ListObjects operation: The AWS Access Key Id you provided does not exist in our records.
[ec2-user@ip-172-31-6-239 ~]$ AWS_ACCESS_KEY_ID=ASIAI7JVCNUGFT6XGMAQ AWS_SECRET_ACCESS_KEY=PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo aws s3 ls s3://github-amitsaha-bucket/

An error occurred (InvalidAccessKeyId) when calling the ListObjects operation: The AWS Access Key Id you provided does not exist in our records.
[ec2-user@ip-172-31-6-239 ~]$ AWS_SESSION_TOKEN="FQoDYXdzEI7//////////wEaDDqRJAWz11tovnatwSLuAUf1CIjLW0OI5dTCAh610HW7f3fBxglofbntqxCSJVyei1DafEjriLIskDzKoCdz6Y7F5Z/uyv/Ue7dCCCvXFpVYExwt82hE7yTGrYJB/oQl+bkMIzPhlHyegDa3/+vxdFu2kbcve8a1VlNhZE8fnpaRLGMoEr9/Ll+NQLjtRyysQ7DuN0GuMVIDiUzqOZHVDFDt4/c5LBHd2VZNfZ2t/rfPTkIwfkI9JQUVON+lcrk5W+FH16Onp1vuZXX4cmraMWQ1ROGf2x4fHGPIcMqaw674sgOnMSllyCUONLIaSPOeJLfOSDIrM/Xfv0PvslgotNrK1AU=" AWS_ACCESS_KEY_ID=ASIAI7JVCNUGFT6XGMAQ AWS_SECRET_ACCESS_KEY=PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo aws s3 ls s3://github-amitsaha-bucket/
[ec2-user@ip-172-31-6-239 ~]$ vim hello
[ec2-user@ip-172-31-6-239 ~]$ AWS_SESSION_TOKEN="FQoDYXdzEI7//////////wEaDDqRJAWz11tovnatwSLuAUf1CIjLW0OI5dTCAh610HW7f3fBxglofbntqxCSJVyei1DafEjriLIskDzKoCdz6Y7F5Z/uyv/Ue7dCCCvXFpVYExwt82hE7yTGrYJB/oQl+bkMIzPhlHyegDa3/+vxdFu2kbcve8a1VlNhZE8fnpaRLGMoEr9/Ll+NQLjtRyysQ7DuN0GuMVIDiUzqOZHVDFDt4/c5LBHd2VZNfZ2t/rfPTkIwfkI9JQUVON+lcrk5W+FH16Onp1vuZXX4cmraMWQ1ROGf2x4fHGPIcMqaw674sgOnMSllyCUONLIaSPOeJLfOSDIrM/Xfv0PvslgotNrK1AU=" AWS_ACCESS_KEY_ID=ASIAI7JVCNUGFT6XGMAQ AWS_SECRET_ACCESS_KEY=PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo aws s3 cp hello  s3://github-amitsaha-bucket/
upload: ./hello to s3://github-amitsaha-bucket/hello             
[ec2-user@ip-172-31-6-239 ~]$ AWS_SESSION_TOKEN="FQoDYXdzEI7//////////wEaDDqRJAWz11tovnatwSLuAUf1CIjLW0OI5dTCAh610HW7f3fBxglofbntqxCSJVyei1DafEjriLIskDzKoCdz6Y7F5Z/uyv/Ue7dCCCvXFpVYExwt82hE7yTGrYJB/oQl+bkMIzPhlHyegDa3/+vxdFu2kbcve8a1VlNhZE8fnpaRLGMoEr9/Ll+NQLjtRyysQ7DuN0GuMVIDiUzqOZHVDFDt4/c5LBHd2VZNfZ2t/rfPTkIwfkI9JQUVON+lcrk5W+FH16Onp1vuZXX4cmraMWQ1ROGf2x4fHGPIcMqaw674sgOnMSllyCUONLIaSPOeJLfOSDIrM/Xfv0PvslgotNrK1AU=" AWS_ACCESS_KEY_ID=ASIAI7JVCNUGFT6XGMAQ AWS_SECRET_ACCESS_KEY=PzFA0bJxxeB+i4kWjowpM6VTQTQfIiejbRxXkZdo aws s3 ls s3://github-amitsaha-bucket/
2018-02-25 12:38:32         12 hello
[ec2-user@ip-172-31-6-239 ~]$ 

```

