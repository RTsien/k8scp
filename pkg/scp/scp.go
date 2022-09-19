package scp

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/rtsien/k8scp/pkg/common"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Copy struct {
	Src string

	ServerURL string
	Namespace string
	Pod       string
	Container string
	Dst       string
}

func (c *Copy) Do() error {
	pR, pW := io.Pipe()
	defer func() { _ = pW.Close() }()

	bodyWriter := multipart.NewWriter(pW)
	contentType := bodyWriter.FormDataContentType()

	go func() {
		_ = bodyWriter.WriteField("namespace", c.Namespace)
		_ = bodyWriter.WriteField("pod", c.Pod)
		_ = bodyWriter.WriteField("container", c.Container)
		_ = bodyWriter.WriteField("dst", c.Dst)
		fileWriter, err := bodyWriter.CreateFormFile("file", c.Src)
		if err != nil {
			fmt.Printf("writing to buffer failed, err: %s", err.Error())
			os.Exit(1)
		}
		tarWriter := tar.NewWriter(fileWriter)

		srcReader, err := os.Open(c.Src)
		common.AssertErr(err, "open src file %s error", c.Src)
		fileInfo, err := srcReader.Stat()
		common.AssertErr(err, "get src file %s info error", c.Src)
		hdr, err := tar.FileInfoHeader(fileInfo, "")

		err = tarWriter.WriteHeader(hdr)
		common.AssertErr(err, "write tar header info of file %s error", c.Src)
		_, err = io.Copy(tarWriter, srcReader)
		if err != nil {
			fmt.Printf("io copy failed, err: %s", err.Error())
			os.Exit(1)
		}
		_ = tarWriter.Close()
		_ = bodyWriter.Close()
		_ = pW.Close()
	}()

	resp, err := http.Post(c.ServerURL, contentType, pR)
	if err != nil {
		fmt.Printf("\nupload failed, err: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\nread response failed, err: %s\n", err.Error())
		os.Exit(1)
	}
	var respJson struct {
		Code    int
		Message string
	}
	err = json.Unmarshal(respBody, &respJson)
	common.AssertErr(err, "Response: %s\n", string(respBody))
	if respJson.Code != 0 {
		fmt.Printf("\nupload failed, code: %d, errMsg: %s\n", respJson.Code, respJson.Message)
		os.Exit(1)
	}

	return nil
}
