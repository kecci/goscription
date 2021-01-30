package service

import (
	"github.com/kecci/goscription/internal/outbound"
	"github.com/kecci/goscription/models"
)

type (
	// DomainService represent the service of the domain
	DomainService interface {
		GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error)
	}

	// DomainServiceImpl represent the service of the domain
	DomainServiceImpl struct {
		godaddyOutbound outbound.GodaddyOutbound
	}
)

// NewDomainService will create new an domainService object representation of service.DomainService interface
func NewDomainService(gd outbound.GodaddyOutbound) DomainService {
	return &DomainServiceImpl{godaddyOutbound: gd}
}

// GetDomainAvailable ...
func (d *DomainServiceImpl) GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error) {
	return d.godaddyOutbound.GetDomainAvailable(domain)
}
