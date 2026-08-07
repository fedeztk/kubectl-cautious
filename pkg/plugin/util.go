package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

func getContext(kubeconfig string) (string, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{
		Precedence: splitKubeconfigPath(kubeconfig),
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()

	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

func getKubeconfigPath() (string, error) {
	if os.Getenv("KUBECONFIG") != "" {
		return os.Getenv("KUBECONFIG"), nil
	}

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		return "", fmt.Errorf("kubeconfig file not found at %s", kubeconfig)
	}
	return kubeconfig, nil
}

// splitKubeconfigPath splits a KUBECONFIG value by the OS path list separator
// (colon on Linux/macOS, semicolon on Windows).
func splitKubeconfigPath(kubeconfig string) []string {
	sep := string(os.PathListSeparator)
	parts := strings.Split(kubeconfig, sep)
	var paths []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}
