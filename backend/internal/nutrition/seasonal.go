package nutrition

import "time"

// SeasonalItem represents a seasonal produce item.
type SeasonalItem struct {
	Name     string `json:"name"`
	Category string `json:"category"` // Obst, Gemuese
	ImageURL string `json:"image_url,omitempty"`
}

// seasonalMap maps month (1-12) to seasonal produce in Germany.
var seasonalMap = map[time.Month][]SeasonalItem{
	time.January: {
		{Name: "Grünkohl", Category: "Gemüse"},
		{Name: "Rosenkohl", Category: "Gemüse"},
		{Name: "Feldsalat", Category: "Gemüse"},
		{Name: "Pastinake", Category: "Gemüse"},
		{Name: "Schwarzwurzel", Category: "Gemüse"},
		{Name: "Lauch", Category: "Gemüse"},
		{Name: "Wirsing", Category: "Gemüse"},
		{Name: "Äpfel", Category: "Obst"},
	},
	time.February: {
		{Name: "Grünkohl", Category: "Gemüse"},
		{Name: "Rosenkohl", Category: "Gemüse"},
		{Name: "Feldsalat", Category: "Gemüse"},
		{Name: "Chicorée", Category: "Gemüse"},
		{Name: "Lauch", Category: "Gemüse"},
		{Name: "Wirsing", Category: "Gemüse"},
		{Name: "Äpfel", Category: "Obst"},
	},
	time.March: {
		{Name: "Bärlauch", Category: "Gemüse"},
		{Name: "Feldsalat", Category: "Gemüse"},
		{Name: "Spinat", Category: "Gemüse"},
		{Name: "Lauch", Category: "Gemüse"},
		{Name: "Rhabarber", Category: "Obst"},
	},
	time.April: {
		{Name: "Bärlauch", Category: "Gemüse"},
		{Name: "Spargel", Category: "Gemüse"},
		{Name: "Spinat", Category: "Gemüse"},
		{Name: "Radieschen", Category: "Gemüse"},
		{Name: "Rhabarber", Category: "Obst"},
	},
	time.May: {
		{Name: "Spargel", Category: "Gemüse"},
		{Name: "Erdbeeren", Category: "Obst"},
		{Name: "Radieschen", Category: "Gemüse"},
		{Name: "Kohlrabi", Category: "Gemüse"},
		{Name: "Spinat", Category: "Gemüse"},
		{Name: "Rhabarber", Category: "Obst"},
	},
	time.June: {
		{Name: "Spargel", Category: "Gemüse"},
		{Name: "Erdbeeren", Category: "Obst"},
		{Name: "Kirschen", Category: "Obst"},
		{Name: "Blaubeeren", Category: "Obst"},
		{Name: "Erbsen", Category: "Gemüse"},
		{Name: "Zucchini", Category: "Gemüse"},
		{Name: "Kohlrabi", Category: "Gemüse"},
	},
	time.July: {
		{Name: "Tomaten", Category: "Gemüse"},
		{Name: "Gurken", Category: "Gemüse"},
		{Name: "Zucchini", Category: "Gemüse"},
		{Name: "Bohnen", Category: "Gemüse"},
		{Name: "Himbeeren", Category: "Obst"},
		{Name: "Johannisbeeren", Category: "Obst"},
		{Name: "Kirschen", Category: "Obst"},
		{Name: "Blaubeeren", Category: "Obst"},
		{Name: "Pflaumen", Category: "Obst"},
	},
	time.August: {
		{Name: "Tomaten", Category: "Gemüse"},
		{Name: "Paprika", Category: "Gemüse"},
		{Name: "Auberginen", Category: "Gemüse"},
		{Name: "Mais", Category: "Gemüse"},
		{Name: "Bohnen", Category: "Gemüse"},
		{Name: "Pflaumen", Category: "Obst"},
		{Name: "Birnen", Category: "Obst"},
		{Name: "Brombeeren", Category: "Obst"},
		{Name: "Mirabellen", Category: "Obst"},
	},
	time.September: {
		{Name: "Kürbis", Category: "Gemüse"},
		{Name: "Tomaten", Category: "Gemüse"},
		{Name: "Paprika", Category: "Gemüse"},
		{Name: "Pilze", Category: "Gemüse"},
		{Name: "Zwetschgen", Category: "Obst"},
		{Name: "Äpfel", Category: "Obst"},
		{Name: "Birnen", Category: "Obst"},
		{Name: "Trauben", Category: "Obst"},
	},
	time.October: {
		{Name: "Kürbis", Category: "Gemüse"},
		{Name: "Rote Bete", Category: "Gemüse"},
		{Name: "Grünkohl", Category: "Gemüse"},
		{Name: "Pilze", Category: "Gemüse"},
		{Name: "Sellerie", Category: "Gemüse"},
		{Name: "Äpfel", Category: "Obst"},
		{Name: "Birnen", Category: "Obst"},
		{Name: "Quitten", Category: "Obst"},
	},
	time.November: {
		{Name: "Grünkohl", Category: "Gemüse"},
		{Name: "Rosenkohl", Category: "Gemüse"},
		{Name: "Kürbis", Category: "Gemüse"},
		{Name: "Rote Bete", Category: "Gemüse"},
		{Name: "Pastinake", Category: "Gemüse"},
		{Name: "Feldsalat", Category: "Gemüse"},
		{Name: "Äpfel", Category: "Obst"},
	},
	time.December: {
		{Name: "Grünkohl", Category: "Gemüse"},
		{Name: "Rosenkohl", Category: "Gemüse"},
		{Name: "Feldsalat", Category: "Gemüse"},
		{Name: "Pastinake", Category: "Gemüse"},
		{Name: "Schwarzwurzel", Category: "Gemüse"},
		{Name: "Wirsing", Category: "Gemüse"},
		{Name: "Äpfel", Category: "Obst"},
	},
}

// SeasonalResponse is returned by the seasonal endpoint.
type SeasonalResponse struct {
	Month string         `json:"month"`
	Items []SeasonalItem `json:"items"`
}

// GetSeasonalProduce returns the seasonal produce for the current month.
func GetSeasonalProduce() SeasonalResponse {
	month := time.Now().Month()
	return SeasonalResponse{
		Month: monthName(month),
		Items: seasonalMap[month],
	}
}

// GetSeasonalIngredientNames returns a list of seasonal ingredient names for the current month.
func GetSeasonalIngredientNames() []string {
	items := seasonalMap[time.Now().Month()]
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}
	return names
}

func monthName(m time.Month) string {
	names := map[time.Month]string{
		time.January:   "Januar",
		time.February:  "Februar",
		time.March:     "März",
		time.April:     "April",
		time.May:       "Mai",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "August",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Dezember",
	}
	return names[m]
}
