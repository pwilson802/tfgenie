package grafana_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pwilson802/tfgenie/pkg/grafana"
	"github.com/stretchr/testify/assert"
)

func TestClientCreation(t *testing.T) {
	t.Parallel()
	testClient := grafana.GrafanaClient{
		Hostname: "test.example.com",
	}
	want := "test.example.com"
	got := testClient.Hostname
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestGrafanaClient_GetAlert(t *testing.T) {
	// Create a mock alert response
	mockAlert := &grafana.Alert{
		ID:        529,
		UID:       "cbb410f8-3d68-4eb6-9334-9e2b6e05226f",
		OrgID:     1,
		FolderUID: "NiPuuXZVk",
		RuleGroup: "devops",
		Title:     "ECS Looping Test",
		Condition: "C",
		Data: []grafana.AlertData{
			{
				RefID:             "A",
				QueryType:         "",
				RelativeTimeRange: grafana.RelativeTimeRange{From: 600, To: 0},
				DatasourceUID:     "hsppouZ4z",
				Model: grafana.AlertModel{
					Datasource:    grafana.DataSource{Type: "cloudwatch", UID: "hsppouZ4z"},
					IntervalMs:    1000,
					MaxDataPoints: 43200,
					MetricName:    "PendingTaskCount",
					Namespace:     "ECS/ContainerInsights",
					Statistic:     "Maximum",
				},
			},
		},
		Updated:      time.Now(),
		NoDataState:  "NoData",
		ExecErrState: "Error",
		For:          "5m",
		Annotations:  map[string]string{"Account": "Development", "description": "ECS Task has been pending for over 15 minutes. Check the service."},
		Labels:       map[string]string{"alert": "test"},
		IsPaused:     false,
	}

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check the requested URL
		if r.URL.Path != "/api/v1/provisioning/alert-rules/cbb410f8-3d68-4eb6-9334-9e2b6e05226f" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Respond with the mock alert in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAlert)
	}))
	defer server.Close()

	// Set up the Grafana client with the mock server URL and a test API key
	client := &grafana.GrafanaClient{
		Hostname: server.URL,
		ApiKey:   "test-api-key",
	}

	// Call GetAlert
	alert, err := client.GetAlert("cbb410f8-3d68-4eb6-9334-9e2b6e05226f")

	// Assert no error
	assert.NoError(t, err)

	// Assert that the alert matches the mock alert, ignoring the Updated field
	assert.Equal(t, mockAlert.ID, alert.ID)
	assert.Equal(t, mockAlert.UID, alert.UID)
	assert.Equal(t, mockAlert.OrgID, alert.OrgID)
	assert.Equal(t, mockAlert.FolderUID, alert.FolderUID)
	assert.Equal(t, mockAlert.RuleGroup, alert.RuleGroup)
	assert.Equal(t, mockAlert.Title, alert.Title)
	assert.Equal(t, mockAlert.Condition, alert.Condition)
	assert.Equal(t, mockAlert.Data, alert.Data)
	assert.Equal(t, mockAlert.NoDataState, alert.NoDataState)
	assert.Equal(t, mockAlert.ExecErrState, alert.ExecErrState)
	assert.Equal(t, mockAlert.For, alert.For)
	assert.Equal(t, mockAlert.Annotations, alert.Annotations)
	assert.Equal(t, mockAlert.Labels, alert.Labels)
	assert.Equal(t, mockAlert.IsPaused, alert.IsPaused)
}
