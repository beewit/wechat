package account

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/beewit/wechat/util"
)

// 二维码图片的URL, 可以GET此URL下载二维码或者在线显示此二维码.
func QRCodePicURL(ticket string) string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
func (clt *Client) QRCodeDownload(ticket, filepath string) (written int64, err error) {
	if ticket == "" {
		err = errors.New("empty ticket")
		return
	}
	i := strings.LastIndex(filepath, "/")
	dir := string([]rune(filepath)[0:i])
	err = util.MkDirAll(dir)
	if err != nil {
		return
	}
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(filepath)
		}
	}()

	if clt.HttpClient == nil {
		clt.HttpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, clt.HttpClient)
}

// 通过ticket换取二维码, 写入到 writer.
func (clt *Client) QRCodeDownloadToWriter(ticket string, writer io.Writer) (written int64, err error) {
	if ticket == "" {
		err = errors.New("empty ticket")
		return
	}
	if writer == nil {
		err = errors.New("nil writer")
		return
	}
	if clt.HttpClient == nil {
		clt.HttpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, clt.HttpClient)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
//  如果 clt == nil 则默认用 http.DefaultClient
func QRCodeDownload(ticket, filepath string, clt *http.Client) (written int64, err error) {
	if ticket == "" {
		err = errors.New("empty ticket")
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(filepath)
		}
	}()

	if clt == nil {
		clt = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, clt)
}

// 通过ticket换取二维码, 写入到 writer.
//  如果 clt == nil 则默认用 http.DefaultClient.
func QRCodeDownloadToWriter(ticket string, writer io.Writer, clt *http.Client) (written int64, err error) {
	if ticket == "" {
		err = errors.New("empty ticket")
		return
	}
	if writer == nil {
		err = errors.New("nil writer")
		return
	}
	if clt == nil {
		clt = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, clt)
}

func qrcodeDownloadToWriter(ticket string, writer io.Writer, clt *http.Client) (written int64, err error) {
	httpResp, err := clt.Get(QRCodePicURL(ticket))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	return io.Copy(writer, httpResp.Body)
}
