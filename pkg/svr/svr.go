package svr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/rtsien/k8scp/pkg/svr/k8s"
)

func UploadHandler(kubeconfig string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n------------\n")

		code := -1
		message := "success"
		defer func() { writeResponse(w, code, message) }()

		reader, err := r.MultipartReader()
		if err != nil {
			message = fmt.Sprintf("get a error: %v", err)
			return
		}

		params := make(map[string]string)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				message = fmt.Sprintf("reader NextPart error: %v", err)
				return
			}
			if part.FileName() == "" {
				switch part.FormName() {
				case "namespace", "pod", "container", "dst":
					r, err := io.ReadAll(part)
					if err != nil {
						message = fmt.Sprintf("read post parameter error: %v", err)
						return
					}
					params[part.FormName()] = string(r)
				default:
					_ = part.Close()
				}
			} else {
				k8sCli, err := k8s.NewClient(kubeconfig)
				if err != nil {
					message = fmt.Sprintf("get k8sCli error: %v", err)
					return
				}
				fmt.Printf("upload file [%s]\n", part.FileName())
				err = k8sCli.CopyFileToPod(params["pod"], params["container"], params["namespace"], part,
					fmt.Sprintf("%s/%s", params["dst"], path.Base(part.FileName())))
				if err != nil {
					message = fmt.Sprintf("upload file [%s] failed: %v", part.FileName(), err)
					return
				}
			}
		}
		code = 0
	}
}

func writeResponse(w http.ResponseWriter, code int, message string) {
	fmt.Println(message)
	resp, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})
	_, _ = w.Write(resp)
}
