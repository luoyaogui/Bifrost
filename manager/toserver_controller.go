/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package manager
import (
	"net/http"
	"github.com/Bifrost/toserver"
	"github.com/Bifrost/toserver/driver"
	"fmt"
	"html/template"
)

func init()  {
	AddRoute("/toserver/add",toserver_add_controller)
	AddRoute("/toserver/del",toserver_del_controller)
	AddRoute("/toserver/list",toserver_list_controller)
	AddRoute("/toserver/check_uri",toserver_checkuri_controller)
}

func toserver_add_controller(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	toServerName := req.Form.Get("toserverkey")
	Type := req.Form.Get("type")
	Notes := req.Form.Get("notes")
	ConnUri := req.Form.Get("connuri")
	if toServerName == "" || Type == "" || ConnUri==""{
		w.Write(returnResult(false,"toserverkey,type,connuri muest be not empty"))
		return
	}
	toserver.SetToServerInfo(toServerName, Type, ConnUri, Notes)
	w.Write(returnResult(true,"success"))
}

func toserver_list_controller(w http.ResponseWriter,req *http.Request){
	defer func() {
		if err := recover();err!=nil{
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()
	req.ParseForm()
	type toServerInfo struct {
		TemplateHeader
		ToServerList map[string]*toserver.ToServer
		Drivers []map[string]string
	}
	var data toServerInfo
	data = toServerInfo{ToServerList: toserver.ToServerMap,Drivers:driver.Drivers()}
	data.Title = "ToServer List - Bifrost"
	t, _ := template.ParseFiles("manager/template/toserver.list.html","manager/template/header.html","manager/template/footer.html")
	t.Execute(w, data)

}

func toserver_del_controller(w http.ResponseWriter,req *http.Request){
	defer func() {
		if err := recover();err!=nil{
			w.Write([]byte(fmt.Sprint(err)))
		}
	}()
	req.ParseForm()
	toServerName := req.Form.Get("toserverkey")
	toserver.DelToServerInfo(toServerName)
	w.Write(returnResult(true,"success"))
}

func toserver_checkuri_controller(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	Type := req.Form.Get("type")
	ConnUri := req.Form.Get("connuri")
	if Type == "" || ConnUri==""{
		w.Write(returnResult(false,"type,connuri must be not empty"))
		return
	}
	err := driver.CheckUri(Type,ConnUri)
	if err !=nil{
		w.Write(returnResult(false,err.Error()))
		return
	}
	w.Write(returnResult(true,"success"))
}