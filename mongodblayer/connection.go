package mongodblayer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dblog struct {
	Action    string
	Timestamp int64
	Status    string
	TxID      string
}

var connStr string
var dbname string
var dbClient *mongo.Client

// Init initilize Connection
func Init(connectionString string, dbName string) {
	connStr = connectionString
	dbname = dbName

	if connStr == "" {
		fmt.Println("The connection string is not properly initialized")
		return
	}

	clientOptions := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	dbClient = client
}

// Close : cierra conexión, llamar despues de init con defer
func Close() {
	err := dbClient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Closing connection")
}

// TestConnection : prueba inicial
func TestConnection() {

	if connStr == "" {
		fmt.Println("The connection string is not properly initialized")
		return
	}

	err := dbClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexion exitosa a : " + connStr)
}

// InsertOneDocument inserta uno
func InsertOneDocument(col string, document interface{}) {
	collection := dbClient.Database(dbname).Collection(col)
	insertResult, err := collection.InsertOne(context.TODO(), document)
	if err == nil {
		insertOperationLog("inserting in "+dbname+"."+col+" : "+fmt.Sprintf("%v", insertResult.InsertedID), "success")
	} else {
		inserted := insertOperationLog("inserting in "+dbname+"."+col, "failure")
		fmt.Println("Operation Failure: log : " + inserted)
		log.Fatal(err)
	}
}

// DeleteOneDocument tbd
func DeleteOneDocument(col string, ID string) {
	collection := dbClient.Database(dbname).Collection(col)
	objID, erx := primitive.ObjectIDFromHex(ID)
	if erx != nil {
		inserted := insertOperationLog("deleting in "+dbname+"."+col, "failure")
		fmt.Println("Operation Failure: log : " + inserted)
		log.Fatal(erx)
		return
	}
	insertResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err == nil {
		insertOperationLog("deleting in "+dbname+"."+col+" : "+fmt.Sprintf("%v", insertResult.DeletedCount), "success")
	} else {
		inserted := insertOperationLog("deleting in "+dbname+"."+col, "failure")
		fmt.Println("Operation Failure: log : " + inserted)
		log.Fatal(err)
	}
}

// FindOneDocument  tbd
func FindOneDocument(col string, ID string) interface{} {
	var result interface{}
	collection := dbClient.Database(dbname).Collection(col)
	objID, erx := primitive.ObjectIDFromHex(ID)

	if erx != nil {
		inserted := insertOperationLog("deleting in "+dbname+"."+col, "failure")
		fmt.Println("Operation Failure: log : " + inserted)
		log.Fatal(erx)
		return nil
	}

	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
	if err == nil {
		fmt.Printf("Found a single document: %+v\n", result)
		return result
	} else {
		fmt.Println("Error Buscando")
		log.Fatal(err)
		return nil
	}
}

func insertOperationLog(action string, status string) string {
	collection := dbClient.Database(dbname).Collection("LOGS")
	inLog := dblog{action, time.Now().Unix(), status, "test"}
	insertResult, _ := collection.InsertOne(context.TODO(), inLog)
	return fmt.Sprintf("%v", insertResult.InsertedID)
}
