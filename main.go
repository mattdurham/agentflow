package main

import (
	"agentflow/config"
	"agentflow/orchestration"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	as := actor.NewActorSystem()
	root := actor.NewRootContext(as, nil)
	cfgStr, err := ioutil.ReadFile("/Users/mdurham/Utils/agent_flow_configs/agent_flow_prom.yml")
	if err != nil {
		panic(err)
	}
	cfg := &config.Config{}
	err = yaml.Unmarshal(cfgStr, cfg)
	if err != nil {
		panic(err)
	}
	orch := orchestration.NewOrchestrator(*cfg)
	err = orch.StartActorSystem(as, root)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/mermaid", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(orch.GenerateMermaid()))
	})

	router.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		bb, _ := yaml.Marshal(orch.NodeList())
		w.Write(bb)
	})

	router.HandleFunc("/nodes/{name}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		key := vars["name"]
		writer.Write(orch.GetNodeStatus(key))

	})
	log.Fatal(http.ListenAndServe(":54321", router))
}
