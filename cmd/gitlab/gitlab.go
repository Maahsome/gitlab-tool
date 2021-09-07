package giblab

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Gitlab struct {
	BaseUrl      string
	ApiPath      string
	RepoFeedPath string
	Token        string
	Client       *resty.Client
}

type PaginationOptions struct {
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type SortOptions struct {
	OrderBy string        `url:"order_by,omitempty"`
	Sort    SortDirection `url:"sort,omitempty"`
}

type ResponseWithMessage struct {
	Message string `json:"message"`
}

type ResponseMeta struct {
	Method     string
	Url        string
	StatusCode int
	RequestId  string
	Page       int
	PerPage    int
	PrevPage   int
	NextPage   int
	TotalPages int
	Total      int
	Runtime    float64
}

const (
	dateLayout = "2006-01-02T15:04:05-07:00"
)

var (
	skipCertVerify = flag.Bool("gitlab.skip-cert-check", false,
		`If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack.`)
)

// New generate a new gitlab client
func New(baseUrl, apiPath, token string) *Gitlab {

	// TODO: Add TLS Insecure && Pass in CA CRT for authentication
	restClient := resty.New()

	if apiPath == "" {
		apiPath = "/api/v4"
	}

	return &Gitlab{
		BaseUrl: baseUrl,
		ApiPath: apiPath,
		Token:   token,
		Client:  restClient,
	}
}

func (r *Gitlab) Get(uri string) (string, error) {

	nextPage := "1"
	combinedResults := ""

	for {
		// TODO: detect if there are no options passed in, ? verus & for page option
		fetchUri := fmt.Sprintf("https://%s%s%s&page=%s", r.BaseUrl, r.ApiPath, uri, nextPage)
		// logrus.Warn(fetchUri)
		resp, resperr := r.Client.R().
			SetHeader("PRIVATE-TOKEN", r.Token).
			Get(fetchUri)

		if resperr != nil {
			logrus.WithError(resperr).Error("Oops")
			return "", resperr
		}

		items := strings.TrimPrefix(string(resp.Body()[:]), "[")
		items = strings.TrimSuffix(items, "]")
		if combinedResults == "" {
			combinedResults += items
		} else {
			combinedResults += fmt.Sprintf(", %s", items)
		}
		currentPage := resp.Header().Get("X-Page")
		nextPage = resp.Header().Get("X-Next-Page")
		totalPages := resp.Header().Get("X-Total-Pages")
		if currentPage == totalPages {
			break
		}
	}
	return fmt.Sprintf("[%s]", combinedResults), nil
}
