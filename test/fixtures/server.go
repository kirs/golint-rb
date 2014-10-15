package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var databaseConfigPath, envName, bindAddress string

func helloHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "http://github.com", http.StatusFound)
}

func fetchItemHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	req.ParseForm()

	force := len(req.FormValue("force")) > 0 && req.UserAgent() == secret_useragent

	item_id := vars["item_id"]
	if item_id != "" {
		var some_item SomeItem
		var err error

		if !force {
			err = dbmap.SelectOne(&some_item, "select id from caches where item_id=$1 limit 1", item_id)

			log.Printf("found id: %d, %s", some_item.Id, err)
			if err == nil {
				respondWithJson(w, &some_item)
				return
			}
		}

		some_item, err = FetchSomeItem(item_id)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%s", err), 422)
			return
		}

		if force {
			log.Println("force fetch, removing cached item")
			dbmap.Exec("delete from some_items_cache where item_id=$1", item_id)
		}

		err = dbmap.Insert(&some_item)
		if err != nil {
			log.Printf("insert error: %s", err)
		}

		respondWithJson(w, &some_item)
	} else {
		http.Error(w, "", 404)
	}
}

func revisionHandler(w http.ResponseWriter, req *http.Request) {
	revision_file := "/home/some/some/current/REVISION"

	response, err := ioutil.ReadFile(revision_file)
	if err != nil {
		log.Printf("failed to read %s: %s", revision_file, err)
	}

	w.Write(response)
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	query := req.Form.Get("query")
	category_id := req.Form.Get("category_id")
	aspect_filters := BuildAspects(req.Form["aspects[][name]"], req.Form["aspects[][value]"])

	if query != "" || category_id != "" {
		per_page, err := strconv.ParseInt(req.Form.Get("per_page"), 10, 64)
		if err != nil {
			per_page = 10
		}

		page_num, err := strconv.ParseInt(req.Form.Get("page"), 10, 64)
		if err != nil {
			page_num = 1
		}

		search_result, err := SearchSomeItems(query, category_id, aspect_filters, per_page, page_num)
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("%s", err), 422)
			return
		}

		// replace with respond json
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		b, err := json.Marshal(search_result)
		if err != nil {
			panic(err)
		}

		w.Write(b)
	}
}

var (
	dbmap            *gorp.DbMap
	secret_useragent string
)

const (
	ROUTE_PREFIX = "/api"
)

func main() {
	flag.Parse()

	dbmap = initDb()
	defer dbmap.Db.Close()

	useragent := []byte("SomeItemSyncer")
	secret_useragent = base64.StdEncoding.EncodeToString(useragent)

	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf("%s/", ROUTE_PREFIX), helloHandler)
	r.HandleFunc(fmt.Sprintf("%s/revision", ROUTE_PREFIX), revisionHandler)
	r.HandleFunc(fmt.Sprintf("%s/items/{item_id}", ROUTE_PREFIX), fetchItemHandler)
	r.HandleFunc(fmt.Sprintf("%s/search", ROUTE_PREFIX), searchHandler)
	http.Handle("/", r)

	fmt.Printf("Starting API on %s\n", bindAddress)

	err := http.ListenAndServe(bindAddress, r)
	if err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}
