package kmgThirdCloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type SoftLayerSDK struct {
	Username string
	APIKey   string
}

func NewSoftLayerSDK(username, apiKey string) *SoftLayerSDK {
	sdk := &SoftLayerSDK{
		Username: username,
		APIKey:   apiKey,
	}
	return sdk
}

func (sdk *SoftLayerSDK) getUrl() string {
	return fmt.Sprintf("https://" + sdk.Username + ":" + sdk.APIKey + "@api.softlayer.com/rest/v3/SoftLayer_Account.json")
}

func (sdk *SoftLayerSDK) GetRegionList() []string {
	resp, err := http.Get(sdk.getUrl())
	handleErr(err)
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	return []string{}
}
