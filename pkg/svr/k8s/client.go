package k8s

import (
	"bytes"
	"fmt"
	"io"
	"path"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// Client wrapper k8s client and namespace
type Client struct {
	kubernetes.Interface
	config *rest.Config
}

// NewClient create a new client
func NewClient(kubeconfig string) (*Client, error) {
	k8sConfig, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}
	config, err := k8sConfig.ClientConfig()
	if err != nil {
		fmt.Printf("create kubeconfig[%s], err[%s].\n", kubeconfig, err.Error())
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{Interface: cli, config: config}, err
}

// CopyFileToPod copy client file to pod
func (c *Client) CopyFileToPod(pod, container, namespace string, file io.Reader, dstPath string) error {
	dstDir := path.Dir(dstPath)
	execCmd := fmt.Sprintf("mkdir -p %s && cd %s && tar x", dstDir, dstDir)

	cmd := []string{
		"sh",
		"-c",
		execCmd,
	}
	fmt.Printf("exec command: %s\n", cmd)
	req := c.CoreV1().RESTClient().Post().
		Resource("pods").Name(pod).
		Namespace(namespace).SubResource("exec")

	req.VersionedParams(
		&v1.PodExecOptions{
			Command:   cmd,
			Container: container,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		},
		scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", req.URL())
	if err != nil {
		return err
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  file,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	fmt.Printf("copy file to pod[%s] container[%s], stdout[%s] stderr[%s]\n",
		pod, container, stdout.String(), stderr.String())
	if err != nil {
		fmt.Printf("copy file to pod[%s] container[%s] failed, err[%s].\n", pod, container, err.Error())
		return err
	}
	return nil
}
