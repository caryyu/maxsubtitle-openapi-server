package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/caryyu/subtitle-open-server/internal/resource"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	var err error
	var subtitles []resource.Subtitle
	var resource = &resource.A4kDotNet{}
	var keyword string = "Matrix"

	if subtitles, err = resource.Search(keyword); err != nil {
		fmt.Println(err)
	}

	for _, s := range subtitles {
		fmt.Println(s.Desc, "-", s.Name)
	}

	//subtitle := subtitles[0]
	//var binary []byte
	//if err = resource.FetchDetail(&subtitles[0]); err != nil {
	//fmt.Errorf("%v", err)
	//fmt.Println(err)
	//}

	//for _, s := range subtitles[0].Binaries {
	//fmt.Println(s.Name)
	//}
	//binary, _ := Decodegbk(subtitles[0].Binaries[0].Bytes)
	//fmt.Println(string(binary))
	//fmt.Println(string(subtitles[0].Binaries[0].Bytes))
	//fmt.Println(subtitle)
	//for _, subtitle := range subtitles[:2] {
	//fmt.Println(subtitle.Desc)
	//}
}

func Decodegbk(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
