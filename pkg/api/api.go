package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"hbtrove/pkg/data"
)

func LoadTroveData() (*data.TroveData, error) {
	var result *data.TroveData = nil
	for i := 0; ; i += 1 {
		chunk, err := LoadTroveDataChunk(i)
		if err != nil {
			return nil, err
		}

		if chunk.IsEmpty() {
			break
		} else {
			result = merge(result, chunk)
		}
	}

	return result, nil
}

func merge(d1 *data.TroveData, d2 *data.TroveData) *data.TroveData {
	if d1 == nil {
		return d2
	}
	return &data.TroveData{Items: append(d1.Items, d2.Items...)}
}

func LoadTroveDataChunk(idx int) (*data.TroveData, error) {
	println("Getting chunk", idx)

	resp, err := http.Get(fmt.Sprintf("https://www.humblebundle.com/api/v1/trove/chunk?index=%d", idx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data.ParseFromHtml(body)
}
