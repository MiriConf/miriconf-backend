package systems

import (
	"net/http"
)

func BuildTemplate(w http.ResponseWriter, r *http.Request) {
	// https://github.com/g-bolmida/miriconf-test/archive/d0df2c29316a008d63793e21d1d7a23f40ad4f16.zip

	//mongoURI := os.Getenv("MONGO_URI")
	//w.Header().Set("Content-Type", "application/json")
	//
	//headerToken := r.Header.Get("Authorization")
	//if headerToken == "" {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//
	//token, err := helpers.ValidateToken(headerToken)
	//if err != nil {
	//	if err == jwt.ErrSignatureInvalid {
	//		w.WriteHeader(http.StatusUnauthorized)
	//		return
	//	}
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//
	//if !token.Valid {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//
	//status, teamID := helpers.GetRequestID("team", r, w)
	//if status == 1 {
	//	return
	//}
	//
	//client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	if err = client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//coll := client.Database("miriconf").Collection("teams")
	//
	//var result teams.Team
	//err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: teamID}}).Decode(&result)
	//if err != nil {
	//	if err == mongo.ErrNoDocuments {
	//		if err != nil {
	//			error := helpers.ErrorMsg("no team matching id requested")
	//			w.Write(error)
	//			helpers.EndpointError("no team matching id requested", r)
	//			return
	//		}
	//	}
	//	log.Fatal(err)
	//}
	//
	//directory := "/mnt/data/" + result.Name
	//_, err = os.Stat(directory)
	//if os.IsNotExist(err) {
	//	panic("directory does not exist... build failed")
	//}
	//
	//openRepo, err := git.PlainOpen(directory)
	//if err != nil {
	//	panic(err)
	//}
	//
	//repoTree, err := openRepo.Worktree()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("pulling latest commits from %v\n", result.SourceRepo)
	//err = repoTree.Pull(&git.PullOptions{Auth: &githttp.BasicAuth{Username: "null", Password: result.SourcePAT}, RemoteName: "origin", Progress: os.Stdout})
	//if err != nil && err == git.NoErrAlreadyUpToDate {
	//	fmt.Println("no changes to pull")
	//} else if err != nil {
	//	panic(err)
	//}
	//
	//repoTree, err = openRepo.Worktree()
	//if err != nil {
	//	panic(err)
	//}
	//
	//generated := `users.users.gbolmida = {
	//	isNormalUser = true;
	//	description = "George";
	//	shell = pkgs.fish;
	//	extraGroups = [ "networkmanager" "wheel" "docker" ];
	//	packages = with pkgs; [
	//	  firefox
	//	  chromium
	//	  unstable.vivaldi
	//	  hugo
	//	  gimp
	//	  vim
	//	  wget
	//	  nano
	//	  tree
	//	  htop
	//	  vscode
	//	  code-server
	//	];
	//  };`
	//
	//filename := filepath.Join(directory, "generated.nix")
	//err = ioutil.WriteFile(filename, []byte(generated), 0644)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = repoTree.Add("generated.nix")
	//if err != nil {
	//	panic(err)
	//}
	//
	//repoStatus, err := repoTree.Status()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(repoStatus)
	//
	//newCommit, err := repoTree.Commit("test commit", &git.CommitOptions{Author: &object.Signature{Name: "miriconf-bot", Email: "miriconf-bot@" + hostName, When: time.Now()}})
	//if err != nil {
	//	panic(err)
	//}
	//
	//obj, err := openRepo.CommitObject(newCommit)
	//if err != nil {
	//	panic(err)
	//}
	//
	//helpers.SuccessLog(r)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode("build commited successfully with message: " + obj.Message)
}
