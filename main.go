package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tomwright/dasel/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

const (
	KUBECONFIG_DIR        = ".kube/config"
	CONTEXT_SELECTOR_YAML = ".contexts.[0].context.namespace"
)

func main() {
	nArgs := len(os.Args)

	switch nArgs {
	case 1:
		homePath, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		kubeconfigPath := homePath + "/" + KUBECONFIG_DIR
		currentNs, err := GetYamlField(kubeconfigPath, CONTEXT_SELECTOR_YAML)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		fmt.Println(currentNs)
	case 2:
		targetNamespace := os.Args[1]
		fmt.Printf("switching to namespace %q\n", targetNamespace)

		homePath, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		kubeconfigPath := homePath + "/" + KUBECONFIG_DIR

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
		kubeClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		myNs, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), targetNamespace, metav1.GetOptions{})
		if err != nil {
			log.Printf("error: %s\n", err)
		}

		err = UpdateYamlField(kubeconfigPath, CONTEXT_SELECTOR_YAML, myNs.Name)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}
	default:
		log.Fatalf("error: incorrect number of arguments (0 or 1!)\n")
	}

}

// UpdateYamlField reads a yaml file and returns one field
func GetYamlField(file, field string) (string, error) {
	// open the config file and transform to json (so we can modify it with dasel)
	yFile, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	jFile, err := yaml.YAMLToJSON(yFile)
	if err != nil {
		return "", err
	}

	// manipulate the data with dasel library
	var data interface{}
	err = json.Unmarshal(jFile, &data)
	if err != nil {
		return "", err
	}

	result, err := dasel.Select(data, field)
	if err != nil {
		return "", err
	}
	currentNs := result[0].String()
	return currentNs, nil
}

// UpdateYamlField reads a yaml file and modify one field with the desired value
func UpdateYamlField(file, field, value string) error {
	// open the config file and transform to json (so we can modify it with dasel)
	yFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	jFile, err := yaml.YAMLToJSON(yFile)
	if err != nil {
		return err
	}

	// manipulate the data with dasel library
	var data interface{}
	err = json.Unmarshal(jFile, &data)
	if err != nil {
		return err
	}

	result, err := dasel.Put(data, field, value)
	if err != nil {
		return err
	}

	// return the file to json, convert it to yaml so metricbeat can use it, and write it to his destination
	jOut, _ := json.Marshal(result.Interface())
	yOut, err := yaml.JSONToYAML(jOut)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, yOut, 0)
	if err != nil {
		return err
	}
	return nil
}
