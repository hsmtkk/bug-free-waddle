package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Getter interface {
	Login(usrename, password string) error
	GetAlbum(url string) (string, error)
	ResolveDataSrc(paths []string) ([]string, error)
	GetThumbnail(urls []string) error
}

const enphotoURL = "https://en-photo.net"
const loginURL = enphotoURL + "/login"
const cookieFile = "cookie.json"

type getterImpl struct{}

func New() (Getter, error) {
	return &getterImpl{}, nil
}

func (g *getterImpl) Login(username, password string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize cookie jar; %w", err)
	}
	client := &http.Client{
		Jar: jar,
	}

	token, err := g.getLogin(client)
	if err != nil {
		return err
	}
	if err := g.postLogin(client, token, username, password); err != nil {
		return err
	}

	parsed, err := url.Parse(enphotoURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL; %s; %w", enphotoURL, err)
	}
	cs := client.Jar.Cookies(parsed)
	if err := saveCookie(cs); err != nil {
		return err
	}
	return nil
}

func (g *getterImpl) getLogin(client *http.Client) (string, error) {
	resp, err := client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("GetLogin failed; %w", err)
	}
	defer resp.Body.Close()
	if err := httpStatusError(resp); err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML; %w", err)
	}
	selection := doc.Find("[name=\"_token\"]")
	token, ok := selection.Attr("value")
	if !ok {
		return "", fmt.Errorf("failed to find token; %w", err)
	}
	return token, nil
}

func (g *getterImpl) postLogin(client *http.Client, token, username, password string) error {
	values := url.Values{
		"_token":   {token},
		"email":    {username},
		"password": {password},
	}
	resp, err := client.PostForm(loginURL, values)
	if err != nil {
		return fmt.Errorf("PostLogin failed; %w", err)
	}
	defer resp.Body.Close()
	return httpStatusError(resp)
}

func (g *getterImpl) GetAlbum(albumURL string) (string, error) {
	client, err := getClientWithCookie()
	if err != nil {
		return "", err
	}
	resp, err := client.Get(albumURL)
	if err != nil {
		return "", fmt.Errorf("GetAlbum failed; %w", err)
	}
	defer resp.Body.Close()
	if err := httpStatusError(resp); err != nil {
		return "", err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response; %w", err)
	}
	return string(respBytes), nil
}

func (g *getterImpl) ResolveDataSrc(paths []string) ([]string, error) {
	client, err := getClientWithCookie()
	if err != nil {
		return nil, err
	}
	urls := []string{}
	for _, path := range paths {
		url := enphotoURL + path
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("ResolveDataSrc failed; %w", err)
		}
		defer resp.Body.Close()
		if err := httpStatusError(resp); err != nil {
			return nil, err
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response; %w", err)
		}
		urls = append(urls, string(respBytes))
	}
	return urls, nil
}

func (g *getterImpl) GetThumbnail(urls []string) error {
	client, err := getClientWithCookie()
	if err != nil {
		return err
	}
	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err != nil {
			return fmt.Errorf("failed to parse URL; %s; %w", u, err)
		}
		elems := strings.Split(parsed.Path, "/")
		name := elems[len(elems)-1]
		path := filepath.Join("photo", name)
		resp, err := client.Get(u)
		if err != nil {
			return fmt.Errorf("GetThumbnail failed; %w", err)
		}
		defer resp.Body.Close()
		if err := httpStatusError(resp); err != nil {
			return fmt.Errorf("failed to get thumbnail; %s; %w", u, err)
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response; %w", err)
		}
		if err := os.WriteFile(path, respBytes, 0644); err != nil {
			return fmt.Errorf("failed to write file; %s; %w", path, err)
		}
	}
	return nil
}

func httpStatusError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non 2XX HTTP status code; %d; %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func saveCookie(cookie []*http.Cookie) error {
	jsBytes, err := json.Marshal(cookie)
	if err != nil {
		return fmt.Errorf("failed to marshal as JSON; %w", err)
	}
	if err := os.WriteFile(cookieFile, jsBytes, 0644); err != nil {
		return fmt.Errorf("failed to write file; %s; %w", cookieFile, err)
	}
	return nil
}

func loadCookie() ([]*http.Cookie, error) {
	jsBytes, err := os.ReadFile(cookieFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file; %s; %w", cookieFile, err)
	}
	var cs []*http.Cookie
	if err := json.Unmarshal(jsBytes, &cs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal as JSON; %w", err)
	}
	return cs, nil
}

func getClientWithCookie() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cookie jar; %w", err)
	}
	parsed, err := getParsedURL()
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL; %s; %w", enphotoURL, err)
	}
	cookie, err := loadCookie()
	if err != nil {
		return nil, err
	}
	jar.SetCookies(parsed, cookie)
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
}

func getParsedURL() (*url.URL, error) {
	parsed, err := url.Parse(enphotoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL; %s; %w", enphotoURL, err)
	}
	return parsed, nil
}
