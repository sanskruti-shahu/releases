package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	const githubApiEndpoint = "https://api.github.com/repos/opentofu/opentofu/releases"

	mainPageContent := "<ul>\n<li>\n<a href=\"../\">../</a></li>\n"

	response, err := http.Get(githubApiEndpoint)
	if err != nil {
		fmt.Println("Error fetching releases: ", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	// fmt.Println(string(body))
	fmt.Println("Data fetched")

	// file, err := os.Open("releases.json")
	// if err != nil {
	// 	fmt.Println("Error opening JSON file:", err)
	// 	return
	// }
	// defer file.Close()

	// body, err := io.ReadAll(file)
	// if err != nil {
	// 	fmt.Println("Error reading JSON file:", err)
	// 	return
	// }
	releases := []map[string]interface{}{}
	if err := json.Unmarshal(body, &releases); err != nil {
		fmt.Println("Error unmarshalling JSON: ", err)
		return
	}

	for _, release := range releases {
		versionTrimmed := release["name"].(string)[1:]
		version := release["name"].(string)
		// child page
		path := versionTrimmed + "/"
		if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
			fmt.Println("Error creating directory: ", err)
			return
		}
		childPageContent := "<ul>\n<li>\n<a href=\"../\">../</a></li>\n"
		if assets, ok := release["assets"].([]interface{}); ok {
			for _, asset := range assets {
				if assetMap, ok := asset.(map[string]interface{}); ok {
					fileName := assetMap["name"]
					childPageContent += "<li>\n"
					childPageContent += fmt.Sprintf("<a href=\"https://github.com/opentofu/opentofu/releases/download/%s/%s\">%s</a>\n", version, fileName, fileName)
					childPageContent += "</li>\n"
				}
			}
		}

		childPageContent += "</ul>\n"
		if err := os.WriteFile(path+"/index.html", []byte(childPageContent), 0644); err != nil {
			fmt.Println("Error writing child page: ", err)
			return
		}
		// main page
		mainPageContent += "<li>\n"
		mainPageContent += fmt.Sprintf("<a href=\"./%s\">tofu_%s</a>\n", path, versionTrimmed)
		mainPageContent += "</li>\n"
	}
	mainPageContent += "</ul>\n"

	if err := os.WriteFile("index.html", []byte(mainPageContent), 0644); err != nil {
		fmt.Println("Error writing main page: ", err)
		return
	}

	fmt.Println("All files written and now server started")

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}
