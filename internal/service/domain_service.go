package service

import (
	"github.com/abyanjksatu/goscription/internal/outbound"
	"github.com/abyanjksatu/goscription/models"
)

type domainService struct {
	godaddyOutbound outbound.GodaddyOutbound
}

//DomainService represent the service of the domain
type DomainService interface {
	GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error)
}

// NewDomainService will create new an domainService object representation of service.DomainService interface
func NewDomainService(gd outbound.GodaddyOutbound) DomainService {
	return &domainService{godaddyOutbound: gd}
}

func (d *domainService) GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error) {
	return d.godaddyOutbound.GetDomainAvailable(domain)
}
