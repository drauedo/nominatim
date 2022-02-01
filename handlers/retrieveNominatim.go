package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"redis/data"
	"time"

	"github.com/go-redis/redis/v8"
)

func (n *API) RetrieveNominatim(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")
	n.l.Println("Retrieving data from: ", q)
	data, cache, err := getData(q, r.Context(), *n)

	if err != nil {
		n.l.Println("There was an error with the query")
		n.l.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := APIresponse{
		Cache: cache,
		Data:  data,
	}
	n.l.Println(resp)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		n.l.Println("There was an error enconding the response")
		n.l.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func getData(q string, ctx context.Context, a API) ([]data.NominatimResponse, bool, error) {

	cachedQuery, err := a.cache.Get(ctx, q).Result()

	if err == redis.Nil {
		escapedQ := url.PathEscape(q)
		address := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapedQ)
		resp, err := http.Get(address)

		if err != nil {
			return nil, false, err
		}

		data := make([]data.NominatimResponse, 0)

		err = json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			return nil, false, err
		}

		b, err := json.Marshal(data)

		if err != nil {
			return nil, false, err
		}
		a.l.Println("Caching solution")
		err = a.cache.Set(ctx, q, bytes.NewBuffer(b).Bytes(), time.Second+30).Err()

		return data, false, nil

	} else if err != nil {
		a.l.Println("error calling redis: %v\n", err)
		panic(err)
	} else {
		data := make([]data.NominatimResponse, 0)

		err := json.Unmarshal(bytes.NewBufferString(cachedQuery).Bytes(), &data)

		if err != nil {
			return nil, false, err
		}
		return data, true, nil
	}

}
