package rethinkdb

import (
  "testing"
  "github.com/allen13/golerta/app/models"
  r "gopkg.in/dancannon/gorethink.v2"
)

//Integration test for alert CRUD operations
//docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb
func TestRethinkDB_CRUDAlert(t *testing.T) {
  db := getTestDB(t)

  alert := &models.Alert{
    Event: "cpu usage idle",
    Resource: "testServer01",
    Environment: "syd01",
    Severity: "CRITICAL",
    Origin: "consul-syd01",
  }
  alert.GenerateDefaults()

  //Create Alert
  id, err := db.CreateAlert(*alert)
  if err != nil {
    t.Fatal(err)
  }

  //Update Alert and check that it was actually updated
  alertUpdate := map[string]interface{}{
    "duplicateCount": r.Row.Field("duplicateCount").Add(1),
  }
  err = db.UpdateAlert(id, alertUpdate)
  if err != nil {
    t.Fatal(err)
  }

  //Get Alert and check that it is the same in addition to the duplciateCount field being updated
  dbAlert, err := db.GetAlert(id)
  if err != nil {
    t.Fatal(err)
  }
  if alert.Event != dbAlert.Event && alert.DuplicateCount == dbAlert.DuplicateCount + 1{
    t.Fatal("Failed to get back the correct alert")
  }

  //Delete Alert and check that it was deleted
  err = db.DeleteAlert(id)
  if err != nil {
    t.Fatal(err)
  }
  _, err = db.GetAlert(id)
  if err == nil {
    t.Fatal("Alert was not properly deleted")
  }
}

func getTestDB(t *testing.T)(db* RethinkDB){
  db = &RethinkDB{
    Address: "localhost:28015",
    Database: "alerta",
  }
  err := db.Connect()

  if err != nil{
    t.Fatal(err)
  }

  return db
}