package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MathieuMoalic/amumax/src/version"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/minio/selfupdate"
)

func doUpdate(tag string) {
	color.Green("Downloading version %s", tag)
	resp, err := http.Get(fmt.Sprintf("https://github.com/MathieuMoalic/amumax/releases/download/%s/amumax", tag))
	if err != nil {
		color.Red(fmt.Sprintf("Error downloading the binary: %s", err))
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		color.Red(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		os.Exit(1)
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
		color.Red(fmt.Sprintf("Error fetching tags: %s", err))
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		color.Red(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		os.Exit(1)
	}

	tempTags := []Tag{}
	if err := json.NewDecoder(resp.Body).Decode(&tempTags); err != nil {
		color.Red(fmt.Sprintf("Error decoding tags: %s", err))
		os.Exit(1)
	}

	for _, tag := range tempTags {
		tags = append(tags, tag.Name)
	}
	return
}

func ShowUpdateMenu() {
	tags := getTags()

	// Create the prompt
	prompt := promptui.Select{
		Label: fmt.Sprintf("Current version : [%s] | Select the amumax version to update to", version.VERSION),
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
