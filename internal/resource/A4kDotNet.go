package resource

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
	astisub "github.com/asticode/go-astisub"
	"github.com/mholt/archiver/v3"
)

const site string = "https://www.a4k.net"

type File struct {
	Name  string `json:"name"`
	Bytes []byte `json:"bytes"`
}

type A4kDotNet struct {
	lock      *sync.Mutex
	CachePath string
}

func NewA4kDotNet() *A4kDotNet {
	return &A4kDotNet{
		lock:      &sync.Mutex{},
		CachePath: "/tmp/data-cache",
	}
}

func (r *A4kDotNet) Search(keyword string) (subtitles []Subtitle, err error) {
	soup.Header("User-Agent", "curl/7.64.1")
	resp, err := soup.Get(fmt.Sprintf("%s/search?term=%s", site, keyword))
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(resp)
	items := doc.FindStrict("ul", "class", "ui relaxed divided list").FindAll("li")

	//var funcGetLanguages func(nodes []soup.Root) []string = func(nodes []soup.Root) []string {
	//var result []string
	//for _, item := range nodes {
	//language := item.Attrs()["data-content"]
	//result = append(result, language)
	//}
	//return result
	//}

	size := int(math.Min(3, float64(len(items))))
	// Only load top 3 items
	for _, item := range items[:size] {
		i := item.FindStrict("div", "class", "content").Find("h3").Find("a")
		id := i.Attrs()["href"][strings.LastIndex(i.Attrs()["href"], "/")+1:]
		name := i.Text()

		if resp, err = soup.Get(fmt.Sprintf("%s/subtitle/%s", site, id)); err != nil {
			log.Printf("Can't load the detail page of %s: %v\n", id, err)
			continue
		}

		doc := soup.HTMLParse(resp)

		// /system/files/subtitle/2021-10/a4k.net_1634455869.rar
		url := doc.FindStrict("div", "class", "download").FindStrict("a", "class", "ui green button").Attrs()["href"]
		url = fmt.Sprintf("%s/%s", site, url)

		var items []Subtitle
		if items, err = r.download(id, name, url); err != nil {
			log.Printf("Can't download the file of %s: %v\n", url, err)
			continue
		}

		subtitles = append(subtitles, items...)

		//subtitle := &Subtitle{
		//Id:        i.Attrs()["href"][strings.LastIndex(i.Attrs()["href"], "/")+1:],
		//Desc:      i.Text(),
		//Languages: funcGetLanguages(item.FindStrict("div", "class", "language").Find("span", "class", "h4").FindAll("i")),
		//}

		//subtitles = append(subtitles, *subtitle)
	}

	return subtitles, nil
}

//func (r *A4kDotNet) FetchDetail(model *Subtitle) (err error) {
//soup.Header("User-Agent", "curl/7.64.1")
//var resp string

//if resp, err = soup.Get(fmt.Sprintf("%s/subtitle/%s", site, model.Id)); err != nil {
//return err
//}

//doc := soup.HTMLParse(resp)

//// /system/files/subtitle/2021-10/a4k.net_1634455869.rar
//url := doc.FindStrict("div", "class", "download").FindStrict("a", "class", "ui green button").Attrs()["href"]
//model.Url = fmt.Sprintf("%s/%s", site, url)

//if err = r.binaryDownload(model); err != nil {
//return err
//}

//return nil
//}

func (r *A4kDotNet) download(id string, name string, url string) (subtitles []Subtitle, err error) {
	var resp *http.Response
	var req *http.Request

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.Header.Set("User-Agent", "curl/7.64.1")
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var bytes []byte
	var files []File
	switch ext := filepath.Ext(url); strings.ToLower(ext) {
	case ".ass", ".ssa":
		if bytes, err = r.fromASSToSRT(resp.Body); err != nil {
			return nil, err
		}
		files = []File{{Name: "default-" + id + ".srt", Bytes: bytes}}
	case ".srt":
		if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		}
		files = []File{{Name: "default-" + id + ext, Bytes: bytes}}
	case ".rar", ".zip":
		files = r.extract(resp.Body, ext)
	}

	if len(files) > 0 {
		r.cacheFiles(&files)

		for _, file := range files {
			hash := sha256.Sum256([]byte(file.Name))

			subtitle := Subtitle{
				Name:       file.Name,
				Id:         hex.EncodeToString(hash[:]),
				OriginalId: id,
				Desc:       name,
				Url:        url,
				Format:     filepath.Ext(file.Name)[1:],
			}
			subtitles = append(subtitles, subtitle)
		}
	}

	return subtitles, nil
}

func (r *A4kDotNet) fromASSToSRT(reader io.Reader) ([]byte, error) {
	s, err := astisub.ReadFromSSA(reader)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	err = s.WriteToSRT(buffer)
	if err != nil {
		return nil, err
	}

	bytes := buffer.Bytes()

	return bytes, nil
}

/**
* Extracting Zip/Rar archives from Http downloaded payload
 */
func (r *A4kDotNet) extract(reader io.Reader, ext string) []File {
	r.lock.Lock()
	defer r.lock.Unlock()

	var err error

	var in *os.File
	var source string = "/tmp/data.bin"
	var dest string = "/tmp/data"
	if in, err = os.Create(source); err != nil {
		panic(err)
	}
	defer func() {
		in.Close()
		if err = os.RemoveAll(dest); err != nil {
			panic(err)
		}
		if err = os.RemoveAll(source); err != nil {
			panic(err)
		}
	}()

	if _, err = os.Stat(dest); os.IsNotExist(err) {
		os.Mkdir(dest, os.ModePerm)
	}

	if _, err = io.Copy(in, reader); err != nil {
		log.Println("Copying data is wrong")
		panic(err)
	}

	var i interface{}
	if i, err = archiver.ByExtension(ext); err != nil {
		log.Println("Getting the unachiver encountered an issue")
		panic(err)
	}

	u, _ := i.(archiver.Unarchiver)

	if err = u.Unarchive(source, dest); err != nil {
		panic(err)
	}

	files := make([]File, 0)
	r.flattenToMemory(dest, &files)

	return files
}

/**
* Looking into every directory and flattening the file list
* Note: .txt will be ignored
 */
func (r *A4kDotNet) flattenToMemory(dest string, files *[]File) {
	var entries []fs.DirEntry
	var err error
	var out []byte

	if entries, err = os.ReadDir(dest); err != nil {
		panic(err)
	}

	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if ext == ".txt" || ext == ".torrent" || ext == ".jpg" || ext == ".png" {
			continue
		}
		path := dest + "/" + entry.Name()
		if entry.IsDir() {
			r.flattenToMemory(path, files)
			continue
		}
		if out, err = ioutil.ReadFile(path); err != nil {
			panic(err)
		}

		name := entry.Name()
		switch strings.ToLower(ext) {
		case ".ssa", ".ass":
			in := bytes.NewReader(out)
			if out, err = r.fromASSToSRT(in); err != nil {
				log.Println("Skipped as it cannot be converted")
				continue
			}
			name = name[0:len(name)-len(ext)] + ".srt"
		}

		file := File{Name: name, Bytes: out}
		*files = append(*files, file)
	}
}

func (r *A4kDotNet) cacheFiles(files *[]File) {
	path := r.CachePath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	for _, file := range *files {
		dest := fmt.Sprintf("%s/%s", path, file.Name)
		if err := os.WriteFile(dest, file.Bytes, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func (r *A4kDotNet) GetFromCache(id string) (file *File, err error) {
	path := r.CachePath
	var entries []fs.DirEntry
	var bytes []byte

	if entries, err = os.ReadDir(path); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		hash := sha256.Sum256([]byte(entry.Name()))
		hashHex := hex.EncodeToString(hash[:])

		if hashHex == id {
			dest := fmt.Sprintf("%s/%s", path, entry.Name())
			if bytes, err = ioutil.ReadFile(dest); err != nil {
				return nil, err
			}
			file = &File{Bytes: bytes, Name: entry.Name()}
			break
		}
	}

	return file, nil
}
