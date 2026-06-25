package main

import (
	"grd/ui"
)

func main() {

	// releases, _ := utils.GetReleases("https://github.com/peass-ng/PEASS-ng")
	// asset := utils.GetAssetsFromRelease(releases[0])

	// fmt.Printf("Assets of: %s\n", releases[0].TagName)
	// for _, v := range asset {
	// 	fmt.Printf("File name: %s, download url: %s\n", v.Name, v.BrowserDownloadURL)
	// }
	//
	app := ui.New()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
