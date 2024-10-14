package config

import (
	"flag"
	"os"
)

// AWSConfig holds AWS S3 configuration
type AWSConfig struct {
	BucketName string
	Region     string
}

// AzureConfig holds Azure Blob configuration
type AzureConfig struct {
	ContainerName string
	AccountName   string
	AccountKey    string
}

// Config holds overall service configuration
type Config struct {
	StorageType   string
	FilePath      string
	AWSConfig     AWSConfig
	AzureConfig   AzureConfig
	ServerAddress string
}

// LoadConfig loads the configuration from command line flags
func LoadConfig() Config {
	var cfg Config

	// Define command line flags
	flag.StringVar(&cfg.StorageType, "storage-type", "memory", "Type of storage to use (memory, file, aws, azure)")
	flag.StringVar(&cfg.FilePath, "file-path", "events.log", "File path for local file storage")
	flag.StringVar(&cfg.AWSConfig.BucketName, "aws-bucket-name", "", "AWS S3 bucket name")
	flag.StringVar(&cfg.AWSConfig.Region, "aws-region", "", "AWS region")
	flag.StringVar(&cfg.AzureConfig.ContainerName, "azure-container-name", "", "Azure Blob container name")
	flag.StringVar(&cfg.AzureConfig.AccountName, "azure-account-name", "", "Azure storage account name")
	flag.StringVar(&cfg.AzureConfig.AccountKey, "azure-account-key", "", "Azure storage account key")
	flag.StringVar(&cfg.ServerAddress, "server-address", ":8080", "Server address")

	// Show help if needed
	showHelp := flag.Bool("h", false, "Display help")
	flag.BoolVar(showHelp, "help", false, "Display help")

	// Parse the flags
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	return cfg
}
