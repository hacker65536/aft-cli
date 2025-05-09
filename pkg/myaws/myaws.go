package myaws

import (
	"context"
	"fmt"
	"sync"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	org "github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/hacker65536/aft-cli/pkg/logger"
	"github.com/hacker65536/aft-cli/pkg/util"
)

func (m *MyAWS) ListAccounts() {

	cacheTTL := 1 * time.Hour

	cache, err := util.NewFileCache(cacheTTL)
	if err != nil {
		logger.ZapLog.Error(err.Error())
	}
	key := "accounts"

	if data, err := cache.Get(key); err == nil {
		// Use cached data
		m.Accounts = data.Accounts
		logger.ZapLog.Debug("data from cache")
	} else {
		// Fetch or compute data
		client := org.NewFromConfig(m.Config)

		input := org.ListAccountsInput{}

		p := org.NewListAccountsPaginator(client, &input, func(o *org.ListAccountsPaginatorOptions) {
			//	o.Limit = 10
		})
		pageNum := 0

		for p.HasMorePages() {
			pageNum++
			output, err := p.NextPage(context.TODO())
			if err != nil {
				logger.ZapLog.Fatal(err.Error())
			}
			for _, account := range output.Accounts {
				if account.Status == "ACTIVE" {
					m.Accounts = append(m.Accounts, account)
				}

			}
		}
		for _, v := range m.Accounts {
			Accounts[*v.Id+"-customizations-pipeline"] = *v.Name
			//		fmt.Println(*v.Id, *v.Name, *v.JoinedTimestamp)
		}
		//		data := []byte("some data to cache")

		// Create AWSAccounts object and populate it with the account data
		awsAccounts := &util.AWSAccounts{
			Accounts: m.Accounts,
		}

		// Store data in cache
		err := cache.Set(key, awsAccounts)
		if err != nil {
			panic(err)
		}
		logger.ZapLog.Debug("data stored in cache")
	}

}

func (m *MyAWS) ListCodePipelines() {
	client := codepipeline.NewFromConfig(m.Config)

	input := codepipeline.ListPipelinesInput{}
	p := codepipeline.NewListPipelinesPaginator(client, &input, func(o *codepipeline.ListPipelinesPaginatorOptions) {
		//o.Limit = 10
	})

	pageNum := 0

	for p.HasMorePages() {
		pageNum++
		output, err := p.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}

		//m.Pipelines = append(m.Pipelines, output.Pipelines...)
		for _, pipeline := range output.Pipelines {
			if strings.HasSuffix(*pipeline.Name, "-customizations-pipeline") {
				m.Pipelines = append(m.Pipelines, pipeline)
			}

		}
	}
	/*
		for _, v := range m.Pipelines {
			fmt.Println(*v.Name,)
		}
	*/
}

func (m *MyAWS) AftPipelineStatus() {

	m.ListAccounts()
	//	fmt.Println("listaccount")
	m.ListCodePipelines()
	//	fmt.Println("listpipeline")

	ch := make(chan struct{}, ConcurrencyNum)
	var wg sync.WaitGroup

	for _, v := range m.Pipelines {
		ch <- struct{}{}
		wg.Add(1)
		id := *v.Name
		go func(id string) {

			defer func() {
				<-ch
				wg.Done()
			}()
			//	fmt.Println("doing")
			//	fmt.Println(*v.Name)
			//  time.Sleep(300 * time.Millisecond)
			//	fmt.Println("done")

			//			m.ListPipelineExecution(id)
			m.GetPipelineStateWithAccount(id)
		}(id)
	}
	wg.Wait()

	for _, v := range m.PipelineWithAccounts {
		//		fmt.Println("Execution")
		/*
			for _, execution := range v.Executions {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\t", Accounts[v.PipelineName], v.PipelineName, *execution.PipelineExecutionId, execution.Status, *execution.LastUpdateTime)
			}
		*/

		fmt.Printf("%s\t%s\t", Accounts[v.PipelineName], v.PipelineName)
		reason := ""
		flg := false
		for _, status := range v.StageState {
			fmt.Printf("%s\t%s\t", *status.StageName, status.LatestExecution.Status)
			if status.LatestExecution.Status == "Failed" {
				flg = true
				reason += "  StageName: " + *status.StageName + "\n"
			}

			for _, action := range status.ActionStates {
				if action.LatestExecution.Status == "Failed" {

					message := *action.LatestExecution.ErrorDetails.Message
					messagePreview := message
					if len(message) > 20 {
						messagePreview = message[:20]
					}
					reason += "    ActionName: " + *action.ActionName + " ErrorCode: " + *action.LatestExecution.ErrorDetails.Code + " ErrorMessage: " + messagePreview + "\n"
				}
			}
		}
		fmt.Println()
		if flg {
			fmt.Printf("%s\n", reason)
		}
		flg = false
		reason = ""
	}

}

func (m *MyAWS) GetPipelineStateWithAccount(pipelineName string) {
	client := codepipeline.NewFromConfig(m.Config)
	optFn := func(o *codepipeline.Options) {
		o.RetryMaxAttempts = 5
		o.RetryMode = aws.RetryModeStandard

	}
	input := codepipeline.GetPipelineStateInput{
		Name: &pipelineName,
	}
	output, err := client.GetPipelineState(context.TODO(), &input, optFn)
	if err != nil {
		panic(err)
	}
	//fmt.Println(output)
	m.PipelineWithAccounts = append(m.PipelineWithAccounts, PipelineWithAccount{
		PipelineName: pipelineName,
		StageState:   output.StageStates,
	})

}

func (m *MyAWS) ListPipelineExecution(pipelineName string) {

	client := codepipeline.NewFromConfig(m.Config)
	optFn := func(o *codepipeline.Options) {
		o.RetryMaxAttempts = 5
		o.RetryMode = aws.RetryModeStandard
	}
	maxResults := int32(1)
	input := codepipeline.ListPipelineExecutionsInput{
		PipelineName: &pipelineName,
		MaxResults:   &maxResults,
	}
	output, err := client.ListPipelineExecutions(context.TODO(), &input, optFn)
	if err != nil {
		panic(err)
	}
	m.PipelineWithAccounts = append(m.PipelineWithAccounts, PipelineWithAccount{
		PipelineName: pipelineName,
		Executions:   output.PipelineExecutionSummaries,
	})

	//	pageNum := 0
	//	for _, pipeline := range m.Pipelines {
	//		maxResults := int32(1)
	//		input := codepipeline.ListPipelineExecutionsInput{
	//			PipelineName: pipeline.Name,
	//			MaxResults:   &maxResults,
	//		}
	//		/*
	//			p := codepipeline.NewListPipelineExecutionsPaginator(client, &input, func(o *codepipeline.ListPipelineExecutionsPaginatorOptions) {
	//				o.Limit = 5
	//			})
	//
	//			pageNum := 0
	//
	//			for p.HasMorePages() {
	//				pageNum++
	//				output, err := p.NextPage(context.TODO())
	//				if err != nil {
	//					panic(err)
	//				}
	//
	//				m.PipelineExecutionsWithAccounts = append(m.PipelineExecutionsWithAccounts, PipelineExecutionsWithAccount{
	//					PipelineName: *pipeline.Name,
	//					Executions:   output.PipelineExecutionSummaries,
	//				})
	//			}
	//		*/
	//		output, err := client.ListPipelineExecutions(context.TODO(), &input)
	//		if err != nil {
	//			panic(err)
	//		}
	//		m.PipelineExecutionsWithAccounts = append(m.PipelineExecutionsWithAccounts, PipelineExecutionsWithAccount{
	//			PipelineName: *pipeline.Name,
	//			Executions:   output.PipelineExecutionSummaries,
	//		})
	//		pageNum++
	//		fmt.Println(pageNum)
	//	}

}
