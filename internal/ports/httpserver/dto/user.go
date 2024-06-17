package dto

import "github.com/boldyrevn/mod-example/internal/model"

type CreateUserRequest struct {
    FirstName string                   `json:"firstName,omitempty"`
    LastName  string                   `json:"lastName,omitempty"`
    Age       uint32                   `json:"age,omitempty"`
    Profile   model.ProfileDescription `json:"profile"`
}

type UpdateUserRequest struct {
    ID        string                   `json:"ID,omitempty"`
    FirstName string                   `json:"firstName,omitempty"`
    LastName  string                   `json:"lastName,omitempty"`
    Age       uint32                   `json:"age,omitempty"`
    Profile   model.ProfileDescription `json:"profile"`
}

type RequestMessage struct {
    Message string `json:"message"`
}
