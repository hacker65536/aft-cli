package myaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	pipelinetypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/hacker65536/aft-cli/pkg/logger"
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
	Accounts       = map[string]string{}
	ConcurrencyNum = 7
)

func New() *MyAWS {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.ZapLog.Fatal(err.Error())
	}
	return &MyAWS{
		Config: cfg,
	}
}
