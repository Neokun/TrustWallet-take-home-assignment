package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/TrustWallet/tx-parser/internal/types"
	"github.com/TrustWallet/tx-parser/internal/utils"
)

type Client interface {
	BlockNumber(ctx context.Context) (uint64, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
}

type ethereumClient struct {
	rpcNode string
}

func NewEthereumClient(rpcNode string) *ethereumClient {
	return &ethereumClient{rpcNode: rpcNode}
}

func (c *ethereumClient) BlockNumber(ctx context.Context) (uint64, error) {
	var result utils.HexUint64
	err := c.callMethod(ctx, &result, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  blockNumberMethod,
		"params":  []string{},
		"id":      83,
	})

	if err != nil {
		return 0, err
	}
	return uint64(result), nil
}

func (c *ethereumClient) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	number := utils.EncodeUint64(blockNumber)
	var raw json.RawMessage
	err := c.callMethod(ctx, &raw, map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  getBlockByNumberMethod,
		"params":  []interface{}{number, true},
		"id":      1,
	})
	if err != nil {
		return nil, err
	}

	var block types.Block
	err = json.Unmarshal(raw, &block)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(block.Transactions); i++ {
		block.Transactions[i].Timestamp = block.Timestamp
	}

	return &block, nil
}

func (c *ethereumClient) callMethod(ctx context.Context, result interface{}, callBody interface{}) error {
	if result != nil && reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("call result parameter must be pointer or nil interface: %v", result)
	}
	respBody, err := c.doRequest(ctx, callBody)
	if err != nil {
		return err
	}
	defer respBody.Close()

	var respmsg jsonrpcMessage
	if err = json.NewDecoder(respBody).Decode(&respmsg); err != nil {
		return err
	}
	if respmsg.Error != nil {
		return respmsg.Error
	}
	if len(respmsg.Result) == 0 {
		return ErrNoResult
	}

	if result == nil {
		return nil
	}
	return json.Unmarshal(respmsg.Result, result)
}

func (c *ethereumClient) doRequest(ctx context.Context, msg interface{}) (io.ReadCloser, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rpcNode, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		var body []byte
		if _, err := buf.ReadFrom(resp.Body); err == nil {
			body = buf.Bytes()
		}

		return nil, HTTPError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}
	return resp.Body, nil
}
