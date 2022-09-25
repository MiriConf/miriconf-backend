package model

// Team model info
// @Description Team information
type Team struct {
	// ID of team
	ID int `json:"id" example:"1" format:"int64"`
	// Name of team
	Name string `json:"name" example:"cyber operations" format:"string"`
	// Parent department of team
	Department string `json:"department" example:"cybersecurity" format:"string"`
	// Tags assigned to the team using the UI
	Tags string `json:"tags" example:"security" format:"string"`
	// Creation timestamp
	CreatedAt string `json:"createdat" example:"08-23-2022" format:"string"`
}

func TeamsAll() Team {
	testData := Team{ID: 4, Name: "backend development team", Department: "development", Tags: "development", CreatedAt: "08-5-2022"} 
    
	return testData
	//newClient := database.MongoConnection()
	//database.GetTeams(newClient)
}

//func TeamsAll() []Team {
//	return teams
//}
