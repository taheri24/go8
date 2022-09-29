package report

import (
	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/server"
	"github.com/gmhafiz/go8/internal/utility/respond"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ReportRouter(r chi.Router) {
	db := server.Serv.DB()
	sqlQuery := `WITH C1 AS (
SELECT  instrument_id,MAX(date_en) as date_en
FROM Trade  GROUP BY instrument_id
 )
SELECT t.instrument_id,t.date_en,t.low,t.high,t.open,t.close FROM Trade t
INNER JOIN C1 ON t.instrument_id=c1.instrument_id AND t.date_en=c1.date_en
`
	r.Get("/report1", func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		rows, err := db.QueryxContext(ctx, sqlQuery)
		if err != nil {
			respond.Error(writer, http.StatusInternalServerError, err)

			return
		}
		userResponse := make([]gen.Trade, 0, 100)
		for rows.Next() {
			currTrade := gen.Trade{}
			rows.StructScan(&currTrade)
			userResponse = append(userResponse, currTrade)
		}
		respond.Json(writer, http.StatusOK, userResponse)
	})
}
