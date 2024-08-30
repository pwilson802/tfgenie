package accessanalyzer

import (
	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer"
)

type AnalyzerClient struct {
	UserName   string
	AwsProfile string
	From       int
	To         int
	auth       *accessanalyzer.Client
}

// func (a *Analyzer) Auth {

// }
