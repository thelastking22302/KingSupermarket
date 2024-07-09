package requsermodel

type SigninModel struct {
	Email    string `json:"email" validate:"required" bson:"email"`
	Password string `json:"password" validate:"required" bson:"password"`
}
