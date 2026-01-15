// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package client

import "time"

// DnsRecord represents a DNS record in the Combell API
type DnsRecord struct {
	ID         string `json:"id,omitempty"`
	Type       string `json:"type"`
	RecordName string `json:"record_name,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	Content    string `json:"content,omitempty"`
	Priority   *int   `json:"priority,omitempty"`
	Service    string `json:"service,omitempty"`
	Weight     *int   `json:"weight,omitempty"`
	Target     string `json:"target,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Port       *int   `json:"port,omitempty"`
}

// Domain represents a domain in the Combell API (list view)
type Domain struct {
	DomainName     string     `json:"domain_name"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
	WillRenew      *bool      `json:"will_renew,omitempty"`
}

// DomainDetail represents detailed domain information
type DomainDetail struct {
	DomainName     string       `json:"domain_name"`
	ExpirationDate *time.Time   `json:"expiration_date,omitempty"`
	WillRenew      *bool        `json:"will_renew,omitempty"`
	CanToggleRenew *bool        `json:"can_toggle_renew,omitempty"`
	NameServers    []NameServer `json:"name_servers,omitempty"`
	Registrant     *Registrant  `json:"registrant,omitempty"`
}

// NameServer represents a nameserver
type NameServer struct {
	Name string `json:"name"`
	IP   string `json:"ip,omitempty"`
}

// Registrant represents domain registrant information
type Registrant struct {
	FirstName        string `json:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty"`
	Address          string `json:"address,omitempty"`
	PostalCode       string `json:"postal_code,omitempty"`
	City             string `json:"city,omitempty"`
	CountryCode      string `json:"country_code,omitempty"`
	Email            string `json:"email,omitempty"`
	Fax              string `json:"fax,omitempty"`
	Phone            string `json:"phone,omitempty"`
	LanguageCode     string `json:"language_code,omitempty"`
	CompanyName      string `json:"company_name,omitempty"`
	EnterpriseNumber string `json:"enterprise_number,omitempty"`
}

// EditNameServersRequest represents a request to update nameservers
type EditNameServersRequest struct {
	DomainName  string   `json:"domain_name"`
	NameServers []string `json:"name_servers"`
}

// EditDomainWillRenewRequest represents a request to update domain renewal state
type EditDomainWillRenewRequest struct {
	WillRenew bool `json:"will_renew"`
}
