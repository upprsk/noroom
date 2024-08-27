package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ContainerState struct {
	Status     string
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  string
	FinishedAt string
}

type ContainerInspectResult struct {
	Id         string
	Name       string
	Path       string
	Args       []string
	Image      string
	Created    string
	SizeRw     *int64
	SizeRootFs *int64
	State      ContainerState
}

type RpcHandler interface {
	Create(ctx context.Context, name, image string, cmd []string, env []string) (string, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Kill(ctx context.Context, id, signal string) error
	Delete(ctx context.Context, id string) error
	Inspect(ctx context.Context, id string) (*ContainerInspectResult, error)
}

type RpcServer struct {
	stream  io.ReadWriter
	handler RpcHandler

	timeout time.Duration
}

func NewRpcServer(stream io.ReadWriter, timeout time.Duration, handler RpcHandler) *RpcServer {
	return &RpcServer{
		stream:  stream,
		handler: handler,
		timeout: timeout,
	}
}

func (rpc *RpcServer) HandleOne(ctx context.Context) (bool, error) {
	var req RpcRequest
	if err := json.NewDecoder(rpc.stream).Decode(&req); err != nil {
		return false, err
	}

	switch req.Method {
	case "create":
		return false, rpc.methodCreate(ctx, req.Params)
	case "start":
		return false, rpc.methodStart(ctx, req.Params)
	case "stop":
		return false, rpc.methodStop(ctx, req.Params)
	case "kill":
		return false, rpc.methodKill(ctx, req.Params)
	case "delete":
		return false, rpc.methodDelete(ctx, req.Params)
	case "inspect":
		return false, rpc.methodInspect(ctx, req.Params)
	default:
		return false, fmt.Errorf("invalid method: %s", req.Method)
	}
}

func (rpc *RpcServer) methodCreate(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcCreateRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	env := make([]string, 0, len(params.Env))
	for k, v := range params.Env {
		env = append(env, fmt.Sprintf(`%s="%s"`, k, v))
	}

	cmd := []string{"sh"}
	if len(params.Cmd) > 0 {
		cmd = params.Cmd
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	id, err := rpc.handler.Create(ctx, params.Name, params.Image, cmd, env)
	if err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcIdResponse{Id: id})
}

func (rpc *RpcServer) methodStart(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcStartRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	if err := rpc.handler.Start(ctx, params.Id); err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcEmptyResponse{})
}

func (rpc *RpcServer) methodStop(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcStopRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	if err := rpc.handler.Stop(ctx, params.Id); err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcEmptyResponse{})
}

func (rpc *RpcServer) methodKill(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcKillRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	if err := rpc.handler.Kill(ctx, params.Id, ""); err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcEmptyResponse{})
}

func (rpc *RpcServer) methodDelete(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcDeleteRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	if err := rpc.handler.Delete(ctx, params.Id); err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcEmptyResponse{})
}

func (rpc *RpcServer) methodInspect(ctx context.Context, rawParams json.RawMessage) error {
	var params RpcIdRequestParams
	if err := json.Unmarshal(rawParams, &params); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, rpc.timeout)
	defer cancel()

	data, err := rpc.handler.Inspect(ctx, params.Id)
	if err != nil {
		return rpc.sendResponse(NewRpcError(err))
	}

	return rpc.sendResponse(RpcInspectResponse{Data: data})
}

func (rpc *RpcServer) sendResponse(res any) error {
	return json.NewEncoder(rpc.stream).Encode(res)
}
