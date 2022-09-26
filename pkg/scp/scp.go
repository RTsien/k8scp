package scp

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/juju/ratelimit"
	progressbar "github.com/schollz/progressbar/v3"

	"github.com/rtsien/k8scp/pkg/common"
)

type Copy struct {
	Src       string
	ServerURL string
	Namespace string
	Pod       string
	Container string
	Dst       string
}

const rate = 12 * 1024 * 1024 // byte/s

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
			fmt.Fprintf(os.Stderr, "writing to buffer failed, err: %s", err.Error())
			os.Exit(1)
		}
		tarWriter := tar.NewWriter(fileWriter)

		err = filepath.Walk(c.Src, func(file string, fi os.FileInfo, err error) error {
			relFile, err := filepath.Rel(c.Src, file)
			common.AssertErr(err, "get file %s relative path error", file)
			// if src is a file path
			if relFile == "." {
				relFile = filepath.Base(c.Src)
			}

			hdr, err := tar.FileInfoHeader(fi, relFile)
			common.AssertErr(err, "get src file %s info error", c.Src)

			// modify os.PathSeparator to slash
			hdr.Name = filepath.ToSlash(relFile)
			err = tarWriter.WriteHeader(hdr)
			common.AssertErr(err, "write tar header info of file %s error", c.Src)

			if !fi.IsDir() {
				bar := progressbar.DefaultBytes(
					fi.Size(),
					relFile,
				)
				progressbar.OptionUseANSICodes(true)(bar)

				srcReader, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err = io.Copy(io.MultiWriter(tarWriter, bar),
					ratelimit.Reader(srcReader, ratelimit.NewBucketWithRate(rate, rate))); err != nil {
					fmt.Fprintf(os.Stderr, "io copy failed, err: %s", err.Error())
					os.Exit(1)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("writing to buffer failed, err: %s", err.Error())
			os.Exit(1)
		}

		_ = tarWriter.Close()
		_ = bodyWriter.Close()
		_ = pW.Close()
	}()

	resp, err := http.Post(c.ServerURL, contentType, pR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nupload failed, err: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nread response failed, err: %s\n", err.Error())
		os.Exit(1)
	}
	var respJson struct {
		Code    int
		Message string
	}
	err = json.Unmarshal(respBody, &respJson)
	common.AssertErr(err, "Response: %s\n", string(respBody))
	if respJson.Code != 0 {
		fmt.Fprintf(os.Stderr, "\nupload failed, code: %d, errMsg: %s\n", respJson.Code, respJson.Message)
		os.Exit(1)
	}

	return nil
}
