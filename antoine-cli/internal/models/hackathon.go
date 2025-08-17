package models

import (
	"time"
)

type Hackathon struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	URL              string                 `json:"url"`
	StartDate        time.Time              `json:"start_date"`
	EndDate          time.Time              `json:"end_date"`
	RegistrationURL  string                 `json:"registration_url"`
	Location         Location               `json:"location"`
	Technologies     []string               `json:"technologies"`
	Categories       []string               `json:"categories"`
	PrizePool        PrizeInfo              `json:"prize_pool"`
	Organizer        Organizer              `json:"organizer"`
	Difficulty       string                 `json:"difficulty"`
	TeamSize         TeamSizeInfo           `json:"team_size"`
	Requirements     []string               `json:"requirements"`
	Themes           []string               `json:"themes"`
	Status           string                 `json:"status"` // upcoming, active, completed
	ParticipantCount int                    `json:"participant_count"`
	ProjectCount     int                    `json:"project_count"`
	Tags             []string               `json:"tags"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type Location struct {
	Type        string  `json:"type"` // online, in-person, hybrid
	City        string  `json:"city,omitempty"`
	Country     string  `json:"country,omitempty"`
	Venue       string  `json:"venue,omitempty"`
	Address     string  `json:"address,omitempty"`
	Timezone    string  `json:"timezone"`
	Coordinates *LatLng `json:"coordinates,omitempty"`
}

type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PrizeInfo struct {
	Total       int      `json:"total"`
	Currency    string   `json:"currency"`
	Breakdown   []Prize  `json:"breakdown"`
	Sponsors    []string `json:"sponsors"`
	NonMonetary []string `json:"non_monetary"`
}

type Prize struct {
	Position    string `json:"position"` // 1st, 2nd, 3rd, best-in-category
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Sponsor     string `json:"sponsor,omitempty"`
}

type Organizer struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"` // company, community, academic
	Website string  `json:"website"`
	Social  Social  `json:"social"`
	Contact Contact `json:"contact"`
}

type Social struct {
	Twitter  string `json:"twitter,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	Discord  string `json:"discord,omitempty"`
	Telegram string `json:"telegram,omitempty"`
}

type Contact struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type TeamSizeInfo struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
