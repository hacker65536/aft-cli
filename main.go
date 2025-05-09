/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/hacker65536/aft-cli/cmd"
	"github.com/hacker65536/aft-cli/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	lvl := zap.FatalLevel
	//lvl := zap.DebugLevel
	logger.InitializeLogger(lvl)
	defer logger.ZapLog.Sync() // Flush buffered logs
	cmd.Execute()
}
