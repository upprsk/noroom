package rpc

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type RpcClient struct {
	stream io.ReadWriter
}

func NewRpcClient(stream io.ReadWriter) *RpcClient {
	return &RpcClient{
		stream: stream,
	}
}

func (rpc *RpcClient) Create(name, image string) (string, error) {
	req, err := NewRpcCreateRequest(RpcCreateRequestParams{Name: name, Image: image})
	if err != nil {
		return "", err
	}

	var res RpcCreateResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return "", err
	}

	return res.Id, nil
}

func (rpc *RpcClient) Start(id string, timeout time.Duration) error {
	req, err := NewRpcStartRequest(RpcStartRequestParams{
		Id:      id,
		Timeout: timeout,
	})
	if err != nil {
		return err
	}

	var res RpcEmptyResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return err
	}

	return nil
}

func (rpc *RpcClient) Stop(id string, timeout time.Duration) error {
	req, err := NewRpcStopRequest(RpcStopRequestParams{
		Id:      id,
		Timeout: timeout,
	})
	if err != nil {
		return err
	}

	var res RpcEmptyResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return err
	}

	return nil
}

func (rpc *RpcClient) Kill(id string, timeout time.Duration) error {
	req, err := NewRpcKillRequest(RpcKillRequestParams{
		Id:      id,
		Timeout: timeout,
	})
	if err != nil {
		return err
	}

	var res RpcEmptyResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return err
	}

	return nil
}

func (rpc *RpcClient) Delete(id string) error {
	req, err := NewRpcDeleteRequest(RpcDeleteRequestParams{Id: id})
	if err != nil {
		return err
	}

	var res RpcEmptyResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return err
	}

	return nil
}

func (rpc *RpcClient) Inspect(id string) (*ContainerInspectResult, error) {
	req, err := NewRpcInspectRequest(RpcInspectRequestParams{Id: id})
	if err != nil {
		return nil, err
	}

	var res RpcInspectResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (rpc *RpcClient) Attach(id string) error {
	req, err := NewRpcAttachRequest(RpcInspectRequestParams{Id: id})
	if err != nil {
		return err
	}

	var res RpcEmptyResponse
	if err := sendMessage(rpc.stream, req, &res); err != nil {
		return err
	}

	// we don't want to use this for RPC anymore
	rpc.stream = nil

	return nil
}

var (
	ErrNilStream = errors.New("nil stream")
)

func sendMessage[R interface{ GetErr() error }](stream io.ReadWriter, req RpcRequest, res R) error {
	if stream == nil {
		return ErrNilStream
	}

	if err := json.NewEncoder(stream).Encode(req); err != nil {
		return err
	}

	if err := json.NewDecoder(stream).Decode(res); err != nil {
		return err
	}

	if err := res.GetErr(); err != nil {
		return err
	}

	return nil
}
