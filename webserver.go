package arachne

import (
	"fmt"
	"github.com/bmeg/arachne/ophion"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

func HandleError(w http.ResponseWriter, req *http.Request, error string, code int) {
	fmt.Println(error)
	fmt.Println(req.URL)
	debug.PrintStack()
	http.Error(w, error, code)
}

//MarshalClean is a shim class to 'fix' outgoing streamed messages
//in the default implementation, grpc-gateway wraps the individual messages
//of the stream with a {"result" : <value>}. The cleaner idendifies that and
//removes the wrapper
type MarshalClean struct {
	m jsonpb.Marshaler
}

func (self *MarshalClean) ContentType() string {
	return "application/json"
}

func (self *MarshalClean) Marshal(v interface{}) ([]byte, error) {
	log.Printf("Marshal: %#v", v)
	if x, ok := v.(map[string]proto.Message); ok {
		o, err := self.m.MarshalToString(x["result"])
		return []byte(o), err
	}
	o, err := self.m.MarshalToString(v.(proto.Message))
	return []byte(o), err
}

func (self *MarshalClean) NewDecoder(r io.Reader) runtime.Decoder {
	return json.NewDecoder(r)
}

func (self *MarshalClean) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

func (self *MarshalClean) Unmarshal(data []byte, v interface{}) error {
	return jsonpb.UnmarshalString(string(data), v.(proto.Message))
}

func StartHttpProxy(rpcPort string, httpPort string, contentDir string) {
	//setup RESTful proxy
	marsh := MarshalClean{m:jsonpb.Marshaler{OrigName:true}}
	grpcMux := runtime.NewServeMux(runtime.WithMarshalerOption("*", &marsh))
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	log.Println("HTTP proxy connecting to localhost:" + rpcPort)
	err := ophion.RegisterQueryHandlerFromEndpoint(ctx, grpcMux, "localhost:"+rpcPort, opts)
	if err != nil {
		fmt.Println("Register Error", err)

	}
	r := mux.NewRouter()

	runtime.OtherErrorHandler = HandleError
	// Routes consist of a path and a handler function
	r.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(contentDir, "index.html"))
		})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(contentDir))))

	r.PathPrefix("/v1/").Handler(grpcMux)
	log.Printf("HTTP API listening on port: %s\n", httpPort)
	http.ListenAndServe(":"+httpPort, r)
}
