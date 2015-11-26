package main 

import (
	"github.com/julienschmidt/httprouter"
    "net/http"
    "log"
    "strconv"
    "encoding/json"
    "strings"
)

var server1_kvalue map[int] string
var server2_kvalue map[int] string
var server3_kvalue map[int] string

type Response struct{
  Key  int 		`json:"key"`
  Value string    `json:"value"`
}

func main() {

  server1_kvalue = make(map[int] string)
  server2_kvalue = make(map[int] string)
  server3_kvalue = make(map[int] string)

  //server 1 code
  go func() {
    mux1 := httprouter.New()	
    mux1.PUT("/keys/:id/:value",put)
    mux1.GET("/keys/:id",get)
    mux1.GET("/keys",getall)
    server := http.Server{
      Addr: "127.0.0.1:3000",
      Handler: mux1,
    }
    server.ListenAndServe()
  } ()
    
  //server 2 code
  go func(){
    mux2 := httprouter.New()
    mux2.PUT("/keys/:id/:value",put)
    mux2.GET("/keys/:id",get)
    mux2.GET("/keys",getall)
    server2 := http.Server{
      Addr:        "127.0.0.1:3001",
      Handler: mux2,
    }
    server2.ListenAndServe()
  }()
  
    //server 3 code
	mux3 := httprouter.New()
  mux3.PUT("/keys/:id/:value",put)
  mux3.GET("/keys/:id",get)
  mux3.GET("/keys",getall)
  server3 := http.Server{
          Addr:        "127.0.0.1:3002",
          Handler: mux3,
  }
  server3.ListenAndServe()


}

  //PUT function
func put(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
	key := p.ByName("id")
	value := p.ByName("value")
  var port []string
	key_int, _ := strconv.Atoi(key)
  port = strings.Split(req.Host,":")
  if(port[1]=="3000"){
      server1_kvalue[key_int] = value    
  } else if (port[1]=="3001"){
      server2_kvalue[key_int] = value 

  } else{
      server3_kvalue[key_int] = value  

  }
}

  //GET function
func get(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

	key := p.ByName("id")
	key_int, _ := strconv.Atoi(key)
  var port []string
	var response Response
  port = strings.Split(req.Host,":")
  
  if(port[1]=="3000"){
      response.Key = key_int
      response.Value = server1_kvalue[key_int]   
  } else if (port[1]=="3001"){
      response.Key = key_int
      response.Value = server2_kvalue[key_int] 
  } else{
      response.Key = key_int
      response.Value = server3_kvalue[key_int] 
  }
  
  payload, err := json.Marshal(response)  
  if err != nil {
    http.Error(rw,"Bad Request" , http.StatusInternalServerError)
     return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(payload)
}

  //GETALL function
func getall(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

	var response []Response
	var key_pair Response
  var port []string
  port = strings.Split(req.Host,":")

  if(port[1]=="3000") {
    for key, value := range server1_kvalue {
      key_pair.Key = key
      key_pair.Value = value
      response = append(response, key_pair)
    }    
  } else if (port[1]=="3001"){
    for key, value := range server2_kvalue {
      key_pair.Key = key
      key_pair.Value = value
      response = append(response, key_pair)
    } 
  } else{
    for key, value := range server3_kvalue {
      key_pair.Key = key
      key_pair.Value = value
      response = append(response, key_pair)
    } 
  }
	
      //error handling
  payload, err := json.Marshal(response)  
  if err != nil {
    http.Error(rw,"Bad Request" , http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(payload)
}
