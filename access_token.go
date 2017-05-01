package clover

type AccessToken struct {
	AccessToken string                 `json:"access_token" bson:"_id" msg:"a"`
	ClientID    string                 `json:"client_id" bson:"c" msg:"c"`
	UserID      string                 `json:"user_id" bson:"u" msg:"u"`
	Expires     int64                  `json:"expires" bson:"e" msg:"e"`
	Scope       []string               `json:"scope" bson:"s" msg:"s"`
	Data        map[string]interface{} `json:"data" bson:"d" msg:"d"`
}

type AccessTokenCreator interface {
	Create() (string, error)
}
