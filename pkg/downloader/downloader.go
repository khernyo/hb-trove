package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"hbtrove/pkg/checker"
	"hbtrove/pkg/data"
)

type Config struct {
	SessionCookie string `toml:"session-cookie"`
}

func NewConfigFromFile(file string) (*Config, error) {
	config := Config{}
	_, err := toml.DecodeFile(file, &config)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &config, nil
}

func (c *Config) cookie() string {
	return fmt.Sprintf("_simpleauth_sess=%v", c.SessionCookie)
}

func Download(config *Config, td *data.TroveData, dir string, checkContents bool, dryRun bool) error {
	results := checker.Check(td, dir, checkContents)

	count := 0
	for _, result := range results {
		if result.Status != checker.Same {
			count += 1
		}
	}

	fmt.Printf("Downloading %v files\n", count)
	for _, r := range results {
		if r.Status != checker.Same {
			fmt.Printf("Downloading %v %v %v ... ", r.Platform, r.Method, r.Product.HumanName)
			err := downloadItem(config, dryRun, r.Product.HumanName, r.Path, r.Download, r.Platform, r.Method)
			if err != nil {
				return errors.Wrap(err, "")
			}
			fmt.Println("Done.")
		}
	}
	return nil
}

func downloadItem(config *Config, dryRun bool, name string, path string, download *data.Download,
	platform data.Platform, method data.DownloadMethod) error {

	downloadUrl, err := getSignedUrl(config.cookie(), name, download, method)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if dryRun {
		fmt.Printf("Dry run. Not downloading [%v] [%v] %v ", platform, method, download.Url[method])
		return nil
	} else {
		err = downloadFile(downloadUrl, path, fmt.Sprintf("[%v] [%v] %v", platform, method, name))
		if err != nil {
			return errors.Wrap(err, "")
		}
		return nil
	}
}

func getSignedUrl(cookie string, name string, download *data.Download, method data.DownloadMethod) (string, error) {
	requestBody := url.Values{}
	requestBody.Set("machine_name", download.MachineName)
	requestBody.Set("filename", string(download.Url[method]))
	req, err := http.NewRequest(http.MethodPost, "https://www.humblebundle.com/api/v1/user/download/sign",
		strings.NewReader(requestBody.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	req.Header.Add("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	sdu := SignedDownloadUrl{}
	err = json.Unmarshal(responseBody, &sdu)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	var signedUrl string
	switch method {
	case data.Web:
		signedUrl = sdu.SignedUrl
	case data.BitTorrent:
		signedUrl = *sdu.SignedTorrentUrl
	default:
		panic("Unknown method: " + method)
	}

	return signedUrl, nil
}

type SignedDownloadUrl struct {
	SignedUrl        string  `json:"signed_url"`
	SignedTorrentUrl *string `json:"signed_torrent_url"`
}

func downloadFile(url string, targetPath string, msg string) error {
	dir := path.Dir(targetPath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "can't create directory '%v'", dir)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return errors.Wrapf(err, "can't create file '%v'", targetPath)
	}

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
