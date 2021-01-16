package handler

import (
	"encoding/json"
	"github.com/zale144/ports/clientapi/internal/api"
	"github.com/zale144/ports/clientapi/pkg/apierror"
	"io"
	"net/http"
)

func PortJSON(up api.PortService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var lastRet apierror.ErrorMessage
		status := http.StatusOK
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			lastRet = up.Process(r.Context(), part)
			if lastRet.R.HasErrors() {
				status = http.StatusInternalServerError
				break
			}
		}

		jsn, _ := json.Marshal(&lastRet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(jsn)
	}
}

func GetPorts(svc api.PortService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		var (
			resp   interface{}
			errRsp apierror.ErrorMessage
		)

		resp, errRsp = svc.GetPorts(r.Context())
		if errRsp.R.HasErrors() {
			status = http.StatusInternalServerError
			resp = errRsp
		}
		jsn, _ := json.Marshal(&resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(jsn)
	}
}
