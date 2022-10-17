package boxofficemojo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type BoxOfficeMojo struct {
	//IMDB Identifier
	Id string

	//URL to BoxOfficeMojo for the Movie
	Url string

	//# of US Theaters
	UsaScreens int
	// URL to the page with US Theater information
	UsaScreensUrl string

	//# of UK Theaters
	UkScreens int
	// URL to the page with UK Theater information
	UKScreensUrl string
}

// Searches and Parse BoxOfficeMojo to retrieve the # of theaters for a movie based on its IMDB identifier.
// Will return an error when the movie does not have screen data and screens will be defaulted to zero as a safe value.
func Search(id string) (BoxOfficeMojo, error) {

	moviePageUrl := getBomUrl(id)
	// fmt.Println("Movie Page url ", moviePageUrl)

	// Retrieve Page #1 - Movie Details
	rg, err := getReleaseGroupUrl(moviePageUrl)
	if err != nil {
		return BoxOfficeMojo{
			Id:  id,
			Url: moviePageUrl,
		}, err
	}
	// fmt.Println("Release group url ", rg)

	// Retrieve Page #2 - Summary of Original Release - Extract Domestic page URL
	domesticUrl, err := getDomesticUrl(rg)
	if err != nil {
		return BoxOfficeMojo{
			Id:  id,
			Url: moviePageUrl,
		}, err
	}
	// fmt.Println("Domestic url ", domesticUrl)

	theaters, err := extractScreens(domesticUrl)
	if err != nil {
		return BoxOfficeMojo{}, err
	}
	fmt.Printf("Domestic Theaters %d ", theaters)

	return BoxOfficeMojo{
		Id:            id,
		Url:           moviePageUrl,
		UsaScreens:    theaters,
		UsaScreensUrl: domesticUrl,
		UkScreens:     0, // Not Implemented
		UKScreensUrl:  "",
	}, nil
}

func getBomUrl(id string) string {
	return fmt.Sprintf("https://www.boxofficemojo.com/title/%s", id)
}

// Retrieve and parse web page for the movie to extract the URL to the Box Office details
func getReleaseGroupUrl(moviePageUrl string) (string, error) {

	res, err := http.Get(moviePageUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading contents from %s, Error: %s", moviePageUrl, err)
	}
	contentStr := string(content[:])
	rgx := regexp.MustCompile(`(?im)(/releasegroup/[a-z0-9]+/)`)
	if rgx.MatchString(contentStr) {
		match := rgx.FindAllStringSubmatch(contentStr, -1)
		urlComponent := match[0][1]

		return fmt.Sprintf("https://www.boxofficemojo.com%s", urlComponent), nil
	}

	return "", fmt.Errorf("unable to find a url matching /releasegroup/xxxxxxxx/ on url %s", moviePageUrl)
}

// Parse web page to extract the URL to the domestic stats for the original release
func getDomesticUrl(url string) (string, error) {

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading contents from %s, Error: %s", url, err)
	}
	//Domestic - https://www.boxofficemojo.com/release/rl2500036097/?ref_=bo_gr_rls
	contentStr := string(content[:])
	rgx := regexp.MustCompile(`(?im)(value="/release\/(.*?)/"[^>]*?>Domestic</option>)`)
	if rgx.MatchString(contentStr) {
		match := rgx.FindAllStringSubmatch(contentStr, -1)
		urlComponent := match[0][2]

		return fmt.Sprintf("https://www.boxofficemojo.com/release/%s", urlComponent), nil
	}

	return "", fmt.Errorf("unable to find a url matching /release/xxxxxxxx/ for domestic stats on url %s", url)

}

// Parse HTML page to extract the # of screens
func extractScreens(url string) (int, error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading contents from %s, Error: %s", url, err)
	}
	//Extract # of theaters
	contentStr := string(content[:])
	rgx := regexp.MustCompile(`(?im)(Widest\s+Release<\/span><span>([\d,]+))`)
	if rgx.MatchString(contentStr) {
		match := rgx.FindAllStringSubmatch(contentStr, -1)
		theathersTxt := strings.ReplaceAll(match[0][2], ",", "")

		theaters, _ := strconv.Atoi(theathersTxt) //ignore error as the regex is digit scoped
		return theaters, nil
	}

	return 0, fmt.Errorf("unable to extract the # of theathers on page %s", url)
}
