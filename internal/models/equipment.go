package models

type EquipmentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EquipmentTypeInput struct {
	Name string `json:"name" binding:"required"`
}

type GymEquipment struct {
	ID              int      `json:"id"`
	GymID           int      `json:"gym_id"`
	EquipmentTypeID int      `json:"equipment_type_id"`
	Weight          *float64 `json:"weight,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

type GymEquipmentWithDetails struct {
	ID              int      `json:"id"`
	GymID           int      `json:"gym_id"`
	EquipmentTypeID int      `json:"equipment_type_id"`
	EquipmentName   string   `json:"equipment_name"`
	Weight          *float64 `json:"weight,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

type GymEquipmentWithHistory struct {
	GymEquipmentWithDetails
	LastSessionWeight     *float64 `json:"last_session_weight,omitempty"`
	PreviousSessionWeight *float64 `json:"previous_session_weight,omitempty"`
}

type GymEquipmentInput struct {
	EquipmentTypeID int      `json:"equipment_type_id" binding:"required"`
	Weight          *float64 `json:"weight,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}
