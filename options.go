package main

type Options struct {
	Port string `long:"port" default:"8080" env:"PORT"`

	AllowedEmailRegex string `long:"allowed-email-regex" default:".*" description:"The regex of allowed emails"`

	StorageType string `long:"storage-type" default:"inmem"`

	// BigQuery Config options
	BigQuery struct {
		ProjectID string `long:"bq-project-id" description:"The GCP project id"`

		DatasetName string `long:"bq-dataset" default:"golinks"`

		TableName string `long:"bq-table" default:"kvs"`
	} `group:"BigQuery Storage Options"`
}
