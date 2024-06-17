package model

type User struct {
    ID        string             `json:"ID,omitempty" bson:"ID"`
    FirstName string             `json:"firstName,omitempty"`
    LastName  string             `json:"lastName,omitempty"`
    Age       uint32             `json:"age,omitempty"`
    Profile   ProfileDescription `json:"profile"`
}

type ProfileDescription struct {
    Bio       string   `json:"bio,omitempty"`
    Interests []string `json:"interests,omitempty"`
}
