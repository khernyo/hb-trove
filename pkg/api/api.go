package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"hbtrove/pkg/data"
)

func LoadTroveData() (jsons [][]byte, td *data.TroveData, err error) {
	for i := 0; ; i += 1 {
		json, chunk, err := LoadTroveDataChunk(i)
		if err != nil {
			return nil, nil, err
		}

		if chunk.IsEmpty() {
			break
		} else {
			jsons = append(jsons, json)
			td = merge(td, chunk)
		}
	}

	return jsons, td, nil
}

func merge(d1 *data.TroveData, d2 *data.TroveData) *data.TroveData {
	if d1 == nil {
		return d2
	}
	return &data.TroveData{Items: append(d1.Items, d2.Items...)}
}

func LoadTroveDataChunk(idx int) (json []byte, td *data.TroveData, err error) {
	println("Getting chunk", idx)

	resp, err := http.Get(fmt.Sprintf("https://www.humblebundle.com/api/v1/trove/chunk?index=%d", idx))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	json, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	td, err = data.ParseFromJson(json)
	return json, td, err
}
