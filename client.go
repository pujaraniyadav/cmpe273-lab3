package main 

import (
    "net/http"
    "log"
    "io/ioutil"
    "strconv"
    "fmt"
    "os"
    "strings"
     "hash/fnv"
)


type Response struct{
	Key  int 		`json:"key"`
	Value string    `json:"value"`

}

func h(val int) uint32 {
	s := strconv.Itoa(val)

    h := fnv.New32a()
    h.Write([]byte(s))
    ret := h.Sum32()
    return ret
}

var key_value map[int] string

func main(){

	var port string
	var hash uint32
	if(len(os.Args)==1){
		log.Println("Incorrect number of entries.")
		os.Exit(1)
	}

	if(len(os.Args)==2 && os.Args[1] == "GET"){
		for port := 3000; port < 3003; port += 1 {
			url := fmt.Sprintf("http://localhost:%d/keys", port)
			get, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			data, err := ioutil.ReadAll(get.Body)
			get.Body.Close()
			log.Println("[ localhost:%d ]", port)
			log.Println("Key Value pairs are  : ", string(data))
		}
	}else if os.Args[1]== "GET" {
		request_string := os.Args[2]
		key := strings.Split(request_string,"/")
		key_integer,_ := strconv.Atoi(key[2])
		hash = h(key_integer) % 3
		if(hash == 0){
			port = "3000"
		}else if(hash == 1){
			port = "3001"
		}else {
			port = "3002"
		}
		url := fmt.Sprintf("http://localhost:%s/keys/%s",port,key[2])
		get_key, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}	
		data, err := ioutil.ReadAll(get_key.Body)
		get_key.Body.Close()
		log.Println("Key Value pairs are  : ", string(data))
	}else if os.Args[1] == "PUT" {
		req_string := os.Args[2]
		key_put := strings.Split(req_string,"/")
		key_int,_ := strconv.Atoi(key_put[2])
		hash = h(key_int) % 3
		if(hash == 0){
			port = "3000"
		}else if(hash == 1){
			port = "3001"
		}else {
			port = "3002"
		}
		put_url:= fmt.Sprintf("http://localhost:%s/keys/%s/%s",port,key_put[2],key_put[3])
		
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", put_url, nil)
		resp, _ := client.Do(req)
		//out, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		log.Println(" Response : ", 200)
	}
	
}	
