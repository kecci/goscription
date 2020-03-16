package usecase

import (
	"github.com/abyanjksatu/goscription/models"
	"github.com/abyanjksatu/goscription/outbound"
)

type domainUsecase struct {
	godaddyOutbound outbound.GodaddyOutbound
}

// NewDomainUsecase will create new an articleUsecase object representation of usecase.ArticleUsecase interface
func NewDomainUsecase(gd outbound.GodaddyOutbound) DomainUsecase {
	return &domainUsecase{godaddyOutbound: gd}
}

//DomainUsecase represent the usecase of the domain
type DomainUsecase interface {
	GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error)
}

func (d *domainUsecase) GetDomainAvailable(domain string) (res models.DomainAvailableResponse, err error) {
	return d.godaddyOutbound.GetDomainAvailable(domain)
}
