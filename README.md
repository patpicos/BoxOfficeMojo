# Overview
This library is used to retrieve the # of screens a movie was showcased in theaters. It focuses on US/UK screens.

It parses the information from BoxOfficeMojo, a sub-site of IMDB using HTML parsing as the API's for this are under license of IMDB Pro.


# Approach
Searching for Box Office information is done using the IMDB movie identifier in the format of `tt######`. From this page, a `releasegroup` URL is parsed out of the HTML.
This 2nd page is parsed to retrive the URL for Domestic and UK links for the actual Box Office data

# Usage

```go
id := "tt1745960"
bom, err := bom.Search(id) // Top Gun Maverick
if err != nil {
    fmt.Printf("Error retrieving data from BoxOfficeMojo for movie id %s, Error: %s", id, err)
}
fmt.Println(bom)
```

# Alternatives
[The-Numbers](https://www.the-numbers.com) provides similar data. It requires a search (which is rather slow) and parsing a result page for the data (not super friendly as the page does not use any classes)