Title: AWS EC2 assume role with Terraform
Date: 2018-02-26 
Category: infrastructure
Status: Draft

Consider the following scenario for N services running on their own AWS EC2 instances in
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

## Solution
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

## Implementation

## Demos using AWS CLI


