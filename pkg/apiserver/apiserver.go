package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/juicedata/juicefs/pkg/utils"
	"io"
	"net/url"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) (*ApiServer, error) {
	u, err := url.JoinPath(addr)
	if err != nil {
		return nil, err
	}
	return &ApiServer{addr: u}, nil
}

func (s *ApiServer) AuthMount(volume string, ak, sk string) (AuthVolume, error) {
	authPath, err := url.JoinPath(s.addr, "/api/v1/juicefs/auth", volume)
	if err != nil {
		return AuthVolume{}, nil
	}
	params := make(map[string]string)
	params["access-key"] = ak
	params["secret-key"] = sk
	res, err := utils.HTTPPost(authPath, nil, params, nil)
	if err != nil {
		return AuthVolume{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return AuthVolume{}, err
	}

	fmt.Println(string(resBytes))

	api := apiAuthVolumeRes{}
	if err = json.Unmarshal(resBytes, &api); err != nil {
		return AuthVolume{}, err
	}

	if api.Code != 200 {
		return AuthVolume{}, fmt.Errorf("auth mount request faild, code is %d, msg: %s", api.Code, api.Message)
	}

	return api.Data, nil
}

type AuthVolume struct {
	BucketName string `json:"bucketName"`
	Auth       bool   `json:"auth"`
	MetaStore  string `json:"metaStore"`
}

type apiAuthVolumeRes struct {
	Code    int        `json:"code"`
	Data    AuthVolume `json:"data"`
	Message string     `json:"message"`
}
