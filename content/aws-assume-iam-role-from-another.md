Title: Setting up AWS EC2 Assume Role with Terraform
Date: 2018-02-26 
Category: infrastructure
Status: Draft

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
on a single EC2 instance.

## The Problem

As I note in the diagram above, the individual services will get an access denied error since
the EC2 instance above is using a different IAM profile than the ones used when the services
are running in the production setup. There are two solutions to this problem:

- The first is to create an IAM profile which will have all the IAM policies of the constituent services
- The second is to use [AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)

We are going to see how we can implement the second approach. 

## Solution: Infrastructure setup

There are two stages to implement this solution. The first stage is to setup the infrastructure to allow the
assume role operation to succeed. This basically means that if an IAM role, `role1` wants to assume another
role, `role2`, then:

- `role1` should be allowed to perform the `sts:AssumeRole` action on `role2`
- `role2` should allow `role1` to assume itself


```
resource "aws_iam_instance_profile" "iam_profile1" {
  name  = "test_profile1"
  role = "${aws_iam_role.role1.name}"
}

resource "aws_iam_role" "role1" {
  name = "test_profile1_role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

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

resource "aws_iam_instance_profile" "iam_profile2" {
  name  = "test_profile2"
  role = "${aws_iam_role.role2.name}"
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

data "aws_iam_policy_document" "conf_bucket_access" {
  statement {
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      "${aws_s3_bucket.conf_bucket.arn}",
    ]

  }

  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "${aws_s3_bucket.conf_bucket.arn}",
      "${aws_s3_bucket.conf_bucket.arn}/*",
    ]
  }
}

resource "aws_iam_role_policy" "conf_bucket_access_policy" {
  name   = "AccessConfigurationBucket"
  role = "${aws_iam_role.role2.name}"
  policy = "${data.aws_iam_policy_document.conf_bucket_access.json}"
}
```


## Solution: How does it work?


## Demos using AWS CLI


