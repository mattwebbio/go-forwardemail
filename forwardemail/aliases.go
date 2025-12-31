package forwardemail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type AccountOrID struct {
	Account *Account
	ID      string
}

func (a *AccountOrID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		a.ID = s
		return nil
	}

	var acc Account
	if err := json.Unmarshal(data, &acc); err == nil {
		a.Account = &acc
		a.ID = acc.Id
		return nil
	}

	return fmt.Errorf("cannot unmarshal user field: %s", string(data))
}

type DomainOrID struct {
	Domain *Domain
	ID     string
}

func (d *DomainOrID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		d.ID = s
		return nil
	}

	var dom Domain
	if err := json.Unmarshal(data, &dom); err == nil {
		d.Domain = &dom
		d.ID = dom.Id
		return nil
	}

	return fmt.Errorf("cannot unmarshal domain field: %s", string(data))
}

type Alias struct {
	User                     AccountOrID `json:"user"`
	Domain                   DomainOrID  `json:"domain"`
	Name                     string      `json:"name"`
	Description              string      `json:"description"`
	Labels                   []string    `json:"labels"`
	IsEnabled                bool        `json:"is_enabled"`
	HasRecipientVerification bool        `json:"has_recipient_verification"`
	Recipients               []string    `json:"recipients"`
	Id                       string      `json:"id"`
	Object                   string      `json:"object"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
}

type AliasParameters struct {
	Recipients               *[]string
	Description              string `json:"description"`
	Labels                   *[]string
	HasRecipientVerification *bool
	IsEnabled                *bool
}

type GeneratePasswordParameters struct {
	NewPassword         *string
	Password            *string
	IsOverride          *bool
	EmailedInstructions *string
}

type GeneratedPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) GetAliases(domain string) ([]Alias, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/v1/domains/%s/aliases", domain))
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var items []Alias

	err = json.Unmarshal(res, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Client) GetAlias(domain string, alias string) (*Alias, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/v1/domains/%s/aliases/%s", domain, alias))
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item Alias

	err = json.Unmarshal(res, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) CreateAlias(domain string, alias string, parameters AliasParameters) (*Alias, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/v1/domains/%s/aliases", domain))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("name", alias)
	if parameters.Description != "" {
		params.Add("description", parameters.Description)
	}

	for k, v := range map[string]*bool{
		"has_recipient_verification": parameters.HasRecipientVerification,
		"is_enabled":                 parameters.IsEnabled,
	} {
		if v != nil {
			params.Add(k, strconv.FormatBool(*v))
		}
	}

	for k, v := range map[string]*[]string{
		"recipients[]": parameters.Recipients,
		"labels[]":     parameters.Labels,
	} {
		if v != nil {
			for _, vv := range *v {
				params.Add(k, vv)
			}
		}
	}

	req.Body = io.NopCloser(strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item Alias

	err = json.Unmarshal(res, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) UpdateAlias(domain string, alias string, parameters AliasParameters) (*Alias, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/v1/domains/%s/aliases/%s", domain, alias))
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("name", alias)
	if parameters.Description != "" {
		params.Add("description", parameters.Description)
	}

	for k, v := range map[string]*bool{
		"has_recipient_verification": parameters.HasRecipientVerification,
		"is_enabled":                 parameters.IsEnabled,
	} {
		if v != nil {
			params.Add(k, strconv.FormatBool(*v))
		}
	}

	for k, v := range map[string]*[]string{
		"recipients[]": parameters.Recipients,
		"labels[]":     parameters.Labels,
	} {
		if v != nil {
			for _, vv := range *v {
				params.Add(k, vv)
			}
		}
	}

	req.Body = io.NopCloser(strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item Alias

	err = json.Unmarshal(res, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) DeleteAlias(domain string, alias string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/v1/domains/%s/aliases/%s", domain, alias))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GenerateAliasPassword(domain string, alias string, parameters GeneratePasswordParameters) (*GeneratedPassword, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/v1/domains/%s/aliases/%s/generate-password", domain, alias))
	if err != nil {
		return nil, err
	}

	params := url.Values{}

	if parameters.NewPassword != nil {
		params.Add("new_password", *parameters.NewPassword)
	}
	if parameters.Password != nil {
		params.Add("password", *parameters.Password)
	}
	if parameters.IsOverride != nil {
		params.Add("is_override", strconv.FormatBool(*parameters.IsOverride))
	}
	if parameters.EmailedInstructions != nil {
		params.Add("emailed_instructions", *parameters.EmailedInstructions)
	}

	req.Body = io.NopCloser(strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var item GeneratedPassword

	err = json.Unmarshal(res, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
