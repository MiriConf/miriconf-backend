package database

//func GetTeams(client *mongo.Client) {
//  teamsColl := client.Database("miriconf").Collection("teams")
//
//  cursor, err := teamsColl.Find(context.TODO(), bson.D{})
//  if err != nil {
//    panic(err)
//  }
//
//  for cursor.Next(context.TODO()) {
//    var result bson.D
//    if err := cursor.Decode(&result); err != nil {
//      log.Fatal(err)
//    }
//	u, err := json.Marshal(Team[result])
//	if err != nil {
//		panic(err)
//	}
//
//  }
//  if err := cursor.Err(); err != nil {
//    log.Fatal(err)
//  }
//
//
//}