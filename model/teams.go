package model

import (
	uuid "github.com/gofrs/uuid"
)

var teams = []Team{
	{ID: 1, Name: "Frontend Developers"},
	{ID: 2, Name: "Backend Developers"},
	{ID: 3, Name: "DevOps Team"},
}

type Team struct {
	ID   int       `json:"id" example:"1" format:"int64"`
	Name string    `json:"name" example:"team name"`
	UUID uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

func TeamsAll() []Team {
	return teams
}
