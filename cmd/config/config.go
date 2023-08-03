package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"time"
)

type Slack struct {

	// BotToken is the bot user OAuth access token.
	//
	// Can be found on the OAuth & Permissions page:
	// https://api.slack.com/apps/<your-app-id>/oauth
	BotToken string `yaml:"slack.botToken" env:"SLACK_BOT_TOKEN" env-required:"true"`

	// SigningSecret is the signing secret used to verify requests from Slack.
	//
	// Can be found on the Basic Information page:
	// https://api.slack.com/apps/<your-app-id>/general
	SigningSecret string `yaml:"slack.signingSecret" env:"SLACK_SIGNING_SECRET" env-required:"true"`

	// Custom is the custom configuration for the bot.
	Custom struct {

		// SuccessMsg is the message that the bot will send to the user if the data is found.
		SuccessMsg string `yaml:"slack.custom.successMsg" env:"SLACK_CUSTOM_SUCCESS_MSG" env-default:"Here is your data:"`
	}
}

type Google struct {

	// Credentials is the string containing the Google service account credentials.
	//
	// Shortcut for passing credentialsFile content directly from environment variables.
	// Can be downloaded from the Google Cloud Console:
	// https://console.cloud.google.com/iam-admin/serviceaccounts
	Credentials []byte `yaml:"google.credentialString" env:"GOOGLE_CREDENTIALS_STRING" env-required:"true"`

	// SpreadsheetID is the ID of the Google spreadsheet.
	//
	// Can access using the URL of the spreadsheet:
	// https://docs.google.com/spreadsheets/d/<spreadsheetID>/edit#gid=0
	SpreadsheetID string `yaml:"google.spreadsheetID" env:"GOOGLE_SPREADSHEET_ID" env-required:"true"`

	// SpreadsheetRange is the range of the Google spreadsheet.
	//
	// For example, "Sheet1!A1:B2" will select the range between A1 and B2 on Sheet1 page
	SpreadsheetRange string `yaml:"google.spreadsheetRange" env:"GOOGLE_SPREADSHEET_RANGE" env-required:"true"`

	// Custom is the custom configuration for the bot.
	Custom struct {

		// SearchingValueFrom is the name of the header that the bot will search passed value.
		//
		// For example, if the header is "Name", then the bot will search the value from the "Name" column.
		// The header name is case-sensitive.
		SearchingValueFrom string `yaml:"google.searchingValueFrom" env:"GOOGLE_SEARCHING_VALUE_FROM" env-required:"true"`

		// TakingValueFrom is the name of the header that the bot will take the value from.
		//
		// For example, if the header is "Email", then the bot will take the value from the "Email" column.
		//
		// Usecase:
		// 	In combination with SearchingValueFrom, the bot will search the value from the "Name" column
		// 	and take the value from the "Email" column.
		TakingValueFrom string `yaml:"google.takingValueFrom" env:"GOOGLE_TAKING_VALUE_FROM" env-required:"true"`

		// HeaderRowIndex is the index of the row containing the headers.
		//
		// For example, if the headers are on the first row, then the index is 0.
		HeaderRowIndex int `yaml:"google.headerRowIndex" env:"GOOGLE_HEADER_ROW_INDEX" env-default:"0"`

		// UseEmailAsSearchingValue is the flag to use email as searching value.
		//
		// If true, the bot will use the user's email as the searching criteria for the SearchingValueFrom header.
		UseEmailAsSearchingValue bool `yaml:"google.useEmailAsSearchingValue" env:"GOOGLE_USE_EMAIL_AS_SEARCHING_VALUE" env-default:"true"`
	}
}

type Server struct {

	// Port is the port that the server will listen on.
	Port int `yaml:"server.port" env:"SERVER_PORT" env-default:"8080"`

	// ReadTimeout is the maximum duration for reading the entire request, including the body.
	ReadTimeout time.Duration `yaml:"server.readTimeout" env:"SERVER_READ_TIMEOUT" env-default:"5s"`

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration `yaml:"server.writeTimeout" env:"SERVER_WRITE_TIMEOUT" env-default:"5s"`
}

type GlobalConfig struct {
	Server Server
	Slack  Slack
	Google Google
}

func NewConfig() (*GlobalConfig, error) {
	_ = godotenv.Load(".env")

	var cfg GlobalConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
