package main

const AssumeRolePolicy = `{
    "Version": "2012-10-17",
    "Statement": [
       	{
            "Action": "sts:AssumeRole",
            "Principal": {
       				"AWS": "%s"
            },
            "Effect": "Allow",
        	"Sid": ""
        }
    ]
}
`

const AdminPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "*",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
`
