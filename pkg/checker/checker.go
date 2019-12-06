package checker

import (
	"crypto/md5"
	"io"
	"os"
	"path"

	"hbtrove/pkg/data"
)

type CheckResult struct {
	Product *data.StandardProduct
	Results []ProductCheckResult
}

type ProductCheckResult struct {
	Download *data.Download
	Platform data.Platform
	Results  []FileCheckResult
}

type FileCheckResult struct {
	Method data.DownloadMethod
	Path   string
	Status DownloadStatus
}

func Check(data *data.TroveData, dir string, checkContents bool) []CheckResult {
	var result []CheckResult
	for i := range data.Items {
		item := &data.Items[i]
		itemDir := path.Join(dir, item.MachineName)
		result = append(result, checkItem(itemDir, item, checkContents))
	}
	checkCollisions()
	return result
}

func checkItem(dir string, item *data.StandardProduct, checkContents bool) CheckResult {
	var results []ProductCheckResult
	for platform, download := range item.Downloads {
		results = append(results, checkDownload(path.Join(dir, string(platform)), download, platform, item.HumanName,
			checkContents))
	}
	return CheckResult{
		Product: item,
		Results: results,
	}
}

func checkDownload(dir string, download *data.Download, platform data.Platform, name string,
	checkContents bool) ProductCheckResult {
	var results []FileCheckResult
	for method, filename := range download.Url {
		filePath := path.Join(dir, string(filename))
		status := checkFile(download, method, filePath, checkContents)
		results = append(results, FileCheckResult{
			Method: method,
			Path:   filePath,
			Status: status,
		})
	}

	return ProductCheckResult{
		Download: download,
		Platform: platform,
		Results:  results,
	}
}

func checkFile(download *data.Download, method data.DownloadMethod, filePath string,
	checkContents bool) DownloadStatus {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return Missing
	}
	if err != nil || !fileInfo.Mode().IsRegular() {
		panic(err)
	}

	switch method {
	case data.Web:
		if fileInfo.Size() != download.FileSize {
			return Differ
		}
		if checkContents && computeMd5(filePath) != download.Md5 {
			return Differ
		}
		return Same
	case data.BitTorrent:
		// TODO extract data from existing torrent file and check size and hash
		return Same
	default:
		panic("Unknown method: " + method)
	}
}

func checkCollisions() {
	// TODO
}

func computeMd5(filePath string) string {
	println("Computing MD5 of", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}
	return toHex(h.Sum(nil))
}

func toHex(bytes []byte) string {
	chars := "0123456789abcdef"

	var result []rune = nil
	for _, b := range bytes {
		result = append(result, rune(chars[b>>4]))
		result = append(result, rune(chars[b&0xf]))
	}
	return string(result)
}

type DownloadStatus string

const (
	Missing DownloadStatus = "Missing"
	Differ  DownloadStatus = "Differ"
	Same    DownloadStatus = "Same"
)
