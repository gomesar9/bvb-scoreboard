package service

import "github.com/gomesar9/bvb-scoreboard/internal/domain/model"

func GetTest() {
	_ = model.Scoreboard{
		Stage:      "Chave Principal - Grupo F - Masc", // Main Draw - Pool F - Central court - Men #24
		LocCity:    "Brasília",
		LocCountry: "Brazil",
		PlayerA1:   "Arthur",
		PlayerA2:   "Adrielson",
		PlayerB1:   "Pedrosa",
		PlayerB2:   "Campos",
		ScoreA:     2,
		ScoreB:     0,
		FlagAURL:   "https://flagcdn.com/w40/br.png", // Bandeira do time A
		FlagBURL:   "https://flagcdn.com/w40/pt.png", // Bandeira do time B
		CountryA:   "Brasil",
		CountryB:   "Portugal",
		SetScores:  []string{"21–16", "21–19"}, // Parciais dos sets
	}
}
