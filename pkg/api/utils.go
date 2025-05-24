package api

import (
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic"
	"hashTest/api/models"
	"io"
	"net/http"
	"strings"
)

// unmarshal is a utility method to parse json body with sonic
func (c *Client) unmarshal(body io.ReadCloser, val interface{}) error {
	gzReader, err := gzip.NewReader(body)
	if err != nil {
		return err
	}
	defer func(gzReader *gzip.Reader) {
		err := gzReader.Close()
		if err != nil {
			panic(err)
		}
	}(gzReader)

	bytes, err := io.ReadAll(gzReader)
	if err != nil {
		return err
	}

	if err := sonic.Unmarshal(bytes, &val); err != nil {
		return err
	}

	return nil
}

// applyHeaders is a utility method to apply request headers, that bypass cloudflare, add auth and etc
func (c *Client) applyHeaders(req *http.Request, bodyLen int) {
	req.Header.Set("User-Agent", "DeepSeek/1.2.1 Android/30")
	req.Header.Set("Content-Type", "application/json")
	if bodyLen > 0 {
		req.Header.Set("Content-Length", fmt.Sprintf("%d", bodyLen))
	}
	if c.ApiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Charset", "UTF-8")
	req.Header.Set("X-Rangers-Id", c.rangersId)
	req.Header.Set("X-Client-Locale", "en")
	req.Header.Set("X-Client-Version", "1.2.1")
	req.Header.Set("X-Client-Platform", "android")
}

// applyPowHeader is a utility method to apply request header for Proof-of-Work result to request. IT NEEDS TO BE APPLIED
func (c *Client) applyPowHeader(req *http.Request, answer int, pow models.PowChallenge) {
	header := fmt.Sprintf(`{"algorithm":"%s","challenge":"%s","salt":"%s","signature":"%s","answer":%d,"target_path":"%s"}`, pow.Algorithm, pow.Challenge, pow.Salt, pow.Signature, answer, pow.TargetPath)
	encodedHeader := base64.StdEncoding.EncodeToString([]byte(header))
	req.Header.Add("X-Ds-Pow-Response", encodedHeader)
}

// execute is a utility method to send a request
func (c *Client) execute(url string, body string, method string, val interface{}) error {
	var req *http.Request
	var err error

	if body == "" || method == http.MethodGet {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	}

	if err != nil {
		return err
	}

	c.applyHeaders(req, len(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if err := c.unmarshal(resp.Body, val); err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

// base sse struct, gets received from completion
type event struct {
	V interface{} `json:"v"`
	P string      `json:"p"`
}

func parseEvents(r io.ReadCloser, tokensCh chan<- string) {
	defer close(tokensCh)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "event: ") {
			// log.Println(line)
			continue
		}
		raw := strings.TrimPrefix(line, "data: ")
		var ev event
		if err := sonic.Unmarshal([]byte(raw), &ev); err != nil {
			continue
		}
		if ev.P != "" {
			switch ev.P {
			case "response/search_status":
				tokensCh <- "\n<searching>\n"
				continue
			case "response/thinking_content":
				tokensCh <- "<searching/>\n" +
					"\n<thinking>\n"
			case "response/content":
				tokensCh <- "\n</thinking>\n" + "\n<answer>\n"
			case "response/status":
				tokensCh <- "\n</answer>"
				continue
			}
		}
		switch v := ev.V.(type) {
		case string:
			tokensCh <- v
		case []interface{}:
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					title := getString(m, "title")
					url := getString(m, "url")
					snippet := getString(m, "snippet")
					var link string

					if title != "" && url != "" {
						if snippet == "" {
							link = fmt.Sprintf("[%s](%s)\n",
								title,
								url)
						} else {
							link = fmt.Sprintf("[%s](%s) - %s\n",
								title,
								url,
								snippet)
						}
					} else {
						continue
					}

					tokensCh <- link
				}
				if s, ok := item.(string); ok {
					tokensCh <- s
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic("error reading input: " + err.Error())
	}
	err := r.Close()
	if err != nil {
		panic(err)
	}
}

func getString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
