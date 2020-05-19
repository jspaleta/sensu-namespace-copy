package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/api/core/v2"
	"gopkg.in/yaml.v2"
	//"strings"
	//"github.com/sensu/sensu-go/types"
	"io/ioutil"
	"log"
	"net/http"
	//"reflect"
	//	"os"
)

type Config struct {
	sensu.PluginConfig
	Namespace          string
	Resource           string
	ResourceType       string
	NewNamespace       string
	SensuApiUrl        string
	SensuAccessToken   string
	SensuTrustedCaFile string
	ResourceString     string
	Verbose            bool
	Dryrun             bool
	Output             bool
	Strip              bool
	Yaml               bool
}

var (
	config = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-namespace-copy",
			Short:    "Copy resources across namespaces",
			Keyspace: "sensu.io/plugins/sensu-namespace-copy/config",
		},
	}
	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "namespace",
			Env:       "SENSU_NAMESPACE", // provided by the sensuctl command plugin execution environment
			Argument:  "namespace",
			Shorthand: "n",
			Default:   "",
			Usage:     "Sensu namespace to copy from",
			Value:     &config.Namespace,
		},
		&sensu.PluginConfigOption{
			Path:      "sensu-api-url",
			Env:       "SENSU_API_URL", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-api-url",
			Shorthand: "",
			Default:   "",
			Usage:     "Sensu API URL (defaults to $SENSU_API_URL)",
			Value:     &config.SensuApiUrl,
		},
		&sensu.PluginConfigOption{
			Path:      "sensu-access-token",
			Env:       "SENSU_ACCESS_TOKEN", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-access-token",
			Shorthand: "",
			Default:   "",
			Usage:     "Sensu API Access Token (defaults to $SENSU_ACCESS_TOKEN)",
			Value:     &config.SensuAccessToken,
		},
		&sensu.PluginConfigOption{
			Path:      "sensu-trusted-ca-file",
			Env:       "SENSU_TRUSTED_CA_FILE", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-trusted-ca-file",
			Shorthand: "",
			Default:   "",
			Usage:     "Sensu API Trusted Certificate Authority File (defaults to $SENSU_TRUSTED_CA_FILE)",
			Value:     &config.SensuTrustedCaFile,
		},
		&sensu.PluginConfigOption{
			Path:      "new-namespace",
			Argument:  "new-namespace",
			Shorthand: "N",
			Default:   "",
			Usage:     "Sensu namespace to copy to",
			Value:     &config.NewNamespace,
		},
		&sensu.PluginConfigOption{
			Path:      "resource",
			Argument:  "resource",
			Shorthand: "r",
			Default:   "",
			Usage:     "Sensu resource to copy ",
			Value:     &config.Resource,
		},
		&sensu.PluginConfigOption{
			Path:      "resource-type",
			Argument:  "resource-type",
			Shorthand: "t",
			Default:   "",
			Usage:     "Sensu resource type to copy",
			Value:     &config.ResourceType,
		},
		&sensu.PluginConfigOption{
			Argument:  "verbose",
			Shorthand: "v",
			Default:   false,
			Usage:     "Enable verbose",
			Value:     &config.Verbose,
		},
		&sensu.PluginConfigOption{
			Argument:  "output",
			Shorthand: "o",
			Default:   false,
			Usage:     "Enable output to stdout",
			Value:     &config.Output,
		},
		&sensu.PluginConfigOption{
			Argument:  "strip",
			Shorthand: "s",
			Default:   false,
			Usage:     "Strip namespace, can only be used with output option",
			Value:     &config.Strip,
		},
		&sensu.PluginConfigOption{
			Argument:  "yaml",
			Shorthand: "y",
			Default:   false,
			Usage:     "Output yaml instead of json",
			Value:     &config.Yaml,
		},
		&sensu.PluginConfigOption{
			Argument:  "dryrun",
			Shorthand: "d",
			Default:   false,
			Usage:     "Perform no action, report configuration to stdout",
			Value:     &config.Dryrun,
		},
	}
)

func main() {
	/* Sensuctl commands are similar to Sensu Checks, so let's use the Sensu SDK's check object as a starting point */
	plugin := sensu.NewGoCheck(&config.PluginConfig, options, checkArgs, copyResource, false)
	/* Let's execute, checking the required arguments beforehand */
	plugin.Execute()
}

/* This is the goCheck executeFunction */
func copyResource(event *v2.Event) (int, error) {
	if config.Dryrun {
		fmt.Println("Dryrun selected, printing out configuration for review and exiting")
		fmt.Println("  SensuAccessToken:", config.SensuAccessToken)
		fmt.Println("  SensuApiUrl:", config.SensuApiUrl)
		fmt.Println("  Namespace:", config.Namespace)
		fmt.Println("  NewNamespace:", config.NewNamespace)
		fmt.Println("  Resource:", config.Resource)
		fmt.Println("  ResourceType:", config.ResourceType)
		fmt.Println("  Verbose:", config.Verbose)
		fmt.Println("  Strip:", config.Strip)
		fmt.Println("  Output:", config.Output)
		fmt.Println("  Dryrun:", config.Dryrun)
		return sensu.CheckStateOK, nil
	}
	resource, err := GetResource()
	metadata := resource["metadata"].(map[string]interface{})
	if config.Strip {
		delete(metadata, "namespace")
	}
	if len(config.NewNamespace) > 0 {
		metadata["namespace"] = config.NewNamespace
	}

	if config.Output {
		err = OutputResource(resource)
	} else {
		err = PostResource(resource)
	}
	if err != nil {
		return sensu.CheckStateCritical, err
	} else {
		return sensu.CheckStateOK, nil
	}
}

/* This used for the goCheck validateFunction */
func checkArgs(event *v2.Event) (int, error) {
	// Don't validate on dryrun
	if config.Dryrun {
		return sensu.CheckStateOK, nil
	}

	// basic validation
	if len(config.SensuApiUrl) == 0 {
		return sensu.CheckStateCritical, errors.New("--sensu-api-url flag or $SENSU_API_URL environment variable must be set")
	}
	if len(config.Namespace) == 0 {
		return sensu.CheckStateCritical, errors.New("--namespace flag or $SENSU_NAMESPACE environment variable must be set")
	}
	if len(config.ResourceType) == 0 {
		return sensu.CheckStateCritical, errors.New("--resource-type flag must be set")
	}
	if len(config.Resource) == 0 {
		return sensu.CheckStateCritical, errors.New("--resource flag must be set")
	}
	if len(config.SensuAccessToken) == 0 {
		return sensu.CheckStateCritical, errors.New("--sensu-access-token flag or $SENSU_ACCESS_TOKEN environment variable must be set")
	}
	if len(config.NewNamespace) == 0 && !config.Strip {
		return sensu.CheckStateCritical, errors.New("must either provide either --new-namespace or --strip option")
	}
	if len(config.NewNamespace) > 0 && config.Strip {
		return sensu.CheckStateCritical, errors.New("cannot use both --new-namespace and --strip options together")
	}
	if !config.Output && config.Strip {
		return sensu.CheckStateCritical, errors.New("--strip requires --output option")
	}
	return sensu.CheckStateOK, nil
}

func LoadCACerts(path string) (*x509.CertPool, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("ERROR: failed to load system cert pool: %s", err)
		return nil, err
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if path != "" {
		certs, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("ERROR: failed to read CA file (%s): %s", path, err)
			return nil, err
		} else {
			rootCAs.AppendCertsFromPEM(certs)
		}
	}
	return rootCAs, nil
}

func initHttpClient() *http.Client {
	certs, err := LoadCACerts(config.SensuTrustedCaFile)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	tlsConfig := &tls.Config{
		RootCAs: certs,
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: tr,
	}
	return client
}

func GetResource() (map[string]interface{}, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/core/v2/namespaces/%s/%s/%s",
			config.SensuApiUrl,
			config.Namespace,
			config.ResourceType,
			config.Resource,
		),
		nil,
	)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	var httpClient *http.Client = initHttpClient()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.SensuAccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return nil, err
	} else if resp.StatusCode == 404 {
		log.Fatalf("ERROR: %v %s (%s)\n", resp.StatusCode, http.StatusText(resp.StatusCode), req.URL)
		return nil, err
	} else if resp.StatusCode >= 300 {
		log.Fatalf("ERROR: %v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil, err
	}

	var result map[string]interface{}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return nil, err
	} else {

		err = json.Unmarshal([]byte(b), &result)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
func OutputResource(resource map[string]interface{}) error {
	if config.Yaml {
		body, err := yaml.Marshal(resource)
		if err != nil {
			log.Fatal("error: ", err)
		}
		fmt.Println("---\n" + string(body) + "...\n")
	} else {
		body, err := json.MarshalIndent(resource, "", "\t")
		if err != nil {
			log.Fatal("error: ", err)
		}
		fmt.Println(string(body))
	}
	return nil
}
func PostResource(resource map[string]interface{}) error {
	postBody, err := json.Marshal(resource)
	if err != nil {
		log.Fatal("error: ", err)
	}
	body := bytes.NewReader(postBody)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/core/v2/namespaces/%s/%s",
			config.SensuApiUrl,
			config.NewNamespace,
			config.ResourceType,
		),
		body,
	)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	var httpClient *http.Client = initHttpClient()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.SensuAccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	fmt.Println("POST Request Status", resp.StatusCode)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return err
	} else if resp.StatusCode == 404 {
		log.Fatalf("ERROR: %v %s (%s)\n", resp.StatusCode, http.StatusText(resp.StatusCode), req.URL)
		return err
	} else if resp.StatusCode == 409 {
		log.Fatalf("ERROR: %v %s Resource already exists in namespace \"%s\"\n", resp.StatusCode, http.StatusText(resp.StatusCode), config.NewNamespace)
		return err
	} else if resp.StatusCode >= 300 {
		log.Fatalf("ERROR: %v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		return err
	} else if resp.StatusCode == 201 {
		log.Printf("copied resource to namespace \"%s\"", config.NewNamespace)
		return nil
	} else {
	}
	return nil
}
