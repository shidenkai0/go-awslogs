# go-awslogs

A simple and reliable AWS Cloudwatch Logs fetcher.

## Prerequisites

Setup IAM credentials allowing the Cloudwatch Logs FilterLogEvents API call on the Log Groups you wish to query, either through the AWS CLI (```aws configure```), or by manually setting your credentials store:(http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html).

## Installation

```
go get github.com/shidenkai0/go-awslogs
cd $GOPATH/src/go-awslogs
go install
```

You should be good to go.

## How to use

``` 
Usage of go-awslogs:
  -end int
    	end unix time (seconds)
  -f string
    	Output file name (optional)
  -filter string
    	Filter pattern (checkout http://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)
  -log-group-name string
    	Log Group Name
  -region string
    	AWS Region (default "eu-west-1")
  -start int
    	start unix time (seconds)
```