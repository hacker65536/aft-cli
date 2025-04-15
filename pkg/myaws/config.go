package myaws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	pipelinetypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type MyAWS struct {
	Accounts             []orgtypes.Account
	Config               aws.Config
	Pipelines            []pipelinetypes.PipelineSummary
	PipelineWithAccounts []PipelineWithAccount
}

type PipelineWithAccount struct {
	PipelineName string
	Executions   []pipelinetypes.PipelineExecutionSummary
	StageState   []pipelinetypes.StageState
}

var (
	Accounts = map[string]string{}
)

func New() *MyAWS {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return &MyAWS{
		Config: cfg,
	}
}
