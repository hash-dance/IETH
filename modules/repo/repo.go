package repo

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var (
	iethAPI = "api"
)

type Repo struct {
	path string
	api  string
}

func (r *Repo) GetApi() ([]byte, error) {
	p := filepath.Join(r.path, iethAPI)
	f, err := os.Open(p)

	if os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // Read only op

	tb, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(tb), nil
}

func (r *Repo) SetApi(api string) error {
	p := filepath.Join(r.path, iethAPI)
	f, err := os.Create(p)

	if os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}
	defer f.Close() //nolint: errcheck // Read only op

	_, err = f.WriteString(api)
	if err != nil {
		return err
	}

	return nil
}

func NewRepo(path string) (*Repo, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	return &Repo{
		path: path,
	}, nil
}
