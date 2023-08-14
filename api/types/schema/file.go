package schema

import "time"

type File struct {
	ID         int64
	Createtime time.Time
	Filename   string
	Owner      string
	Expired    time.Time
	Fileid     string
	Bucketname string
}

type FileInfo struct {
	Fileid     string    `json:"id"`
	Filename   string    `json:"fileName"`
	Owner      string    `json:"owner"`
	Createtime time.Time `json:"createTime,omitempty"`
	Expired    time.Time `json:"expireDate"`
	URL        string    `json:"url,omitempty" gorm:"-:all"`
}

type Files []*FileInfo

type FileQueryResult struct {
	PageResult *PaginationResult `json:"pageResult"`
	PageData   Files             `json:"pageData"`
}

type FileSign struct {
	Filename  string `json:"fileName"`
	AuthParam string `json:"authParam"`
}

// 上传成功后返回文件名和短url
type FileURL struct {
	File string `json:"file"`
	URL  string `json:"url"`
}
