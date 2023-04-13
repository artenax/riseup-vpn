package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"0xacab.org/leap/bitmask-vpn/pkg/bitmask"
)

var vpnTransports = []string{"openvpn", "obfs4", "obfs4-hop", "obfs4-kcp"}

func CheckAuth(handler http.HandlerFunc, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("X-Auth-Token")
		if t == token {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
		}
	}
}

func webOn(w http.ResponseWriter, r *http.Request) {
	log.Println("Web UI: on")
	SwitchOn()
}

func webOff(w http.ResponseWriter, r *http.Request) {
	log.Println("Web UI: off")
	SwitchOff()
}

func webStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ctx.Status.String())
}

func webGatewayGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ctx.bm.GetCurrentGateway())
}

func webGatewaySet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		gwLabel := r.FormValue("transport")
		fmt.Fprintf(w, "selected gateway: %s\n", gwLabel)
		ctx.bm.UseGateway(gwLabel)
		// TODO make sure we don't tear the fw down on reconnect...
		SwitchOff()
		// a little sleep is needed, though, because iptables takes some time
		time.Sleep(500 * time.Millisecond)
		SwitchOn()
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Only POST supported.")
	}
}

func webGatewayList(w http.ResponseWriter, r *http.Request) {
	transport := ctx.bm.GetTransport()
	locationJson, err := json.Marshal(ctx.bm.ListLocationFullness(transport))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting json: %v", err)
	}
	fmt.Fprintf(w, string(locationJson))
}

func webTransportGet(w http.ResponseWriter, r *http.Request) {
	t, err := json.Marshal(ctx.bm.GetTransport())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting json: %v", err)
	}
	fmt.Fprintf(w, string(t))

}

func webTransportSet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		t := r.FormValue("transport")
		if isValidTransport(t) {
			fmt.Fprintf(w, "Selected transport: %s\n", t)
			go ctx.bm.SetTransport(string(t))
		} else {
			fmt.Fprintf(w, "Unknown transport: %s\n", t)
		}
	default:
		fmt.Fprintf(w, "Only POST supported.")
	}
}

func webTransportList(w http.ResponseWriter, r *http.Request) {
	trList, err := json.Marshal(ctx.bm.ListTransport())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting json: %v", err)
	}
	fmt.Fprintf(w, string(trList))
}

func webQuit(w http.ResponseWriter, r *http.Request) {
	log.Println("Web UI: quit")
	Quit()
	os.Exit(0)
}

func enableWebAPI(port int) {
	log.Println("Starting WebAPI in port", port)
	bitmask.GenerateAuthToken()
	token := bitmask.ReadAuthToken()

	handle := func(pattern string, fn func(w http.ResponseWriter, r *http.Request)) {
		http.Handle(pattern, CheckAuth(http.HandlerFunc(fn), token))
	}

	handle("/vpn/start", webOn)
	handle("/vpn/stop", webOff)
	handle("/vpn/gw/get", webGatewayGet)
	handle("/vpn/gw/set", webGatewaySet)
	handle("/vpn/gw/list", webGatewayList)
	handle("/vpn/transport/get", webTransportGet)
	handle("/vpn/transport/set", webTransportSet)
	handle("/vpn/transport/list", webTransportList)
	handle("/vpn/status", webStatus)
	handle("/vpn/quit", webQuit)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func isValidTransport(t string) bool {
	for _, b := range vpnTransports {
		if b == t {
			return true
		}
	}
	return false
}
