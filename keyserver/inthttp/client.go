// Copyright 2020 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inthttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/matrix-org/dendrite/internal/httputil"
	"github.com/matrix-org/dendrite/keyserver/api"
	userapi "github.com/matrix-org/dendrite/userapi/api"
)

// HTTP paths for the internal HTTP APIs
const (
	InputDeviceListUpdatePath         = "/keyserver/inputDeviceListUpdate"
	PerformUploadKeysPath             = "/keyserver/performUploadKeys"             //执行上传秘钥
	PerformClaimKeysPath              = "/keyserver/performClaimKeys"              //执行领取秘钥
	PerformDeleteKeysPath             = "/keyserver/performDeleteKeys"             //执行删除秘钥
	PerformUploadDeviceKeysPath       = "/keyserver/performUploadDeviceKeys"       //执行上传设备key
	PerformUploadDeviceSignaturesPath = "/keyserver/performUploadDeviceSignatures" //执行上传设备消息
	QueryKeysPath                     = "/keyserver/queryKeys"                     //查询key
	QueryKeyChangesPath               = "/keyserver/queryKeyChanges"               //查询key变化
	QueryOneTimeKeysPath              = "/keyserver/queryOneTimeKeys"              //查询一次性key
	QueryDeviceMessagesPath           = "/keyserver/queryDeviceMessages"           //查询设备消息
	QuerySignaturesPath               = "/keyserver/querySignatures"               //查询签名
	PerformMarkAsStalePath            = "/keyserver/markAsStale"                   //标记为过时
)

// NewKeyServerClient creates a KeyInternalAPI implemented by talking to a HTTP POST API.
// If httpClient is nil an error is returned
func NewKeyServerClient(
	apiURL string,
	httpClient *http.Client,
) (api.KeyInternalAPI, error) {
	if httpClient == nil {
		return nil, errors.New("NewKeyServerClient: httpClient is <nil>")
	}
	return &httpKeyInternalAPI{
		apiURL:     apiURL,
		httpClient: httpClient,
	}, nil
}

type httpKeyInternalAPI struct {
	apiURL     string
	httpClient *http.Client
}

func (h *httpKeyInternalAPI) SetUserAPI(i userapi.KeyserverUserAPI) {
	// no-op: doesn't need it
}

func (h *httpKeyInternalAPI) PerformClaimKeys(
	ctx context.Context,
	request *api.PerformClaimKeysRequest,
	response *api.PerformClaimKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"PerformClaimKeys", h.apiURL+PerformClaimKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) PerformDeleteKeys(
	ctx context.Context,
	request *api.PerformDeleteKeysRequest,
	response *api.PerformDeleteKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"PerformDeleteKeys", h.apiURL+PerformDeleteKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) PerformUploadKeys(
	ctx context.Context,
	request *api.PerformUploadKeysRequest,
	response *api.PerformUploadKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"PerformUploadKeys", h.apiURL+PerformUploadKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) QueryKeys(
	ctx context.Context,
	request *api.QueryKeysRequest,
	response *api.QueryKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"QueryKeys", h.apiURL+QueryKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) QueryOneTimeKeys(
	ctx context.Context,
	request *api.QueryOneTimeKeysRequest,
	response *api.QueryOneTimeKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"QueryOneTimeKeys", h.apiURL+QueryOneTimeKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) QueryDeviceMessages(
	ctx context.Context,
	request *api.QueryDeviceMessagesRequest,
	response *api.QueryDeviceMessagesResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"QueryDeviceMessages", h.apiURL+QueryDeviceMessagesPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) QueryKeyChanges(
	ctx context.Context,
	request *api.QueryKeyChangesRequest,
	response *api.QueryKeyChangesResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"QueryKeyChanges", h.apiURL+QueryKeyChangesPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) PerformUploadDeviceKeys(
	ctx context.Context,
	request *api.PerformUploadDeviceKeysRequest,
	response *api.PerformUploadDeviceKeysResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"PerformUploadDeviceKeys", h.apiURL+PerformUploadDeviceKeysPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) PerformUploadDeviceSignatures(
	ctx context.Context,
	request *api.PerformUploadDeviceSignaturesRequest,
	response *api.PerformUploadDeviceSignaturesResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"PerformUploadDeviceSignatures", h.apiURL+PerformUploadDeviceSignaturesPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) QuerySignatures(
	ctx context.Context,
	request *api.QuerySignaturesRequest,
	response *api.QuerySignaturesResponse,
) error {
	return httputil.CallInternalRPCAPI(
		"QuerySignatures", h.apiURL+QuerySignaturesPath,
		h.httpClient, ctx, request, response,
	)
}

func (h *httpKeyInternalAPI) PerformMarkAsStaleIfNeeded(
	ctx context.Context,
	request *api.PerformMarkAsStaleRequest,
	response *struct{},
) error {
	return httputil.CallInternalRPCAPI(
		"MarkAsStale", h.apiURL+PerformMarkAsStalePath,
		h.httpClient, ctx, request, response,
	)
}
