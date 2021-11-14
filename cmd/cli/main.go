package main

import (
	"fmt"

	"github.com/caryyu/subtitle-open-server/internal/resource"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "mast",
	}
	var searchCmd = &cobra.Command{
		Use:   "search",
		Short: "To search the result by a given keyword",
		Long:  "Example: mast search Mulan",
		Run:   search,
	}
	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "To download a single subtitle by an id from the search result",
		Long:  "Example: mast download abcd",
		Run:   download,
	}

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.Execute()

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

//func Decodegbk(s []byte) ([]byte, error) {
//I := bytes.NewReader(s)
//O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
//d, e := ioutil.ReadAll(O)
//if e != nil {
//return nil, e
//}
//return d, nil
//}

func search(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
		panic("No keyword Found")
	}
	var err error
	var subtitles []resource.Subtitle
	var resource = resource.NewA4kDotNet()
	var keyword string = args[0]

	if subtitles, err = resource.Search(keyword); err != nil {
		fmt.Println(err)
	}

	for _, s := range subtitles {
		fmt.Println(s.Id, "-", s.Name)
	}
}

func download(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cmd.Help()
		panic("No id Found")
	}

	var id string = args[0]
	var resource = resource.NewA4kDotNet()
	file, err := resource.GetFromCache(id)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(file.Name)
}
