package main

const (
	SEARCH_ENDPOINT string = "/api/1/search"
)

type (
	SearchResult struct {
		Uid        string    `jsonapi:"primary,search"`
		Count      int       `jsonapi:"attr,count"`
		SearchTerm string    `jsonapi:"attr,search_term"`
		Results    []*Result `jsonapi:"relation,results"`
	}

	Result struct {
		Uid         string `jsonapi:"primary,result"`
		Kind        string `jsonapi:"attr,kind"` // podcast | episode
		Title       string `jsonapi:"attr,title"`
		Description string `jsonapi:"attr,description"`
		Url         string `jsonapi:"attr,url"`
		Feed        string `jsonapi:"attr,feed"`
		ImageUrl    string `jsonapi:"attr,image_url"`

		// metadata
		Score     int   `jsonapi:"attr,score"` // scaled to [0..100]
		Published int64 `jsonapi:"attr,published"`
	}
)
