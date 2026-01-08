package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type NominatimResult struct {
	Lat         string  `json:"lat"`
	Lon         string  `json:"lon"`
	DisplayName string  `json:"display_name"`
	Importance  float64 `json:"importance"`
}

func main() {
	db, err := sql.Open("sqlite3", "zones.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	zones := map[string][]string{
		"Johor": {
			"Pulau Aur", "Pulau Pemanggil", "Johor Bahru", "Kota Tinggi",
			"Mersing", "Kulai", "Kluang", "Pontian", "Batu Pahat",
			"Muar", "Segamat", "Gemas Johor", "Tangkak",
		},
		"Kedah": {
			"Kota Setar", "Kubang Pasu", "Pokok Sena", "Kuala Muda",
			"Yan", "Pendang", "Padang Terap", "Sik", "Baling",
			"Bandar Baharu", "Kulim", "Langkawi", "Puncak Gunung Jerai",
		},
		"Kelantan": {
			"Bachok", "Kota Bharu", "Machang", "Pasir Mas", "Pasir Puteh",
			"Tanah Merah", "Tumpat", "Kuala Krai", "Mukim Chiku",
			"Gua Musang (Daerah Galas Dan Bertam)", "Jeli", "Jajahan Kecil Lojing",
		},
		"Melaka": {
			"Seluruh Negeri Melaka",
		},
		"Negeri Sembilan": {
			"Tampin", "Jempol", "Jelebu", "Kuala Pilah", "Rembau",
			"Port Dickson", "Seremban",
		},
		"Pahang": {
			"Pulau Tioman", "Kuantan", "Pekan", "Muadzam Shah", "Jerantut",
			"Temerloh", "Maran", "Bera", "Chenor", "Jengka", "Bentong",
			"Lipis", "Raub", "Genting Sempah", "Janda Baik", "Bukit Tinggi",
			"Cameron Highlands", "Genting Higlands", "Bukit Fraser",
			"Mukim Rompin", "Mukim Endau", "Mukim Pontian",
		},
		"Perlis": {
			"Kangar", "Padang Besar", "Arau",
		},
		"Pulau Pinang": {
			"Seluruh Negeri Pulau Pinang",
		},
		"Perak": {
			"Tapah", "Slim River", "Tanjung Malim", "Kuala Kangsar",
			"Sg. Siput", "Ipoh", "Batu Gajah", "Kampar", "Lenggong",
			"Pengkalan Hulu", "Grik", "Temengor", "Belum", "Kg Gajah",
			"Teluk Intan", "Bagan Datuk", "Seri Iskandar", "Beruas",
			"Parit", "Lumut", "Sitiawan", "Pulau Pangkor", "Selama",
			"Taiping", "Bagan Serai", "Parit Buntar", "Bukit Larut",
		},
		"Sabah": {
			"Sandakan", "Bukit Garam", "Semawang", "Temanggong",
			"Tambisan", "Bandar Sandakan", "Sukau", "Beluran",
			"Telupid", "Pinangah", "Terusan", "Kuamut", "Lahad Datu",
			"Silabukan", "Kunak", "Sahabat", "Semporna", "Tungku",
			"Bandar Tawau", "Balong", "Merotai", "Kalabakan", "Kudat",
			"Kota Marudu", "Pitas", "Pulau Banggi", "Gunung Kinabalu",
			"Kota Kinabalu", "Ranau", "Kota Belud", "Tuaran",
			"Penampang", "Papar", "Putatan", "Pensiangan", "Keningau",
			"Tambunan", "Nabawan", "Beaufort", "Kuala Penyu", "Sipitang",
			"Tenom", "Long Pasia", "Membakut", "Weston",
		},
		"Selangor": {
			"Gombak", "Petaling", "Sepang", "Hulu Langat",
			"Hulu Selangor", "Shah Alam", "Kuala Selangor",
			"Sabak Bernam", "Klang", "Kuala Langat",
		},
		"Sarawak": {
			"Limbang", "Lawas", "Sundar", "Trusan", "Miri", "Niah",
			"Bekenu", "Sibuti", "Marudi", "Pandan", "Belaga", "Suai",
			"Tatau", "Sebauh", "Bintulu", "Sibu", "Mukah", "Dalat",
			"Song", "Igan", "Oya", "Balingian", "Kanowit", "Kapit",
			"Sarikei", "Matu", "Julau", "Rajang", "Daro", "Bintangor",
			"Belawai", "Lubok Antu", "Sri Aman", "Roban", "Debak",
			"Kabong", "Lingga", "Engkelili", "Betong", "Spaoh", "Pusa",
			"Saratok", "Serian", "Simunjan", "Samarahan", "Sebuyau",
			"Meludam", "Kuching", "Bau", "Lundu", "Sematan", "Kampung Patarikan",
		},
		"Terengganu": {
			"Kuala Terengganu", "Marang", "Kuala Nerus", "Besut",
			"Setiu", "Hulu Terengganu", "Dungun", "Kemaman",
		},
		"Wilayah Persekutuan": {
			"Kuala Lumpur", "Putrajaya", "Labuan",
		},
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	for state, cities := range zones {
		for _, city := range cities {
			query := fmt.Sprintf("%s, %s", city, state)

			fmt.Println("Processing", query)
			result, err := searchNominatim(client, query)
			if err != nil {
				fmt.Printf("skip %s: %s\n", query, err)
				continue
			}

			saveZone(db, city, state, result)
			time.Sleep(1 * time.Second)
		}
	}
}

func createTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS zones (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		city TEXT,
		state TEXT,
		lat REAL,
		lon REAL
	);
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func searchNominatim(client *http.Client, query string) (*NominatimResult, error) {
	endpoint := "https://nominatim.openstreetmap.org/search"

	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("limit", "1")

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:145.0) Gecko/20100101 Firefox/145.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var results []NominatimResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results")
	}

	return &results[0], nil
}

func saveZone(db *sql.DB, city, state string, r *NominatimResult) {
	stmt := `
	INSERT INTO zones (city, state, lat, lon) VALUES (?, ?, ?, ?);
	`

	_, err := db.Exec(stmt, city, state, r.Lat, r.Lon)
	if err != nil {
		fmt.Println("db error", err)
	}
}
