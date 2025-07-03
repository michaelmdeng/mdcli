package rm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	remarkableHost = "10.11.99.1"
	documentsPath  = "/documents/"
	downloadsPath  = "/download/%s/placeholder"

	DocumentType = "DocumentType"
)

func getDocuments() ([]Document, error) {
	u := url.URL{
		Scheme: "http",
		Host:   remarkableHost,
		Path:   documentsPath,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var docs []RMDocument
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		return nil, err
	}

	output := make([]Document, 0, len(docs))
	for _, doc := range docs {
		if doc.Type != DocumentType {
			continue
		}

		modifiedTime, err := time.Parse(time.RFC3339, doc.ModifiedClient)
		if err != nil {
			return nil, err
		}

		output = append(output, Document{
			ID:           doc.ID,
			ModifiedTime: modifiedTime,
			Name:         doc.VisibleName,
		})
	}

	return output, nil
}

func getDocumentIdByName(name string) (string, error) {
	u := url.URL{
		Scheme: "http",
		Host:   remarkableHost,
		Path:   documentsPath,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var docs []RMDocument
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		return "", err
	}

	for _, doc := range docs {
		if doc.VisibleName == name {
			return doc.ID, nil
		}
	}

	return "", errors.New("document not found")
}

func downloadDocument(id string, outputFile string) error {
	u := url.URL{
		Scheme: "http",
		Host:   remarkableHost,
		Path:   fmt.Sprintf(downloadsPath, id),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if _, err := os.Stat(outputFile); err == nil {
		return errors.New("file already exists")
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
