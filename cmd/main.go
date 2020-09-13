package main

import (
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"net/http"
	"time"
)

func main() {
	//链路上报地址
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	//本地链路地址
	endpoint,err := zipkin.NewEndpoint("testSvc","localhost:9001")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v \n",err)
	}

	//zipkin链路上报追踪客户端
	tracer,err := zipkin.NewTracer(reporter,zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n",err)
	}

	//zipkin返回body大小中间件
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer,zipkinhttp.TagResponseSize(true),
	)
	//zipkin客户端
	client,err := zipkinhttp.NewClient(tracer,zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n",err)
	}

	//router
	router := mux.NewRouter()
	router.Use(serverMiddleware)
	router.Methods("GET").Path("/some_function").HandlerFunc(someFunc(client,"http://localhost:9001"))
	router.Methods("POST").Path("/other_function").HandlerFunc(otherFunc(client))

	err = http.ListenAndServe("localhost:9001",router)
	if err != nil {
		log.Fatal(err)
	}
}

func someFunc(client *zipkinhttp.Client,url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("some_function called with method: %s\n",r.Method)

		span := zipkin.SpanFromContext(r.Context())
		span.Tag("custome_key","some value")

		time.Sleep(25 * time.Millisecond)
		span.Annotate(time.Now(),"expensive_calc_done")

		newRequest,err := http.NewRequest("POST",url + "/other_function",nil)
		if err != nil {
			log.Printf("unable to create client /other_function : %+v\n",err)
			http.Error(w,err.Error(),500)
			return
		}

		ctx       := zipkin.NewContext(newRequest.Context(),span)
		newRequest = newRequest.WithContext(ctx)

		res,err := client.DoWithAppSpan(newRequest,"other_function")
		if err != nil {
			log.Printf("call to other_function returned error: %+v\n",err)
			http.Error(w,err.Error(),500)
			return
		}
		_ = res.Body.Close()
		return
	}
}

func otherFunc(client *zipkinhttp.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("other_function called with method: %s\n",r.Method)
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(200)
	}
}