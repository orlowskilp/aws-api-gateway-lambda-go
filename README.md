# REST API with AWS Lambda and DynamoDB integration

![Build and Deploy workflow status](https://github.com/orlowskilp/aws-api-gateway-lambda-go/workflows/Build%20and%20Deploy/badge.svg)
![Quick check workflow status](https://github.com/orlowskilp/aws-api-gateway-lambda-go/workflows/Quick%20check/badge.svg)
![Check IaC workflow status](https://github.com/orlowskilp/aws-api-gateway-lambda-go/workflows/Check%20IaC/badge.svg)

Sample serverless API based using API Gateway, AWS Lambda (implemented in Go) and AWS DynamoDB.  

## Technology stack

This example, demonstrates how to deploy an AWS-native serverless API using:
* [AWS API Gateway](https://aws.amazon.com/api-gateway/) for API deployment
* [AWS Lambda](https://aws.amazon.com/lambda/) for HTTP request handling
* [Go programming language](https://golang.org/) for [Lambda function implementation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
* [AWS DynamoDB](https://aws.amazon.com/dynamodb/) for data storage
* [AWS CloudFormation](https://aws.amazon.com/cloudformation/) to represent infrastructure as code (using YAML)

Toolchain used for development comprises:
* [Visual Studio Code](https://code.visualstudio.com/)
* [Docker](https://www.docker.com/) for [development in containers](https://code.visualstudio.com/docs/remote/containers)
* GNU make
* [AWS CLI](https://aws.amazon.com/cli/) with some bash script to manage DynamoDB tables
* GitHub Actions

## Solution architecture overview

<p align="center"><img src="https://github.com/orlowskilp/aws-api-gateway-lambda-go/blob/master/resources/architecture.svg" /></p>

The architecture of the solution presented in this example is straightforward. A REST API is 
deployed with API Gateway and the endpoints point to a Lambda function handling requests. 
The function, in turn, talks to a DynamoDB table and logs are stored in a CloudWatch log stream.

The API implements a trivial CRUD interface (using appropriate HTTP methods) and specifies
item key using `key` path parameter. The items are stored in a DynamoDB table in which the
hash key is a `string` value named `Key` and the data is stored in the field `Value` as `string`.

Here are some important architectural considerations:
* REST API accepts `GET`, `PUT` and `DELETE` methods
* `GET` method with specified `key` retrieves an item assigned to key
* `PUT` method with specified `key` assigns an item to key
* `DELETE` method with specified `key` assigned an item assigned to key
* API Gateway proxies requests (using `AWS_PROXY`) to the Lambda function for processing.
The function logic selects an appropriate handler depending on type of HTTP method used.
* The requests are proxied to Lambda function with `POST` method (required by AWS Lambda)

## Development environment for Visual Studio Code

My team and I love using [Microsoft Visual Studio Code](https://code.visualstudio.com/). 
We also, typically, develop on different platforms than our target platforms. VSCode 
allows us to encapsulate the development environment, by leveraging the power of containers.

```
.devcontainer
├── Dockerfile
├── devcontainer.json
└── docker-compose.yml
```

The `.devcontainer` directory contains the `Dockerfile`, managing software stack in the
development container, `docker-compose.yml`, for Docker Compose to orchestrate the development 
container and local deployment of AWS DynamoDB, and `devcontainer.json` file, configuring the
development environment.

Install `ms-vscode-remote.remote-containers` in your VSCode IDE and then 
**open the repository directory**. VSCode will offer to reopen the directory in container.
This gives you a development environment with all the dependencies pre-installed.

Here we're going to use Docker Compose, to orchestrate a 2 container environment. One container 
is the development container, whose command line we'll be accessing with VSCode and the other
one is the official _amazon/dynamodb-local_ providing local DynamoDB database.

### Using local DynamoDB

The local DynamoDB accessible under `http://dynamodb:8000`. To check whether it works, you can
do:

```
$ aws dynamodb list-tables --endpoint http://dynamodb:8000
```

If you didn't previously create any tables, you should get the following output:

```
{
    "TableNames": []
}
```

You can read more on DynamoDB CLI [here](https://docs.aws.amazon.com/cli/latest/reference/dynamodb/index.html).

### Helper scripts for DynamoDB table management

The helper scripts take care of managing the table for you. They are extremely simple, so if you
need something more advanced, you may need to modify them.

```
dynamodb
├── common.sh
├── create-table.sh
├── delete-table.sh
└── populate-table.sh
```

The `dynamodb` directory contains the scripts. The `common.sh` file stores variables used by the
helper scripts. If you orchestrate your DynamoDB table under different address or want to use a
different name of the table, you'll have to update that in that file.

`create-table.sh` creates a table by default named `sample-table`. It expects that there is no such
table.

`populate-table.sh` stores a couple of entries in the table, to have something to work with. It
expects that a table already exists.

`delete-table.sh` deletes a table. It expects that the table exists

You can use the scripts to manage your table in AWS. To do so, just run respective scripts with `-r`
or `--remote` flag. Note, that you will need to have your credentials configured to access an
appropriate table in AWS.

**NOTE**: Using helper scripts with the `-r` or `--remote` flag will permanently override previously
existing data.

## API overview

The REST API in this example exposes the following endpoints:

```
/
GET
└── key/
    └── {key}
        GET
        PUT
        DELETE
```

* `GET https://{rest-api-id}.execute-api.{aws-region}.amazonaws.com/{api-stage-name}`  
This method is only for deployment testing. It ignores any query parameters
* `GET https://{rest-api-id}.execute-api.{aws-region}.amazonaws.com/{api-stage-name}/key/*`  
This method retrieves an item stored under specified key. It ignores any query parameters
* `PUT https://{rest-api-id}.execute-api.{aws-region}.amazonaws.com/{api-stage-name}/key/*`  
This method places an item under specified key. Value is passed in request body.  It ignores any
query parameters
* `DELETE https://{rest-api-id}.execute-api.{aws-region}.amazonaws.com/{api-stage-name}/key/*`  
This method deletes an item stored under specified key. It ignores request body.  It ignores any
query parameters

## AWS Lambda function code

The following files are relevant for the Lambda function:

```
main.go
Makefile
go.mod
go.sum
pkg
└── dynamodb
    ├── dynamodb.go
    ├── dynamodb_integration_test.go
    ├── dynamodb_mock.go
    └── dynamodb_test.go
```

The `dynamodb` package is stored in the `pkg` directory.

To build the code run:

```
$ make all
```

This will produce the binary and a zipfile with the compressed binary (for AWS Lambda deployment).  
To provide a custom name to the binary and the zipfile, you can do:

```
$ make TARGET=custom-name all
```

### Unit tests

Unit tests use mocks implemented with [gomock](https://github.com/golang/mock) and
[go-dynamock](https://github.com/gusaul/go-dynamock) (for DynamoDB mocks). The source
file with mocks `dynamodb_mock.go` is generated using `mockgen` and should not be manually
modified! To re-generate a the `dynamodb_mock.go` file run:

```
$ make mock
```

To run unit tests, run:

```
$ make test
```

### Integration tests

Integration tests talk to the local DynamoDB table. To run integration tests, run:

```
$ make integration_test
```

**NOTE**: Make sure that there are no tables in local DynamoDB instance, before running
integration tests, otherwise the integration tests will fail.

## AWS CloudFormation templates for infrastructure hosted in AWS

The `cfn` directory contains all the CloudFormation templates to deploy necessary resources
in AWS.

```
cfn
├── api
│   └── api.yml
├── backend
│   └── lambda-dynamodb-iam.yml
└── cicd
    ├── artifact-bucket.yml
    └── principal.yml
```

### CI/CD resources

The templates in the `cicd` directory deploy resources necessary to get the CI/CD
workflows running.

The `artifact-bucket.yml` template creates an S3 bucket to store build artifacts.

The `principal.yml` template creates and IAM user and policy to allow CI/CD workflow
to access AWS resources. The policy allows the CI/CD user to:
* Write to the S3 bucket for artifacts
* Update code of Lambda function deployed in AWS
* Validate CloudFormation templates

**NOTE**: You will need to manually create access keys for the user created by the
`principal.yml` template.

### Backend resources

The `lambda-dynamodb-iam.yml` template creates all the necessary backend resources,
which are:
* The DynamoDB table
* Policy for the Lambda execution role, allowing access to DynamoDB and CloudWatch
log streams
* Execution role for Lambda function
* Lambda function

### Frontend resources

The `api.yml` template creates the REST API resources:
* API Gateway resources
* Lambda permissions, allowing API Gateway methods execute the Lambda function

**NOTE**: The stack created with `api.yml` depends on the one created with
`lambda-dynamodb-iam.yml`. You need to supply the name of the backend stack
to the frontend stack, to resolve cross-stack references.

## CI/CD workflows

There are three workflows implemented for this repository:
* _Quick check_ - checks if the code builds and unit tests pass
* _IaC check_ - checks of the CloudFormation templates are valid
* _Build and Deploy_ - Tests, builds, packages and deploys code

```
.github
├── CODEOWNERS
└── workflows
    ├── build.yml
    ├── check.yml
    └── infra-check.yml
```

### Environment variables and secrets

The only notable environment variable set in the `build.yml` file is `SERVICE_NAME`.
It controls the name of the Lambda function and Lambda function handler name.

The following secrets need to be set before the CI/CD workflows can run:
* `AWS_ACCESS_KEY_ID` - Access key ID of the CI/CD user
* `AWS_SECRET_ACCESS_KEY` - Secret access key of the CI/CD user
* `AWS_ACCOUNT_NO` - Your AWS account nubmber
* `S3_BUCKET_NAME` - Name of the S3 bucket for artifacts (in the `s3://bucket.name`
format)

## Getting everything running

1. Create stack from the `artifact-bucket.yml` template
2. Create stack from the `principal.yml` template
3. Create access keys for the user created by stack in step (2)
4. Set secrets for the repository
5. Run the _Build and Deploy_ (e.g. by pushing an empty `staging` branch).
It will fail, but before that it will upload a zipfile with Lambda function
code to the S3 bucket for artifacts
6. Create stack from the `lambda-dynamodb-iam.yml` template
7. Create stack from the `api.yml` template