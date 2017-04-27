package daocloudsdk

type BuildFlow struct {
	CreatedAt    int64  `json:"created_at"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	PackageID    string `json:"package_id"`
	Repo         string `json:"repo"`
	SrcOriginURL string `json:"src_origin_url"`
	SrcProvider  string `json:"src_provider"`
	Status       string `json:"status"`
}

type Build struct {
	Author      string `json:"author"`
	CreatedAt   string `json:"created_at"`
	ID          int64  `json:"id"`
	Ref         string `json:"ref"`
	RefIsBranch bool   `json:"ref_is_branch"`
	RefIsTag    bool   `json:"ref_is_tag"`
	Sha         string `json:"sha"`
	Status      string `json:"status"`
	Tag         string `json:"tag"`
}

type BuildParam struct {
	Branch string `json:"branch"`
}
