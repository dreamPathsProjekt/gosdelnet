package client

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	MaxRetryCount    = 3
	RetryWaitTime    = 1 * time.Second
	MaxRetryWaitTime = 5 * time.Second
	Rows             = 10
)

type Response struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       Data           `json:"response"`
}

type ResponseHeader struct {
	Status int64       `json:"status"`
	QTime  int64       `json:"QTime"`
	Params QueryParams `json:"params"`
}

type QueryParams struct {
	Q      string `json:"q"`
	Indent bool   `json:"indent,omitempty"`
	Wt     string `json:"wt,omitempty"`
	Start  int64  `json:"start,omitempty"`
	Rows   string `json:"rows,omitempty"`
}

type Data struct {
	NumFound      int64  `json:"numFound"`
	Start         int64  `json:"start"`
	NumFoundExact bool   `json:"numFoundExact"`
	Docs          []Book `json:"docs"`
}

type Book struct {
	ProductID           string    `json:"product_id"`
	Translated          string    `json:"translated"`
	PublishingStatus    string    `json:"publishing_status"`
	LastUpdate          time.Time `json:"last_update"`
	ProductStatus       bool      `json:"product_status"`
	Lcx                 bool      `json:"lcx"`
	AddProductTs        time.Time `json:"add_product_ts"`
	ProductFormID       string    `json:"product_form_id"`
	ProductForm         string    `json:"product_form"`
	ProductKind         string    `json:"product_kind"`
	PubID               string    `json:"pub_id"`
	PubName             string    `json:"pub_name"`
	Imprint             string    `json:"imprint"`
	Pages               int64     `json:"pages"`
	DistinctiveTitle    string    `json:"distinctive_title"`
	DistinctiveSubtitle string    `json:"distinctive_subtitle"`
	ThemaCode           []string  `json:"thema_code"`
	Classification      []string  `json:"classification"`
	Category            []string  `json:"category"`
	Audience            []string  `json:"audience"`
	Language            []string  `json:"language"`
	EditionNo           int64     `json:"edition_no"`
	City                string    `json:"city"`
	PubMonth            uint8     `json:"pub_month"`
	PubYear             uint16    `json:"pub_year"`
	PubDay              uint8     `json:"pub_day"`
	LicenseNo           int64     `json:"license_no"`
	Lcno                int64     `json:"lcno"`
	GTin13              string    `json:"gtin13"`
	Prefix              string    `json:"prefix"`
	Isbn13              string    `json:"isbn13"`
	ContrID             []string  `json:"contr_id"`
	ContrName           []string  `json:"contr_name"`
	ContrIDRole         []string  `json:"contr_id_role"`
	ContrIDTab          []string  `json:"contr_id_tab"`
	ContrRole           []string  `json:"contr_role"`
	ContrDenoms         []string  `json:"contr_denoms"`
	LastPriceUpdate     time.Time `json:"last_price_update"`
	Price               float64   `json:"price"`
	PriceChange         float64   `json:"price_change"`
	Vat                 float64   `json:"vat"`
	VatChange           float64   `json:"vat_change"`
	PriceType           string    `json:"price_type"`
	PriceValidUntil     time.Time `json:"price_valid_until"`
	Cover               string    `json:"cover"`
	Version             int64     `json:"_version_"`
}

type CompactBook struct {
	ProductID        string    `json:"product_id"`
	PublishingStatus string    `json:"publishing_status"`
	PubID            string    `json:"pub_id"`
	PubName          string    `json:"pub_name"`
	Isbn13           string    `json:"isbn13"`
	DistinctiveTitle string    `json:"distinctive_title"`
	LastUpdate       time.Time `json:"last_update"`
	LastPriceUpdate  time.Time `json:"last_price_update"`
	PriceValidUntil  time.Time `json:"price_valid_until"`
	Price            float64   `json:"price"`
	PriceChange      float64   `json:"price_change"`
	Vat              float64   `json:"vat"`
	VatChange        float64   `json:"vat_change"`
	PriceType        string    `json:"price_type"`
}

func CSVHeader(w io.Writer, compact bool) {
	cw := csv.NewWriter(w)
	if compact {
		cw.Write([]string{
			"product_id",
			"publishing_status",
			"pub_id",
			"pub_name",
			"isbn13",
			"distinctive_title",
			"last_update",
			"last_price_update",
			"price_valid_until",
			"price",
			"price_change",
			"vat",
			"vat_change",
			"price_type",
		})
	} else {
		cw.Write([]string{
			"product_id",
			"translated",
			"publishing_status",
			"last_update",
			"product_status",
			"lcx",
			"add_product_ts",
			"product_form_id",
			"product_form",
			"product_kind",
			"pub_id",
			"pub_name",
			"imprint",
			"pages",
			"distinctive_title",
			"distinctive_subtitle",
			"thema_code",
			"classification",
			"category",
			"audience",
			"language",
			"edition_no",
			"city",
			"pub_month",
			"pub_year",
			"pub_day",
			"license_no",
			"lcno",
			"gtin13",
			"prefix",
			"isbn13",
			"contr_id",
			"contr_name",
			"contr_id_role",
			"contr_id_tab",
			"contr_role",
			"contr_denoms",
			"last_price_update",
			"price",
			"price_change",
			"vat",
			"vat_change",
			"price_type",
			"price_valid_until",
			"cover",
			"version",
		})
	}
	cw.Flush()
}

func (b *Book) CSVRow(w io.Writer, compact bool) {
	cw := csv.NewWriter(w)
	if compact {
		cw.Write([]string{
			b.ProductID,
			b.PublishingStatus,
			b.PubID,
			b.PubName,
			b.Isbn13,
			b.DistinctiveTitle,
			b.LastUpdate.Format(time.RFC1123),
			b.LastPriceUpdate.Format(time.RFC1123),
			b.PriceValidUntil.Format(time.RFC1123),
			strconv.FormatFloat(b.Price, 'f', 2, 64),
			strconv.FormatFloat(b.PriceChange, 'f', 2, 64),
			strconv.FormatFloat(b.Vat, 'f', 2, 64),
			strconv.FormatFloat(b.VatChange, 'f', 2, 64),
			b.PriceType,
		})
	} else {
		cw.Write([]string{
			b.ProductID,
			b.Translated,
			b.PublishingStatus,
			b.LastUpdate.Format(time.RFC1123),
			strconv.FormatBool(b.ProductStatus),
			strconv.FormatBool(b.Lcx),
			b.AddProductTs.Format(time.RFC1123),
			b.ProductFormID,
			b.ProductForm,
			b.ProductKind,
			b.PubID,
			b.PubName,
			b.Imprint,
			strconv.FormatInt(b.Pages, 10),
			b.DistinctiveTitle,
			b.DistinctiveSubtitle,
			strings.Join(b.ThemaCode, " "),
			strings.Join(b.Classification, " "),
			strings.Join(b.Category, " "),
			strings.Join(b.Audience, " "),
			strings.Join(b.Language, " "),
			strconv.FormatInt(b.EditionNo, 10),
			b.City,
			strconv.FormatInt(int64(b.PubMonth), 10),
			strconv.FormatInt(int64(b.PubYear), 10),
			strconv.FormatInt(int64(b.PubDay), 10),
			strconv.FormatInt(b.LicenseNo, 10),
			strconv.FormatInt(b.Lcno, 10),
			b.GTin13,
			b.Prefix,
			b.Isbn13,
			strings.Join(b.ContrID, " "),
			strings.Join(b.ContrName, " "),
			strings.Join(b.ContrIDRole, " "),
			strings.Join(b.ContrIDTab, " "),
			strings.Join(b.ContrRole, " "),
			strings.Join(b.ContrDenoms, " "),
			b.LastPriceUpdate.Format(time.RFC1123),
			strconv.FormatFloat(b.Price, 'f', 2, 64),
			strconv.FormatFloat(b.PriceChange, 'f', 2, 64),
			strconv.FormatFloat(b.Vat, 'f', 2, 64),
			strconv.FormatFloat(b.VatChange, 'f', 2, 64),
			b.PriceType,
			b.PriceValidUntil.Format(time.RFC1123),
			b.Cover,
			strconv.FormatInt(b.Version, 10),
		})
	}
	cw.Flush()
}
