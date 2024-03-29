package healthz

import (
	"bytes"
	"fmt"
	"net/http"
)

func handleRootHealthz(fetcher ChecksFetcher) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		failed := false
		var verboseOut bytes.Buffer
		for _, check := range fetcher() {
			if err := check.Check(); err != nil {
				fmt.Fprintf(&verboseOut, "[-]%v failed\n", check.Name())
				failed = true
			} else {
				fmt.Fprintf(&verboseOut, "[+]%v ok\n", check.Name())
			}
		}

		res.Header().Set("Content-Type", "text/plain")
		if failed {
			http.Error(res, fmt.Sprintf("%vhealthz check failed", verboseOut.String()),
				http.StatusInternalServerError)
			return
		}

		if _, found := req.URL.Query()["verbose"]; !found {
			fmt.Fprint(res, "ok")
			return
		}

		verboseOut.WriteTo(res)
		fmt.Fprint(res, "healthz check passed\n")
	})
}

func (h *Healthz) InstallHandler(mux *http.ServeMux) {
	mux.Handle("/healthz", handleRootHealthz(h.Checks))
}
