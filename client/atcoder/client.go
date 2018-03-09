package client

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	csrfTokenRe = regexp.MustCompile(`name="csrf_token"\svalue=('|")(.*?)('|")`)
)

type AtcoderClient struct {
	name    string
	pass    string
	logined bool
	client  *http.Client
	jar     *cookiejar.Jar
}

func (c *AtcoderClient) Auth(name, pass string) {
	c.name = name
	c.pass = pass
}

func New() *AtcoderClient {
	c := &AtcoderClient{}
	jar, _ := cookiejar.New(nil)
	c.jar = jar
	c.client = &http.Client{Jar: c.jar}
	c.Auth(os.Getenv("ATC_NAME"), os.Getenv("ATC_PASS"))
	c.Login()
	return c
}

func (c *AtcoderClient) Login() error {
	resp, err := c.client.Get("https://beta.atcoder.jp/login")
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	csrfToken, err := getCsrfToken(resp.Body)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Set("csrf_token", csrfToken)
	values.Add("username", c.name)
	values.Add("password", c.pass)
	req, err := http.NewRequest("POST", "https://beta.atcoder.jp/login?", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.logined = true
	return nil
}

func (c *AtcoderClient) Submit(contest, task, code string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://beta.atcoder.jp/%s/tasks/%s_%s", contest, contest, task), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	csrfToken, err := getCsrfToken(resp.Body)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Set("csrf_token", csrfToken)
	values.Add("data.LanguageId", "3013")
	values.Add("sourceCode", code)
	values.Add("data.TaskScreenName", fmt.Sprintf("%s_%s", contest, task))
	req, err = http.NewRequest("POST", fmt.Sprintf("https://beta.atcoder.jp/contests/%s/submit", contest), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	po, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(po))
	return nil
}

func getCsrfToken(body io.ReadCloser) (string, error) {
	html, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	match := csrfTokenRe.FindStringSubmatch(string(html))
	if len(match) < 2 {
		return "", errors.New("get csrf token failed")
	}
	return strings.Replace(match[2], "&#43;", "+", -1), nil
}
