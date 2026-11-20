package admin

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"gin-spa-app-framework-back/common"
	"gin-spa-app-framework-back/initialize"
	"gin-spa-app-framework-back/internal/routes"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	initialize.InitSystem()

	r = routes.Router
}

func TestCreateUser(t *testing.T) {
	postData := url.Values{
		"username":         {"user1"},
		"password":         {"123456"},
		"confirm_password": {"123456"},
	}

	reqBody := strings.NewReader(postData.Encode())
	req, err := http.NewRequest("POST", "/admin/user/create", reqBody)
	if err != nil {
		t.Fatal("http.NewRequest error:", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rdb := common.COM_REDIS
	if rdb == nil {
		t.Error("rdb is nil")
	}
	ctx := context.Background()
	token, err := rdb.Get(ctx, "x-header-token").Result()
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("x-header-token", token)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	res := rec.Result()
	resBody, _ := io.ReadAll(res.Body)
	resBodyString := string(resBody)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status code error,want %d,got %d, msg: %v", http.StatusOK, res.StatusCode, resBodyString)
	}

	defer res.Body.Close()
}

func TestGetUserInfoByToken(t *testing.T) {
	rdb := common.COM_REDIS
	if rdb == nil {
		t.Error("rdb is nil")
	}
	ctx := context.Background()
	token, err := rdb.Get(ctx, "x-header-token").Result()
	if err != nil {
		t.Error(err)
	}
	postData := url.Values{
		"data": {token},
		"type": {"token"},
	}
	reqBody := strings.NewReader(postData.Encode())
	req, err := http.NewRequest("POST", "/public/admin/user/get_user_info", reqBody)
	if err != nil {
		t.Fatal("http.NewRequest error:", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	res := rec.Result()
	resBody, _ := io.ReadAll(res.Body)
	resBodyString := string(resBody)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status code error,want %d,got %d, msg: %v", http.StatusOK, res.StatusCode, resBodyString)
	}
	println("resBody:", resBodyString)

	defer res.Body.Close()
}
