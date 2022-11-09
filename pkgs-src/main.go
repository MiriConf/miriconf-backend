package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var walkFunc = func(filePath string, fileName fs.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("file read error: %v\n", err)
		return err
	}
	if fileName.IsDir() {
		fmt.Printf("not a file... skipping\n")
		return nil
	} else {
		if bytes.Contains([]byte(filePath), []byte(".nix")) {
			getDetails(filePath)
		} else {
			fmt.Println("not a nix file... skipping")
			return nil
		}
	}
	return nil
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("mongo URI is not specified, set with MONGO_URI environment variable")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	getPkgs("/home/pkgs/applications", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/desktops", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/development", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/games", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/misc", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/servers", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/shells", &wg)
	wg.Add(1)
	getPkgs("/home/pkgs/tools", &wg)
	wg.Wait()
	return
}

func getPkgs(pkgDir string, wgrp *sync.WaitGroup) {
	err := filepath.WalkDir(pkgDir, walkFunc)
	if err != nil {
		fmt.Printf("file read error: %v\n", err)
	}
	wgrp.Done()
}

func getDetails(filePath string) {
	pkgFile, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var pkgName []byte
	findPackageNameStart := bytes.Index(pkgFile, []byte("pname = \""))
	if findPackageNameStart == -1 {
		fmt.Println("no package name defined...")
		return
	} else {
		fromPName := pkgFile[findPackageNameStart:]
		findStartQuote := bytes.Index(fromPName, []byte("\""))
		fromQuoteStart := pkgFile[findPackageNameStart+findStartQuote+1:]
		findEndQuote := bytes.Index(fromQuoteStart, []byte("\""))
		pkgName = fromQuoteStart[:findEndQuote]
	}

	var pkgDesc []byte
	findPackageDescStart := bytes.Index(pkgFile, []byte("description = \""))
	if findPackageDescStart == -1 {
		pkgDesc = []byte("no description provided")
	} else {
		fromPDesc := pkgFile[findPackageDescStart:]
		findStartQuote := bytes.Index(fromPDesc, []byte("\""))
		fromQuoteStart := pkgFile[findPackageDescStart+findStartQuote+1:]
		findEndQuote := bytes.Index(fromQuoteStart, []byte("\""))
		pkgDesc = fromQuoteStart[:findEndQuote]
	}

	cmd := exec.Command("nix-env", "-qaA", "nixpkgs."+string(pkgName))

	stdOut, err := cmd.Output()
	if err != nil {
		fmt.Println("not a valid package: ", err)
		return
	}

	findVersionStart := bytes.LastIndex(stdOut, []byte("-"))
	pkgVer := stdOut[findVersionStart+1:]

	publishPackage(string(pkgName), string(pkgVer), string(pkgDesc))
}

func publishPackage(pkgName string, pkgVer string, pkgDesc string) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("miriconf").Collection("applications")

	var nameCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: pkgName}}).Decode(&nameCheck)
	if err != mongo.ErrNoDocuments {
		fmt.Println("application with this name already exists... skipping")
		return
	}

	doc := bson.D{{Key: "name", Value: pkgName}, {Key: "version", Value: pkgVer}, {Key: "description", Value: pkgDesc}}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("success: added package to mongo at %v\n", result.InsertedID)
	return
}
