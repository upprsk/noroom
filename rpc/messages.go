package rpc

import (
	"encoding/json"
	"errors"
	"time"
)

type RpcRequest struct {
	Method string
	Params json.RawMessage
}

type RpcIdRequestParams struct {
	Id string
}

type RpcIdTimeoutRequestParams struct {
	Id      string
	Timeout time.Duration
}

type RpcCreateRequestParams struct {
	Name  string
	Image string
	Cmd   []string
	Env   map[string]string
}

type RpcStartRequestParams = RpcIdTimeoutRequestParams
type RpcStopRequestParams = RpcIdTimeoutRequestParams
type RpcKillRequestParams = RpcIdTimeoutRequestParams
type RpcDeleteRequestParams = RpcIdRequestParams
type RpcInspectRequestParams = RpcIdRequestParams
type RpcAttachRequestParams = RpcIdRequestParams

func NewRpcCreateRequest(params RpcCreateRequestParams) (RpcRequest, error) {
	return NewRpcRequest("create", params)
}

func NewRpcStartRequest(params RpcStartRequestParams) (RpcRequest, error) {
	return NewRpcRequest("start", params)
}

func NewRpcStopRequest(params RpcStopRequestParams) (RpcRequest, error) {
	return NewRpcRequest("stop", params)
}

func NewRpcKillRequest(params RpcKillRequestParams) (RpcRequest, error) {
	return NewRpcRequest("kill", params)
}

func NewRpcDeleteRequest(params RpcDeleteRequestParams) (RpcRequest, error) {
	return NewRpcRequest("delete", params)
}

func NewRpcInspectRequest(params RpcInspectRequestParams) (RpcRequest, error) {
	return NewRpcRequest("inspect", params)
}

func NewRpcAttachRequest(params RpcAttachRequestParams) (RpcRequest, error) {
	return NewRpcRequest("attach", params)
}

func NewRpcRequest(method string, params any) (RpcRequest, error) {
	rawParams, err := json.Marshal(params)
	if err != nil {
		return RpcRequest{}, err
	}

	return RpcRequest{Method: method, Params: rawParams}, nil
}

func NewRpcError(err error) RpcBaseResponse {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	return RpcBaseResponse{Err: errStr}
}

type RpcBaseResponse struct {
	Err string
}

func (r RpcBaseResponse) GetErr() error {
	if r.Err == "" {
		return nil
	}

	return errors.New(r.Err)
}

type RpcIdResponse struct {
	RpcBaseResponse
	Id string
}

type RpcInspectResponse struct {
	RpcBaseResponse
	Data *ContainerInspectResult
}

type RpcEmptyResponse = RpcBaseResponse
type RpcCreateResponse = RpcIdResponse
type RpcKillResponse = RpcIdResponse
