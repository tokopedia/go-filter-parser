# go-filter-parser
Parses filter query string to a Go struct object.

# Example
```go
package main

import (
    "fmt"
    "log"

    // will be exported as "filter"
    "github.com/tokopedia/go-filter-parser"
)

type SearchShopFilter struct {
    Location          filter.String
    Rating            filter.Range
    GoldMerchantOnly  filter.Bool
}

func (f *SearchShopFilter) FilterMap() filter.FilterMap {
    return filter.FilterMap{
        &f.Location:         "location",
        &f.Rating:           "rating",
        &f.GoldMerchantOnly: "gold_merchant_only",
    }
}

func main() {
    // assume this data is your filter query from query string
    fq := `location=Jakarta;rating=3.5;gold_merchant_only=true`
    separator := ";"
    // create the filter object
    f := new(SearchShopFilter)
    // parse the filter, it will return an error if unable to parse
    if err := filter.Parse(fq, separator, f); err != nil {
        // log.Fatal here just as an example, you can do anything with the error
        log.Fatal(err)
    }

    // example to get the value
    fmt.Println("Operator: %s", filter.OperatorText(f.Location.Operator))
    fmt.Println("Value:    %s", f.Location.Value)
}
```

## Contributors
* [Louis Andris](https://github.com/ruizu)
* [Yuwono Bangun Nagoro](https://github.com/SurgicalSteel)
