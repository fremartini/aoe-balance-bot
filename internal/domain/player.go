package domain

type Player struct {
	PlayerId   uint
	SteamId    string
	Name       string
	Rating_1v1 *uint
	Rating_TG  *uint
}
