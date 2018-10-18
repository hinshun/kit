package config

type ConfigLock struct {
	RefLocks []RefLock `json:"refLocks"`
}

type RefLock struct {
	Ref string `json:"ref"`
	Cid string `json:"cid"`
}
