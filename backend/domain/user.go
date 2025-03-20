package domain

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"` // Hashed password
	Online   bool   `json:"online" bson:"online"`
}