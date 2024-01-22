package opensearch

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
	"crypto/tls"
)

var ErrNotFound  = errors.New("Resource not found")

type Client struct {
	httpClient *http.Client
	Host string
}

func NewClient(host string) *Client {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    	}
    	client := &http.Client{Transport: tr}
	return &Client{
		httpClient: client,
		Host: host,
	}
}

func (c *Client) GetTickData(ctx context.Context, tickNumber uint32) (TickDataResponse, error) {
	url := c.Host + "/tick/_doc/" + strconv.Itoa(int(tickNumber))

	var tickData TickDataResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &tickData)
	if err != nil {
		return TickDataResponse{}, errors.Wrap(err, "performing request")
	}

	return tickData, nil
}

func (c *Client) GetTx(ctx context.Context, id string) (TxResponse, error) {
	url := c.Host + "/txid/_doc/" + id

	var tx TxResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &tx)
	if err != nil {
		return TxResponse{}, errors.Wrap(err, "performing request")
	}

	return tx, nil
}

func (c *Client) GetBx(ctx context.Context, id string) (BxResponse, error) {
	url := c.Host + "/bxid/_doc/" + id

	var bx BxResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &bx)
	if err != nil {
		return BxResponse{}, errors.Wrap(err, "performing request")
	}

	return bx, nil
}

func (c *Client) GetStatus(ctx context.Context) (StatusResponse, error) {
	url := c.Host + "/status/_doc/api"

	var status StatusResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &status)
	if err != nil {
		return StatusResponse{}, errors.Wrap(err, "performing request")
	}

	return status, nil
}

func (c *Client) GetComputors(ctx context.Context, epoch int) (ComputorsResponse, error) {
	url := c.Host + "/computors/_doc/" + strconv.Itoa(epoch)

	var computors ComputorsResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &computors)
	if err != nil {
		return ComputorsResponse{}, errors.Wrap(err, "performing request")
	}

	return computors, nil
}

func (c *Client) GetQuorum(ctx context.Context, tick int) (QuorumResponse, error) {
	url := c.Host + "/quorum/_doc/" + strconv.Itoa(tick)

	var quorum QuorumResponse
	err := c.performRequest(ctx,url, http.MethodGet, nil, http.StatusOK, &quorum)
	if err != nil {
		return QuorumResponse{}, errors.Wrap(err, "performing request")
	}

	return quorum, nil
}

func (c *Client) performRequest(ctx context.Context, url string, method string, payload io.Reader, expectedStatusCode int, responseDest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	if res.StatusCode != expectedStatusCode {
		return errors.Errorf("Got unexpected status: %s. Body: %s", res.Status, string(body))
	}

	var osResponse Response

	err = json.Unmarshal(body, &osResponse)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response")
	}

	if !osResponse.Found {
		return ErrNotFound
	}

	source, err := json.Marshal(osResponse.Source)
	if err != nil {
		return errors.Wrap(err, "marshalling source")
	}

	err = json.Unmarshal(source, &responseDest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling body source")
	}

	return nil
}
