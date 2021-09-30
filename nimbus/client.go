package nimbus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	httpClient *http.Client

	user string
}

func NewClient(user string) Client {
	return Client{
		httpClient: cleanhttp.DefaultPooledClient(),
		user:       user,
	}
}

func (c Client) ListTestbeds() ([]Testbed, error) {
	listCmd := []string{"nimbus-ctl", "--testbed", "list"}
	resp, err := c.launchNimbusCommand(context.Background(), &nimbusLaunchRequest{
		User: c.user,
		Cmd:  listCmd,
	})

	if err != nil {
		return nil, err
	}

	commandID := resp.ID
	logs, err := c.fetchCommandLogs(context.Background(), commandID)
	if err != nil {
		return nil, err
	}

	testbeds, err := c.extractTestbedsFromLogs(logs.Logs)
	if err != nil {
		return nil, err
	}

	return testbeds, nil
}

func (c Client) DeleteTestbed(name string) error {
	deleteCmd := []string{"nimbus-ctl", "--testbed", "kill", name}
	_, err := c.launchNimbusCommand(context.Background(), &nimbusLaunchRequest{
		User: c.user,
		Cmd:  deleteCmd,
	})

	if err != nil {
		return err
	}

	return nil
}

// for acceptance tests
func (c Client) CreateTestbed(name, jsonTestbedSpec, resultsDir string) error {
	createCmd := []string{
		"nimbus-testbeddeploy",
		"--testbedSpecJsonFile", jsonTestbedSpec,
		"--runName", name,
		"--resultsDir", resultsDir,
		"--waitForCompletion", "true"}
	_, err := c.launchNimbusCommand(context.Background(), &nimbusLaunchRequest{
		User: c.user,
		Cmd:  createCmd,
		Env: map[string]string{
			"USER": c.user,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create testbed %s: %w", name, err)
	}

	return nil
}

// launchNimbusCommand will execute a remote `nimbus` CLI command using the
// launcher API. It will wait for the operation to complete before returning.
// https://nimbus-api.eng.vmware.com/doc/index.html#GET%20/v1/launcher/:id/info-instance_method
func (c Client) launchNimbusCommand(ctx context.Context, req *nimbusLaunchRequest) (*nimbusLaunchResponse, error) {
	path := "v1/launcher"

	responseBody, err := c.invokeNimbusAPI(ctx, "POST", path, req, backoff.WithContext(&backoff.StopBackOff{}, ctx))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke Nimbus API %s: %w", path, err)
	}

	response := nimbusLaunchResponse{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Nimbus API response: err=%w, body=%s", err, string(responseBody))
	}

	if response.ID == "" || response.Error != "" {
		return nil, fmt.Errorf("Didn't receive valid nimbus task ID: %s %s", response.Error, response.Message)
	}
	return &response, c.waitForNimbusCommand(ctx, response.ID)
}

func (c Client) waitForNimbusCommand(ctx context.Context, commandID string) error {
	for {
		resp, err := c.fetchCommandStatus(ctx, commandID)
		if err != nil {
			return err
		}

		switch resp.Status {
		case "succeeded":
			return nil
		case "failed":
			logs, err := c.fetchCommandLogs(ctx, commandID)
			if err != nil {
				logs = &nimbusLogsResponse{}
			}

			return fmt.Errorf("Nimbus command failed %+v (id %v), logs:\n%s", resp, commandID, logs.Logs)
		default:
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

// https://nimbus-api.eng.vmware.com/doc/index.html#GET%20/v1/launcher/:id/status-instance_method
func (c Client) fetchCommandStatus(ctx context.Context, commandID string) (*nimbusStatusResponse, error) {
	var response nimbusStatusResponse
	path := fmt.Sprintf("v1/launcher/%s/status", commandID)
	responseBody, err := c.invokeNimbusAPI(ctx, "GET", path, nil, backoff.WithContext(backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10), ctx))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke Nimbus API %s: %v", path, err)
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Nimbus API response: err=%v body=%s", err, string(responseBody))
	}
	return &response, nil
}

func (c Client) fetchCommandLogs(ctx context.Context, commandID string) (*nimbusLogsResponse, error) {
	path := fmt.Sprintf("v1/launcher/%s/logs", commandID)
	responseBody, err := c.invokeNimbusAPI(ctx, "GET", path, map[string]string{}, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke Nimbus API %s: %w", path, err)
	}
	var response nimbusLogsResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Nimbus API response: err=%w body=%s", err, string(responseBody))
	}
	return &response, nil
}

func (c Client) extractTestbedsFromLogs(logs string) ([]Testbed, error) {
	pattern := regexp.MustCompile(`.*: (.*)\(ID:.*`)

	matches := pattern.FindAllStringSubmatch(logs, -1)
	if matches == nil {
		return nil, fmt.Errorf("fail to match pattern from nimbus log, pattern = %s, log = %s", pattern, logs)
	}

	testbeds := []Testbed{}
	for _, match := range matches {
		for i, testbed := range match {
			if i == 0 {
				continue // 0 index is the whole match
			}
			testbeds = append(testbeds, Testbed{
				client: c,
				name:   testbed,
			})
		}
	}

	return testbeds, nil
}

func (c Client) invokeNimbusAPI(ctx context.Context, method, path string, body interface{}, b backoff.BackOffContext) ([]byte, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	var out []byte
	err = backoff.Retry(func() error {
		req, err := http.NewRequestWithContext(ctx, method,
			fmt.Sprintf("https://nimbus-api.eng.vmware.com/api/%s", path),
			bytes.NewBuffer(bodyBytes))
		if err != nil {
			return err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		outBody := string(out)
		if strings.Contains(outBody, "408 Request Time-out") ||
			strings.Contains(outBody, "Limit of 500 reached") {
			return errors.New("Body contains error response")
		}

		return nil
	}, b)

	return out, err
}
