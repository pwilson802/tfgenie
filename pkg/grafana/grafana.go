package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

type GrafanaClient struct {
	Hostname string // The hostname of grafana to connect to
	ApiKey   string
}

func CreateNewGrafanaClient(hostname string) *GrafanaClient {
	apiKey := os.Getenv("GRAFANA_API_KEY")
	return &GrafanaClient{Hostname: hostname, ApiKey: apiKey}
}

func (g *GrafanaClient) GetAlert(uid string) (*Alert, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/provisioning/alert-rules/%s", g.Hostname, uid), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Check if the response status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get alert: status code %d", resp.StatusCode)
	}

	var alert Alert
	if err := json.NewDecoder(resp.Body).Decode(&alert); err != nil {
		return nil, err
	}

	return &alert, nil

}

type Alert struct {
	ID           int               `json:"id"`
	UID          string            `json:"uid"`
	OrgID        int               `json:"orgID"`
	FolderUID    string            `json:"folderUID"`
	RuleGroup    string            `json:"ruleGroup"`
	Title        string            `json:"title"`
	Condition    string            `json:"condition"`
	Data         []AlertData       `json:"data"`
	Updated      time.Time         `json:"updated"`
	NoDataState  string            `json:"noDataState"`
	ExecErrState string            `json:"execErrState"`
	For          string            `json:"for"`
	Annotations  map[string]string `json:"annotations"`
	Labels       map[string]string `json:"labels"`
	IsPaused     bool              `json:"isPaused"`
}

// Struct representing each item in the "data" array.
type AlertData struct {
	RefID             string            `json:"refId"`
	QueryType         string            `json:"queryType"`
	RelativeTimeRange RelativeTimeRange `json:"relativeTimeRange"`
	DatasourceUID     string            `json:"datasourceUid"`
	Model             AlertModel        `json:"model"`
}

// Struct representing the relative time range in each "data" item.
type RelativeTimeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Struct representing the "model" object in each "data" item.
type AlertModel struct {
	Datasource       DataSource        `json:"datasource"`
	Dimensions       map[string]string `json:"dimensions,omitempty"`
	Expression       string            `json:"expression,omitempty"`
	Hide             bool              `json:"hide"`
	ID               string            `json:"id,omitempty"`
	IntervalMs       int               `json:"intervalMs"`
	Label            string            `json:"label,omitempty"`
	LogGroups        []string          `json:"logGroups,omitempty"`
	MatchExact       bool              `json:"matchExact,omitempty"`
	MaxDataPoints    int               `json:"maxDataPoints"`
	MetricEditorMode int               `json:"metricEditorMode,omitempty"`
	MetricName       string            `json:"metricName,omitempty"`
	MetricQueryType  int               `json:"metricQueryType,omitempty"`
	Namespace        string            `json:"namespace,omitempty"`
	Period           string            `json:"period,omitempty"`
	QueryMode        string            `json:"queryMode,omitempty"`
	RefID            string            `json:"refId"`
	Region           string            `json:"region,omitempty"`
	SqlExpression    string            `json:"sqlExpression,omitempty"`
	Statistic        string            `json:"statistic,omitempty"`
	Conditions       []Condition       `json:"conditions,omitempty"`
	Reducer          string            `json:"reducer,omitempty"`
	Type             string            `json:"type,omitempty"`
}

// Struct representing the "datasource" object in the "model".
type DataSource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// Struct representing a condition inside the "model".
type Condition struct {
	Evaluator Evaluator `json:"evaluator"`
	Operator  Operator  `json:"operator"`
	Query     Query     `json:"query"`
	Reducer   Reducer   `json:"reducer"`
	Type      string    `json:"type"`
}

// Struct representing the "evaluator" object in a condition.
type Evaluator struct {
	Params []interface{} `json:"params"`
	Type   string        `json:"type"`
}

// Struct representing the "operator" object in a condition.
type Operator struct {
	Type string `json:"type"`
}

// Struct representing the "query" object in a condition.
type Query struct {
	Params []string `json:"params"`
}

// Struct representing the "reducer" object in a condition.
type Reducer struct {
	Params []interface{} `json:"params"`
	Type   string        `json:"type"`
}

// ExportAlertToTerraform converts an Alert to a Terraform resource and writes it to a file.
func (g *GrafanaClient) ExportAlertToTerraform(alert *Alert, filename string, folderName string, ruleName string) error {
	if filename == "" {
		filename = "grafana-alert.txt"
	}
	if folderName == "" {
		folderName = "New Folder"
	}
	if ruleName == "" {
		ruleName = "New Rule"
	}

	tmpl := `
resource "grafana_folder" "{{ .FolderResourceName }}" {
  title = "{{ .FolderTitle }}"
}

resource "grafana_rule_group" "{{ .RuleGroupResourceName }}" {
  name             = "{{ .RuleGroupName }}"
  folder_uid       = grafana_folder.{{ .FolderResourceName }}.uid
  interval_seconds = 240
  rule {
    name           = "{{ .Alert.Title }}"
    for            = "{{ .Alert.For }}"
    condition      = "{{ .Alert.Condition }}"
    no_data_state  = "{{ .Alert.NoDataState }}"
    exec_err_state = "{{ .Alert.ExecErrState }}"
    annotations = {
      {{- range $key, $value := .Alert.Annotations }}
      "{{ $key }}" = "{{ $value }}"
      {{- end }}
    }
    labels = {
      {{- range $key, $value := .Alert.Labels }}
      "{{ $key }}" = "{{ $value }}"
      {{- end }}
    }
    is_paused = {{ .Alert.IsPaused }}
    {{- range $index, $data := .Alert.Data }}
    data {
      ref_id     = "{{ $data.RefID }}"
      query_type = "{{ $data.QueryType }}"
      relative_time_range {
        from = {{ $data.RelativeTimeRange.From }}
        to   = {{ $data.RelativeTimeRange.To }}
      }
      datasource_uid = "{{ $data.DatasourceUID }}"
      model = jsonencode({
        "datasource" : {
          "type" : "{{ $data.Model.Datasource.Type }}",
          "uid" : "{{ $data.Model.Datasource.UID }}"
        },
        "dimensions" : {
          {{- range $dimKey, $dimValue := $data.Model.Dimensions }}
          "{{ $dimKey }}" : "{{ $dimValue }}",
          {{- end }}
        },
        "expression" : "{{ $data.Model.Expression }}",
        "hide" : {{ $data.Model.Hide }},
        "id" : "{{ $data.Model.ID }}",
        "intervalMs" : {{ $data.Model.IntervalMs }},
        "label" : "{{ $data.Model.Label }}",
        "logGroups" : [
          {{- range $lg := $data.Model.LogGroups }}
          "{{ $lg }}",
          {{- end }}
        ],
        "matchExact" : {{ $data.Model.MatchExact }},
        "maxDataPoints" : {{ $data.Model.MaxDataPoints }},
        "metricEditorMode" : {{ $data.Model.MetricEditorMode }},
        "metricName" : "{{ $data.Model.MetricName }}",
        "metricQueryType" : {{ $data.Model.MetricQueryType }},
        "namespace" : "{{ $data.Model.Namespace }}",
        "period" : "{{ $data.Model.Period }}",
        "queryMode" : "{{ $data.Model.QueryMode }}",
        "refId" : "{{ $data.Model.RefID }}",
        "reducer" : "{{ $data.Model.Reducer }}",
        "region" : "{{ $data.Model.Region }}",
        "sqlExpression" : "{{ $data.Model.SqlExpression }}",
        "statistic" : "{{ $data.Model.Statistic }}",
        "conditions" : [
          {{- range $cond := $data.Model.Conditions }}
          {
            "evaluator" : {
              "params" : {{ $cond.Evaluator.Params }},
              "type" : "{{ $cond.Evaluator.Type }}"
            },
            "operator" : {
              "type" : "{{ $cond.Operator.Type }}"
            },
            "query" : {
              "params" : [
				{{- range $param := $cond.Query.Params }}
					"{{ $param }}",
				{{- end }} 
            	]
			},
            "reducer" : {
              "params" : {{ $cond.Reducer.Params }},
              "type" : "{{ $cond.Reducer.Type }}"
            },
            "type" : "{{ $cond.Type }}"
          }
          {{- end }}
        ],
        "type" : "{{ $data.Model.Type }}"
      })
    }
    {{- end }}
  }
}
`

	// Define template data structure
	data := struct {
		FolderResourceName    string
		FolderTitle           string
		RuleGroupResourceName string
		RuleGroupName         string
		Alert                 *Alert
	}{
		FolderResourceName:    strings.ToLower(strings.ReplaceAll(folderName, " ", "_")),
		FolderTitle:           folderName,
		RuleGroupResourceName: strings.ToLower(strings.ReplaceAll(ruleName, " ", "_")),
		RuleGroupName:         ruleName,
		Alert:                 alert,
	}

	// Parse and execute the template
	t, err := template.New("terraform").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Write the output to the file
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
