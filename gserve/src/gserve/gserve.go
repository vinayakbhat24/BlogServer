package main


import (
    	"fmt"
    	"log"
    	"net/http"
        "html/template"
  	"io/ioutil"
  	"encoding/json"
	"bytes"
       	"path/filepath"
        "os"
 	"reflect"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type Profile struct {
  DocID    string
  Document string
  Value string
  Author string
  MetVal string
}




func handle(rw http.ResponseWriter, request *http.Request) {
    println("--->")
    fmt.Println("method:", request.Method) //get request method

    if request.Method == "GET" {


        client := &http.Client{}



//	body := strings.NewReader(`<Scanner batch="10"/>`)

//	rsp, err := http.Post("http://hbase:8080/se2:library/scanner/", "application/json", body)
//	if err != nil {
//		panic(err)
//	}
//	defer rsp.Body.Close()

//	body_byte, err := ioutil.ReadAll(rsp.Body)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(body_byte))



        req, err := http.NewRequest("GET", "http://hbase:8080/se2:library/*", nil)
        if err != nil {
                log.Fatalln(err)
        }

        req.Header.Set("Accept", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                log.Fatalln(err)
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatalln(err)
        }

        log.Println(string(body))

	var encodedRows EncRowsType
	json.Unmarshal(body, &encodedRows)

	val,err :=encodedRows.decode()

	decodedJSON, _ := json.Marshal(val)
	println("decoded:", string(decodedJSON))

	var unencodedRows RowsType
	json.Unmarshal(decodedJSON, &unencodedRows)


	//var decodedrows RowsType
	//json.Unmarshal(decodedJSON, &decodedrows)

	//http.ServeFile(rw, request, request.URL.Path[1:])

	//fmt.Fprintf(rw, string(decodedJSON))

	//Print title SE2 LIBRARY
	cwd, _ := os.Getwd()
        fmt.Println( filepath.Join( cwd, "title.html" ) )
	fpTitle:= filepath.Join( cwd, "title.html" );
  	tmplTitle, err := template.ParseFiles(fpTitle)
  		if err != nil {
    			http.Error(rw, err.Error(), http.StatusInternalServerError)
    			return
  		}
  			if err := tmplTitle.Execute(rw, tmplTitle); err != nil {
    			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}


			j:=0
	for _, value := range unencodedRows.Row {
			s := reflect.ValueOf(&value).Elem()
			typeOfT := s.Type()
			for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			fmt.Print("key   : ", typeOfT.Field(i).Name, "\n")
			fmt.Print("value : ", f.Interface(), "\n")
			fmt.Print("\n")
			fmt.Print("\n")
			fmt.Print("\n")
			fmt.Print("\n")

			}
			//Temporary soution
			Key0 := unencodedRows.Row[j].Key
			doc := unencodedRows.Row[j].Cell[0].Column
			docVal := unencodedRows.Row[j].Cell[0].Value
			author := unencodedRows.Row[j].Cell[1].Column
			metVal := unencodedRows.Row[j].Cell[1].Value

		  	profile := Profile{Key0,doc,docVal,author,metVal}

			fmt.Println( filepath.Join( cwd, "index.html" ) )
			fp:= filepath.Join( cwd, "index.html" );

		  	tmpl, err := template.ParseFiles(fp)

		  	if err != nil {
		    			http.Error(rw, err.Error(), http.StatusInternalServerError)
		    			return
		  		}

		  	if err := tmpl.Execute(rw, profile); err != nil {
		    		http.Error(rw, err.Error(), http.StatusInternalServerError)
				}

			j++;
		}
	fmt.Fprintf(rw, "Proudly served by  %s", os.Getenv("NAME"))

    } else if request.Method == "POST" {
    	println("Inside Post  method")
       	unencodedJSON, err := ioutil.ReadAll(request.Body)
    	if err != nil {
        	panic(err.Error())
    		}
   	println("unencoded:", string(unencodedJSON))

	var unencodedRows RowsType
	json.Unmarshal(unencodedJSON, &unencodedRows)

	encodedRows := unencodedRows.encode()
	// convert encoded Go objects to JSON
	encodedJSON, _ := json.Marshal(encodedRows)
	println("encoded:", string(encodedJSON))

	b := bytes.NewReader(encodedJSON)
	rsp, err := http.Post("http://hbase:8080/se2:library/Col1", "application/json", b)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	body_byte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body_byte))

    }
}


func registerToZookeeper() {

	timer1 := time.NewTimer(time.Second * 20)
 	<-timer1.C

	conn, _, err1 := zk.Connect([]string{"hbase:2181"}, 6*time.Second)
	if err1 != nil {
		fmt.Printf(" connect zk error: %s ", err1)
	} 	
	if err1 != nil {
		fmt.Printf(" connect zk error: %s ", err1)
	} 	


	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	path, err2 := conn.Create("/go_servers", []byte("data"), flags, acl)


	if err2 != nil {
		fmt.Printf(" connect zk error: %s ", err2)
	} 	
	fmt.Printf("create: %+v\n", path)

	path, err3 := conn.Create("/go_servers/"+os.Getenv("NAME"), []byte("data"), flags, acl);
	if err3 != nil {
		fmt.Printf(" connect zk error: %s ", err3)
	} 	
	fmt.Printf("create: %+v\n", path)


}


func main() {
    registerToZookeeper()
    http.HandleFunc("/", handle) // setting router rule
    err := http.ListenAndServe(":9090", nil) // setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

