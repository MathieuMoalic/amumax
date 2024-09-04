package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MathieuMoalic/amumax/engine"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/minio/selfupdate"
)

func doUpdate(tag string) {
	color.Green("Downloading version %s", tag)
	resp, err := http.Get(fmt.Sprintf("https://github.com/MathieuMoalic/amumax/releases/download/%s/amumax", tag))
	if err != nil {
		util.Log.PanicIfError(err)
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		color.Red("Error updating")
		color.Red(fmt.Sprint(err))
	}
	color.Green("Done.")
}

func getTags() (tags []string) {
	type Tag struct {
		Name string `json:"name"`
	}
	resp, err := http.Get("https://api.github.com/repos/mathieumoalic/amumax/tags")
	if err != nil {
		util.Log.PanicIfError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error: HTTP %v", resp.StatusCode)
	}

	tempTags := []Tag{}
	if err := json.NewDecoder(resp.Body).Decode(&tempTags); err != nil {
		util.Log.PanicIfError(err)
	}

	for _, tag := range tempTags {
		tags = append(tags, tag.Name)
	}
	return
}

func showUpdateMenu() {
	tags := getTags()

	// Create the prompt
	prompt := promptui.Select{
		Label: fmt.Sprintf("Current version : [%s] | Select the amumax version to update to", engine.VERSION),
		Items: tags,
		Size:  10,
	}

	// Run the prompt
	_, tag, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	doUpdate(tag)
}