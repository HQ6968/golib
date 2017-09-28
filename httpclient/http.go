package httpclient

import (
	"net/http"
	"net/url"
)

type HttpClient struct {
	client *http.Client
	Header map[string]string
}

func NewHttp() *HttpClient {
	client := &http.Client{}
	return &HttpClient{
		client:client,
		Header : make(map[string]string),
	}
}

func (this *HttpClient)PostForm(url string , data url.Values) (*http.Response , error) {
	request , err := http.NewRequest("Post" , url , nil)
	if err != nil {
		return nil , err
	}
	request.PostForm = data
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return this.client.Do(request)
}