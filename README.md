# go-filter-parser
Filter parser for Go.

# Example
```go
package main

import (
    "fmt"

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
    fq := `any_string==hello world\;foo\;bar;any_bool!=True;any_float==12345.6789;any_range!=5000..10000`

    // create the filter object
    f := new(SearchShopFilter)
    // parse the filter, it will return an error if unable to parse
    if err := filter.Parse(fq, f); err != nil {
        // log.Fatal here just as an example, you can do anything with the error
        log.Fatal(err)
    }

    // example to get the value
    fmt.Println("Operator: %s", filter.OperatorText(f.Location.Operator))
    fmt.Println("Value:    %s", f.Location.Value)
}
```
