package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	var (
		aggHTTPEndpoint  = "http://127.0.0.1" + os.Getenv("AGG_HTTP_ENDPOINT")
		gateHTTPEndpoint = os.Getenv("GATE_HTTP_ENDPOINT")
		client           = client.NewHTTPClient(aggHTTPEndpoint)
		invoiceHandler   = newInvoiceHandler(client)
	)

	http.HandleFunc("/invoice", makeApiFunc(invoiceHandler.HandleGetInvoice))
	logrus.Infof("Starting HTTP server on port %s", gateHTTPEndpoint)
	log.Fatal(http.ListenAndServe(gateHTTPEndpoint, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	vals, ok := r.URL.Query()["obu"]
	if !ok {
		return fmt.Errorf("invalid id")
	}
	if len(vals) == 0 {
		return fmt.Errorf("invalid id")
	}
	id, err := strconv.Atoi(vals[0])
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	inv, err := h.client.GetInvoice(context.Background(), id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"err": err.Error()})
		}
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
