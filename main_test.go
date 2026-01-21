package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []string{
		"/cafe?city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range request {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: 5},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		requset := "/cafe?city=moscow&count=" + strconv.Itoa(v.count)
		req := httptest.NewRequest("GET", requset, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())
		if body == "" {
			assert.Equal(t, 0, v.want)
		} else {
			cafes := strings.Split(body, ",")
			assert.Len(t, cafes, v.want)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		resp := httptest.NewRecorder()
		str := "/cafe?city=moscow&search=" + v.search
		req := httptest.NewRequest("GET", str, nil)
		handler.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)

		s := strings.ToLower(v.search)
		body := resp.Body.String()
		body = strings.TrimSpace(body)
		if body == "" {
			assert.Equal(t, 0, v.wantCount)
			continue
		}
		body = strings.ToLower(body)

		sliceCafe := strings.Split(body, ",")
		for _, j := range sliceCafe {
			assert.Contains(t, strings.ToLower(j), s)
		}
		assert.Len(t, sliceCafe, v.wantCount)
	}
}
