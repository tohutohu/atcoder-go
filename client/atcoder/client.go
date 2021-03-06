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

	"github.com/PuerkitoBio/goquery"
)

var (
	csrfTokenRe = regexp.MustCompile(`name="csrf_token"\svalue=('|")(.*?)('|")`)
)

type Sample struct {
	Input  string
	Output string
}

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

func New() (*AtcoderClient, error) {
	c := &AtcoderClient{}
	jar, _ := cookiejar.New(nil)
	c.jar = jar
	c.client = &http.Client{Jar: c.jar}
	c.Auth(os.Getenv("ATC_NAME"), os.Getenv("ATC_PASS"))
	err := c.Login()
	return c, err
}

func (c *AtcoderClient) Login() error {
	resp, err := c.client.Get("https://beta.atcoder.jp/login")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	csrfToken, err := getCsrfToken(resp.Body)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Set("csrf_token", csrfToken)
	values.Add("username", c.name)
	values.Add("password", c.pass)
	resp, err = c.client.PostForm("https://beta.atcoder.jp/login", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c.logined = true
	return nil
}

func (c *AtcoderClient) Submit(contest, task, code string) error {
	resp, err := c.client.Get(fmt.Sprintf("https://beta.atcoder.jp/%s/tasks/%s_%s", contest, contest, task))
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
	resp, err = c.client.PostForm(fmt.Sprintf("https://beta.atcoder.jp/contests/%s/submit", contest), values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

func (c *AtcoderClient) GetTaskInfo(contest, problem string) ([]Sample, string, error) {
	samples := []Sample{}
	url := fmt.Sprintf("https://beta.atcoder.jp/contests/%s/tasks/%s_%s", contest, contest, problem)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return samples, "", nil
	}
	sample := Sample{}
	doc.Find("div.part>section>pre").Each(func(_ int, s *goquery.Selection) {
		if s.Parent().Parent().Parent().HasClass("io-style") {
			return
		}
		if sample.Input == "" {
			sample.Input = s.Text()
		} else {
			sample.Output = s.Text()
			samples = append(samples, sample)
			sample = Sample{}
		}
	})
	return samples, doc.Find("#task-statement>span>span.lang-ja").Text(), err
}
