package data

import "encoding/json"

type TroveData struct {
	Items []StandardProduct
}

type ProductMachineName string
type Platform string

type DownloadMethod string

const (
	Web        DownloadMethod = "web"
	BitTorrent DownloadMethod = "bittorrent"
)

type Filename string

type StandardProduct struct {
	DateAdded      int64                  `json:"date-added"`
	MachineName    string                 `json:"machine_name"`
	HumbleOriginal bool                   `json:"humble-original"`
	Downloads      map[Platform]*Download `json:"downloads"`
	HumanName      string                 `json:"human-name"`
}

type Download struct {
	UploadedAt  int64                       `json:"uploaded_at"`
	Name        string                      `json:"name"`
	Url         map[DownloadMethod]Filename `json:"url"`
	Timestamp   int64                       `json:"timestamp"`
	MachineName string                      `json:"machine_name"`
	FileSize    int64                       `json:"file_size"`
	Small       int                         `json:"small"`
	Size        string                      `json:"size"`
	Md5         string                      `json:"md5"`
}

func ParseFromHtml(html []byte) (*TroveData, error) {
	var items []StandardProduct
	err := json.Unmarshal(html, &items)
	if err != nil {
		return nil, err
	}
	return &TroveData{items}, err
}

func (td *TroveData) IsEmpty() bool {
	return len(td.Items) == 0
}
