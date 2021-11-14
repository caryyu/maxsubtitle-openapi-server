package resource

type Subtitle struct {
	Id         string `json:"id"`
	OriginalId string `json:"originalId"`
	Desc       string `json:"desc"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Format     string `json:"format"`
}
