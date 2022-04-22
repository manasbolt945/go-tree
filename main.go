package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var map_a map[string]int

type metrics struct {
	webrequests int
	timespent   int
}
type dimension struct {
	device  string
	country string
}
type values struct {
	dim dimension
	met metrics
}
type dimu struct {
	key string
	val string
}
type metu struct {
	key string
	val int
}
type Node struct {
	dim [2]dimu
	met [2]metu
}

type RootNode struct {
	val   metrics
	child []*CountryNode
}

var root *RootNode

type CountryNode struct {
	country string
	met     metrics
	child   []*DeviceNode
}

type DeviceNode struct {
	device string
	met    metrics
}

func GetCountryNode(dim dimension, met metrics) *CountryNode {
	var me *CountryNode = &CountryNode{}
	me.country = dim.country
	me.met = met
	me.child = make([]*DeviceNode, 0)
	return me
}

func getDeviceNode(dim dimension, met metrics) *DeviceNode {
	var me *DeviceNode = &DeviceNode{}
	me.device = dim.device
	me.met = met
	return me
}

func (this *CountryNode) AddtoCountry(dim dimension, met metrics) {
	var me *DeviceNode = getDeviceNode(dim, met)
	var t *CountryNode = FindCountryNode(dim.country)
	UpdateCountryNode(met, t)
	if dim.device != "" {
		if FindDeviceNode(dim.device, t) != nil {
			var a *DeviceNode = FindDeviceNode(dim.device, t)
			UpdateDeviceNode(met, a)
			return
		}
	}
	this.child = append(this.child, me)
}
func FindCountryNode(country string) *CountryNode {
	var i int = 0
	for i < len(root.child) {
		if country == root.child[i].country {
			return root.child[i]
		}
	}
	return nil
}
func FindDeviceNode(device string, node *CountryNode) *DeviceNode {
	var i int = 0
	for i < len(node.child) {
		if device == node.child[i].device {
			return node.child[i]
		}
	}
	return nil
}
func UpdateCountryNode(met metrics, node *CountryNode) {
	//var me *CountryNode = FindCountryNode(dim.country)
	node.met.timespent += met.timespent
	node.met.webrequests += met.webrequests
}

func UpdateDeviceNode(met metrics, node *DeviceNode) {
	node.met.timespent += met.timespent
	node.met.webrequests += met.webrequests
}
func (this *RootNode) AddtoRoot(dim dimension, met metrics) {
	var t *CountryNode = GetCountryNode(dim, met)
	map_a[dim.country] = 1
	root.child = append(this.child, t)
}
func UpdateRootNode(met metrics) {
	root.val.timespent += met.timespent
	root.val.webrequests += met.webrequests

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/insert", insertNode).Methods("POST")
	router.HandleFunc("/v1/query", returnNode).Methods("POST")
	fmt.Println("1. For insert \n2. For Query")
	http.ListenAndServe(":8080", router)
}

func insertNode(w http.ResponseWriter, r *http.Request) {
	var abc Node
	var fgh values
	json.NewDecoder(r.Body).Decode(&fgh)
	var dims [2]dimu
	dims[0].key = abc.dim[0].key
	dims[0].val = abc.dim[0].val
	dims[1].key = abc.dim[1].key
	dims[1].val = abc.dim[1].val
	var mets [2]metu
	mets[0].key = abc.met[0].key
	mets[0].val = abc.met[0].val
	mets[1].key = abc.met[0].key
	mets[1].key = abc.met[0].key
	var d dimension
	var m metrics
	d.country = fgh.dim.country
	d.device = fgh.dim.device
	m.timespent = fgh.met.timespent
	m.webrequests = fgh.met.webrequests
	/*if dims[0].key == "country" {
		d.country = dims[0].val
	}
	if dims[0].key == "device" {
		d.device = dims[0].val
	}
	if dims[1].key == "country" {
		d.country = dims[1].val
	}
	if dims[1].key == "device" {
		d.device = dims[1].val
	}
	m.timespent += mets[1].val
	m.webrequests += mets[0].val
	*/
	UpdateRootNode(m)
	_, f := map_a[d.country]
	if f {
		var node *CountryNode = FindCountryNode(d.country)
		node.AddtoCountry(d, m)
	} else {
		root.AddtoRoot(d, m)
	}
}
func returnNode(w http.ResponseWriter, r *http.Request) {
	var gg dimension
	json.NewDecoder(r.Body).Decode(&gg)
	var node *CountryNode = FindCountryNode(gg.country)
	var node2 *DeviceNode = FindDeviceNode(gg.device, node)
	var vars metrics
	vars.timespent = node.met.timespent
	vars.webrequests = node.met.webrequests
	json.NewEncoder(w).Encode(&vars)
}
