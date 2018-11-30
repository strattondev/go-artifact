package cardset

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const cardSetUrl = "https://playartifact.com/cardset/"

type Receiver struct {
	Client *http.Client
}

type cdnResponse struct {
	CdnRoot    string `json:"cdn_root"`
	Url        string `json:"url"`
	ExpireTime int    `json:"expire_time"`
}

type CardSet struct {
	Set struct {
		Version  int     `json:"version"`
		SetInfo  SetInfo `json:"set_info"`
		CardList []Card  `json:"card_list"`
	} `json:"card_set"`
}

type SetInfo struct {
	SetId       int  `json:"set_id"`
	PackItemDef int  `json:"pack_item_def"`
	Name        Name `json:"name"`
}

type Card struct {
	CardName    Name        `json:"card_name"`
	MiniImage   Image       `json:"mini_image"`
	Illustrator string      `json:"illustrator"`
	GoldCost    int         `json:"gold_cost"`
	IngameImage Image       `json:"ingame_image"`
	CardId      int         `json:"card_id"`
	LargeImage  Image       `json:"large_image"`
	IsRed       bool        `json:"is_red"`
	IsBlue      bool        `json:"is_blue"`
	IsBlack     bool        `json:"is_black"`
	HitPoints   int         `json:"hit_points"`
	References  []Reference `json:"references"`
	BaseCardId  int         `json:"base_card_id"`
	ManaCost    int         `json:"mana_cost"`
	SubType     string      `json:"sub_type"`
	IsGreen     bool        `json:"is_green"`
	Armor       int         `json:"armor"`
	CardText    Name        `json:"card_text"`
	CardType    string      `json:"card_type"`
	Attack      int         `json:"attack"`
}

type Name struct {
	English    string `json:"english"`
	German     string `json:"german"`
	French     string `json:"french"`
	Italian    string `json:"italian"`
	Koreana    string `json:"koreana"`
	Spanish    string `json:"spanish"`
	Schinese   string `json:"schinese"`
	Tchinese   string `json:"tchinese"`
	Russian    string `json:"russian"`
	Thai       string `json:"thai"`
	Japanese   string `json:"japanese"`
	Portuguese string `json:"portuguese"`
	Polish     string `json:"polish"`
	Danish     string `json:"danish"`
	Dutch      string `json:"dutch"`
	Finnish    string `json:"finnish"`
	Norwegian  string `json:"norwegian"`
	Swedish    string `json:"swedish"`
	Hungarian  string `json:"hungarian"`
	Czech      string `json:"czech"`
	Romanian   string `json:"romanian"`
	Turkish    string `json:"turkish"`
	Brazilian  string `json:"brazilian"`
	Bulgarian  string `json:"bulgarian"`
	Greek      string `json:"greek"`
	Ukrainian  string `json:"ukrainian"`
	Latam      string `json:"latam"`
	Vietnamese string `json:"vietnamese"`
}

type Image struct {
	Default   string `json:"default"`
	German    string `json:"german,omitempty"`
	French    string `json:"french,omitempty"`
	Italian   string `json:"italian,omitempty"`
	Koreana   string `json:"koreana,omitempty"`
	Spanish   string `json:"spanish,omitempty"`
	Schinese  string `json:"schinese,omitempty"`
	Tchinese  string `json:"tchinese,omitempty"`
	Russian   string `json:"russian,omitempty"`
	Japanese  string `json:"japanese,omitempty"`
	Brazilian string `json:"brazilian,omitempty"`
	Latam     string `json:"latam,omitempty"`
}

type Reference struct {
	CardId  int    `json:"card_id"`
	RefType string `json:"ref_type"`
	Count   int    `json:"count,omitempty"`
}

func (c Receiver) RetrieveCardSet(setId string) (CardSet, error) {
	var cardSet CardSet

	cdnResponse, err := c.getCdnResponse(setId)

	if err != nil {
		return cardSet, err
	}

	response, err := c.Client.Get(cdnResponse.CdnRoot + cdnResponse.Url)

	if err != nil {
		return cardSet, err
	}

	defer response.Body.Close()
	bytes, err := getJsonBytes(response.Body)

	if err != nil {
		return cardSet, err
	}

	err = json.Unmarshal(bytes, &cardSet)

	return cardSet, err
}

func (c Receiver) getCdnResponse(setId string) (cdnResponse, error) {
	var cdnResponse cdnResponse
	response, err := c.Client.Get(cardSetUrl + setId)

	if err != nil {
		return cdnResponse, err
	}

	defer response.Body.Close()
	bytes, err := getJsonBytes(response.Body)

	if err != nil {
		return cdnResponse, err
	}

	err = json.Unmarshal(bytes, &cdnResponse)

	return cdnResponse, err
}

func getJsonBytes(body io.ReadCloser) ([]byte, error) {
	return ioutil.ReadAll(body)
}
