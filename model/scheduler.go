package model

type Scheduler struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Endpoint    string `json:"endpoint" bson:"endpoint"`
	Schedule    string `json:"schedule" bson:"schedule"`
	Enabled     bool   `json:"enabled" bson:"enabled"`
}
