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
	"log"
	"html/template"
	"github.com/Bifrost/server"
	"encoding/json"
	"github.com/Bifrost/toserver"
)

func init(){
	AddRoute("/db/detail",db_detail_controller)
	AddRoute("/db/tablelist",get_table_List_controller)
	AddRoute("/db/tablefields",get_table_fields_controller)
}

func db_detail_controller(w http.ResponseWriter,req *http.Request){
	defer func() {
		if err := recover();err!=nil{
			log.Println(err)
		}
	}()
	type dbDetail struct {
		TemplateHeader
		DbName string
		DataBaseList []string
		ToServerList string
		ChannelList map[int]*server.Channel
	}
	req.ParseForm()
	dbname := req.Form.Get("dbname")
	dbUri := server.GetDBObj(dbname).ConnectUri
	dbConn := DBConnect(dbUri)
	if dbConn == nil{
		return
	}
	defer dbConn.Close()
	DataBaseList := GetSchemaList(dbConn)
	var Result dbDetail
	b,_:=json.Marshal(toserver.GetToServerMap())
	Result = dbDetail{DataBaseList:DataBaseList,DbName:dbname,ToServerList: string(b),ChannelList:server.GetDBObj(dbname).ListChannel()}
	Result.Title = dbname + " - Detail - Bifrost"
	t, _ := template.ParseFiles("manager/template/db.detail.html","manager/template/header.html","manager/template/footer.html")
	t.Execute(w, Result)
}

func get_table_List_controller(w http.ResponseWriter,req *http.Request){
	defer func() {
		if err := recover();err!=nil{
			log.Println(err)
		}
	}()
	req.ParseForm()
	dbname := req.Form.Get("dbname")
	schema_name := req.Form.Get("schema_name")
	DBObj := server.GetDBObj(dbname)
	dbUri := DBObj.ConnectUri
	dbConn := DBConnect(dbUri)
	if dbConn == nil{
		return
	}
	defer dbConn.Close()
	type ResultType struct{
		TableName string
		AddStatus bool
	}
	var data []ResultType
	data = make([]ResultType,0)
	TableList := GetSchemaTableList(dbConn,schema_name)
	for _,tableName := range TableList{
		if DBObj.GetTable(schema_name,tableName) == nil{
			data = append(data,ResultType{TableName:tableName,AddStatus:false})
		}else{
			data = append(data,ResultType{TableName:tableName,AddStatus:true})
		}
	}
	b,_:=json.Marshal(data)
	w.Write(b)
}


func get_table_fields_controller(w http.ResponseWriter,req *http.Request){
	defer func() {
		if err := recover();err!=nil{
			log.Println(err)
		}
	}()
	req.ParseForm()
	dbname := req.Form.Get("dbname")
	schema_name := req.Form.Get("schema_name")
	table_name := req.Form.Get("table_name")
	DBObj := server.GetDBObj(dbname)
	dbUri := DBObj.ConnectUri
	dbConn := DBConnect(dbUri)
	if dbConn == nil{
		return
	}
	defer dbConn.Close()
	TableFieldsList := GetSchemaTableFieldList(dbConn,schema_name,table_name)
	b,_:=json.Marshal(TableFieldsList)
	w.Write(b)
}