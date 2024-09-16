package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	Status            string
	TenderServiceType string
)

const (
	StatusCreated Status = "Created"
	StatusPublish Status = "Published"
	StatusClosed  Status = "Closed"
	StatusCancel  Status = "Canceled"

	ServiceTypeConstruction TenderServiceType = "Construction"
	ServiceTypeDelivery     TenderServiceType = "Delivery"
	ServiceTypeManufacture  TenderServiceType = "Manufacture"
)

type TenderResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ServiceType string    `json:"serviceType"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TenderRequest struct {
	Name            string            `json:"name" validate:"required,max=100"`
	Description     string            `json:"description" validate:"required,max=500"`
	ServiceType     TenderServiceType `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationID  uuid.UUID         `json:"organizationId" validate:"required,uuid,not=00000000-0000-0000-0000-000000000000"`
	CreatorUsername string            `json:"creatorUsername" validate:"required,min=4"`
}

type TenderGetRequest struct {
	Limit       int               `validate:"omitempty,gt=0"`
	Offset      int               `validate:"omitempty,gt=-1"`
	ServiceType TenderServiceType `validate:"omitempty,oneof=Construction Delivery Manufacture"`
}

type TenderMyGetRequest struct {
	Limit    int    `validate:"omitempty,gt=0"`
	Offset   int    `validate:"omitempty,gt=-1"`
	Username string `validate:"required,min=4"`
}

type TenderStatusRequest struct {
	TenderId uuid.UUID `validate:"required,uuid,not=00000000-0000-0000-0000-000000000000"`
	Status   Status    `validate:"required,oneof=Created Published Closed"`
	Username string    `validate:"required,min=4"`
}

type TenderEditRequest struct {
	TenderId    uuid.UUID         `validate:"required,uuid,not=00000000-0000-0000-0000-000000000000"`
	Username    string            `validate:"required,min=4"`
	Name        string            `validate:"omitempty,max=100"`
	Description string            `validate:"omitempty,max=500"`
	ServiceType TenderServiceType `validate:"required,oneof=Construction Delivery Manufacture"`
}

type TenderEditRequestJSON struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceType TenderServiceType `json:"serviceType"`
}

type TenderRollbackRequest struct {
	TenderId uuid.UUID `validate:"required,uuid,not=00000000-0000-0000-0000-000000000000"`
	Username string    `validate:"required,min=4"`
	Version  int       `validate:"required,numeric,gte=1"`
}

type TenderGetStatusRequest struct {
	TenderId uuid.UUID `validate:"required,uuid,not=00000000-0000-0000-0000-000000000000"`
	Username string    `validate:"required,min=4"`
}
