package config

import (
	"flag"
	"os"
)

type AzureConfig struct {
	AccountName   string
	ContainerName string
	BlobName      string
	AccountKey    string
}

type LocalFileConfig struct {
	FilePath string
}

type Config struct {
	StorageType         string
	EstimatedEventCount uint
	ServerAddress       string
	LocalFileConfig     LocalFileConfig
	AzureConfig         AzureConfig
}

func LoadConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.StorageType, "storage-type", "memory", "Type of storage to use (memory, file, aws, azure)")
	flag.UintVar(&cfg.EstimatedEventCount, "estimated-event-count", 1000000, "Server address")
	flag.StringVar(&cfg.LocalFileConfig.FilePath, "file-path", "events.log", "File path for local file storage")
	flag.StringVar(&cfg.AzureConfig.AccountName, "azure-account-name", "", "Azure storage account name")
	flag.StringVar(&cfg.AzureConfig.ContainerName, "azure-container-name", "", "Azure Blob container name")
	flag.StringVar(&cfg.AzureConfig.BlobName, "azure-blob-name", "", "Azure Blob name")
	// no no no!
	// flag.StringVar(&cfg.AzureConfig.AccountKey, "azure-account-key", "", "Azure storage account key")
	flag.StringVar(&cfg.ServerAddress, "server-address", ":8080", "Server address")

	showHelp := flag.Bool("h", false, "Display help")
	flag.BoolVar(showHelp, "help", false, "Display help")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	return cfg
}
