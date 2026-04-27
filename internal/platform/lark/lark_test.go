package lark

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(srv *httptest.Server) *Client {
	return &Client{appID: "id", appSecret: "secret", baseURL: srv.URL, httpClient: srv.Client()}
}

func TestSendMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/open-apis/auth/v3/tenant_access_token/internal":
			json.NewEncoder(w).Encode(map[string]any{"code": 0, "tenant_access_token": "tok"})
		default:
			json.NewEncoder(w).Encode(map[string]any{"code": 0})
		}
	}))
	defer srv.Close()

	c := newTestClient(srv)
	if err := c.SendMessage("user1", "hello"); err != nil {
		t.Fatal(err)
	}
}

func TestSendMessageAuthFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"code": 99})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	if err := c.SendMessage("user1", "hello"); err == nil {
		t.Fatal("expected error")
	}
}

func TestReadMessages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/open-apis/auth/v3/tenant_access_token/internal":
			json.NewEncoder(w).Encode(map[string]any{"code": 0, "tenant_access_token": "tok"})
		default:
			json.NewEncoder(w).Encode(map[string]any{
				"code": 0,
				"data": map[string]any{
					"items": []map[string]any{
						{"message_id": "m1", "sender": map[string]any{"id": "u1"}, "body": map[string]any{"content": "hi"}},
					},
				},
			})
		}
	}))
	defer srv.Close()

	c := newTestClient(srv)
	msgs, err := c.ReadMessages("chat1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].ID != "m1" {
		t.Fatalf("unexpected messages: %+v", msgs)
	}
}
