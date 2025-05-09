package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/hacker65536/aft-cli/pkg/logger"
	"go.uber.org/zap"
)

type FileCache struct {
	path string
	ttl  time.Duration
	mu   sync.RWMutex
}

type AWSAccounts struct {
	Accounts []orgtypes.Account
}

func NewFileCache(ttl time.Duration) (*FileCache, error) {
	/*
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	*/
	appname := "aft-cli"
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Handle error appropriately, maybe return a default or error
			return nil, err
		}
		cacheDir = filepath.Join(homeDir, ".cache", appname)
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			//fmt.Println("Error creating directory:", err)
			return nil, err
		}
		logger.ZapLog.Debug("Create directory:", zap.String("cachedir", cacheDir))
	} else if err == nil {
		logger.ZapLog.Debug("Directory already exists:", zap.String("cachedir", cacheDir))
	} else {
		//fmt.Println("Error checking directory:", err)
		return nil, err
	}

	return &FileCache{
		path: cacheDir,
		ttl:  ttl,
	}, nil
}

func (c *FileCache) Get(key string) (*AWSAccounts, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := filepath.Join(c.path, key)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err

	}

	if time.Since(fileInfo.ModTime()) > c.ttl {
		os.Remove(filePath)
		return nil, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	var data *AWSAccounts
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return nil, err

	}

	return data, nil
}

func (c *FileCache) Set(key string, a *AWSAccounts) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	filePath := filepath.Join(c.path, key)
	file, err := os.Create(filePath)
	if err != nil {
		// Handle error
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(a)
	if err != nil {
		// Handle error
		return err
	}
	return nil
}
