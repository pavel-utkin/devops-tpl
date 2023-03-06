package handlers

import (
	"devops-tpl/internal/server/storage"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func UpdateGaugePost(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	statName := chi.URLParam(request, "statName")
	statValue := chi.URLParam(request, "statValue")
	statValueInt, err := strconv.ParseFloat(statValue, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Bad request"))
		return
	}

	err = memStatsStorage.UpdateGaugeValue(statName, statValueInt)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Server error"))
		return
	}

	fmt.Println("Update gauge:")
	fmt.Printf("%v: %v\n\n", statName, statValue)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Ok"))
}

func UpdateCounterPost(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	statName := chi.URLParam(request, "statName")
	statValue := chi.URLParam(request, "statValue")
	statCounterValue, err := strconv.ParseInt(statValue, 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	err = memStatsStorage.UpdateCounterValue(statName, statCounterValue)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	fmt.Println("Inc counter:")
	fmt.Printf("%v: %v", statName, statValue)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Ok"))
}

func PrintStatsValues(rw http.ResponseWriter, request *http.Request, memStatsStorage storage.MemStatsMemoryRepo) {
	htmlTemplate := `
		<html>
			<head>
			<title></title>
			</head>
			<body>
				%v
			</body>
		</html>`
	keyValuesHTML := ""

	for k, v := range memStatsStorage.GetDBSchema() {
		keyValuesHTML += fmt.Sprintf("<div><b>%v</b>: %v</div>", k, v)
	}

	htmlPage := fmt.Sprintf(htmlTemplate, keyValuesHTML)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(htmlPage))
}
