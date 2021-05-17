package impl

import (
	"github.com/lexiko/alexandros/pkg/config"
	"github.com/lexiko/plato/middleware"
	"github.com/lexiko/plato/models"
	"net/http"
)

type AlexandrosHandler struct {
	Config *config.AlexandrosConfig
}

// PingPong pongs the ping
func (a *AlexandrosHandler) pingPong(w http.ResponseWriter, req *http.Request) {
	pingPong := models.ResultModel{Result: "pong"}
	middleware.ResponseWithJson(w, pingPong)
}

func (a *AlexandrosHandler) searchWord(w http.ResponseWriter, req *http.Request) {
	//GET demokritos/_search
	//{
	//  "size": 10,
	//  "query": {
	//    "multi_match": {
	//      "query": "αγ",
	//      "type": "bool_prefix",
	//      "fields": [
	//        "greek",
	//        "greek._2gram",
	//        "greek._3gram"
	//      ]
	//    }
	//  }
	//}
}