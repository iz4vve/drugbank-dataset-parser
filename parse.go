package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func main() {
	xmlFile, err := os.Open("drugbank.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)

	drugs := []Drug{}

	count := 0

	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "drug" {
				var d Drug
				decoder.DecodeElement(&d, &startElement)
				drugs = append(drugs, d)
			}
		}
		count++
		if count%20 == 0 {
			break
		}
	}
	for _, d := range drugs {
		fmt.Println("---------------------------")
		s, _ := json.MarshalIndent(d, "", "\t")
		fmt.Println(string(s))
	}
}

type Drug struct {
	ID                   string         `xml:"drugbank-id"`
	Name                 string         `xml:"name"`
	Description          string         `xml:"description"`
	CAS                  string         `xml:"cas-number"`
	UNII                 string         `xml:"unii"`
	State                string         `xml:"state"`
	Groups               []Group        `xml:"groups"`
	References           Reference      `xml:"general-references"`
	Indication           string         `xml:"indication"`
	Pharmacodynamics     string         `xml:"pharmacodynamics"`
	MechanismOfAction    string         `xml:"mechanism-of-action"`
	Toxicity             string         `xml:"toxicity"`
	Metabolism           string         `xml:"metabolism"`
	Absorption           string         `xml:"absorption"`
	HalfLife             string         `xml:"half-life"`
	RouteofElimination   string         `xml:"route-of-elimination"`
	VolumeOfDistribution string         `xml:"volume-of-distribution"`
	Clearance            string         `xml:"clearance"`
	Classification       Classification `xml:"classification"`
	Synonyms             []Synonym      `xml:"synonyms"`
	Products             []Product      `xml:"products"`
	Mixtures             []Mixture      `xml:"mixtures"`
	Packagers            []Packager     `xml:"packagers"`
	Manufacturers        []Manufacturer `xml:"manufacturers"`
	Prices               []Price        `xml:"prices"`
	Categories           []Category     `xml:"categories"`
	AffectedOrganisms    []Organism     `xml:"affected-organisms"`
	Dosages              []Dosage       `xml:"dosages"`
	ATCCodes             []ATCCode      `xml:"atc-codes"`
	FDALabel             string         `xml:"fda-label"`
	MSDS                 string         `xml:"msds"`
	Patents              []Patent       `xml:"patents"`
}

// <synthesis-reference/>
//    <protein-binding/>
//   <salts/>
// <international-brands/>
// <ahfs-codes/>
// <pdb-entries/>

type Patent struct {
	Number    string `xml:"patent>number"`
	Country   string `xml:"patent>country"`
	Approved  string `xml:"patent>approved"`
	Expires   string `xml:"patent>expires"`
	Pediatric bool   `xml:"patent>pediatric-extension"`
}

type Dosage struct {
	Form     string `xml:"dosage>form"`
	Route    string `xml:"dosage>route"`
	Strength string `xml:"dosage>strength"` // TODO
}

type ATCCode struct {
	Code struct {
		code   string `xml:",attr"`
		Levels string `xml:",innerxml"` // TODO parse into a new struct (Level)
	} `xml:"atc-code"`
}

type Level struct {
	Description struct {
		Code        string `xml:"code,attr"`
		Description string `xml:",chardata"`
	} `xml:"level"`
}

type Mixture struct {
	Name        string `xml:"mixture>name"`
	Ingredients string `xml:"mixture>ingredients"`
}

type Packager struct {
	Name string `xml:"packager>name"`
	URL  string `xml:"packager>url"`
}

type Manufacturer struct {
	Name string `xml:"manufacturer"`
	URL  string `xml:"url,attr"`
}

type Price struct {
	Description string `xml:"price>description"`
	Details     struct {
		Amount   float64 `xml:",chardata"`
		Currency string  `xml:"currency,attr"` // TODO
	} `xml:"price>cost"`
	Unit string `xml:"price>unit"`
}

type Category struct {
	Category string `xml:"category>category"`
	MeshID   string `xml:"category>mesh-id"`
}

type Organism struct {
	Description string `xml:"affected-organism"`
}

type Synonym struct {
	Language string `xml:"language,attr"`
	Coder    string `xml:"coder,attr"`
	Synonym  string `xml:"synonym"`
}

type Classification struct {
	Description string `xml:"description"`
	Parent      string `xml:"direct-parent"`
	Kingdom     string `xml:"kingdom"`
	Superclass  string `xml:"superclass"`
	Class       string `xml:"class"`
	Subclass    string `xml:"subclass"`
}

type Product struct {
	Name                 string `xml:"product>name"`
	Labeller             string `xml:"product>labeller"`
	NDCID                string `xml:"product>ndc-id"`
	NDCProductCode       string `xml:"product>ndc-product-code"`
	DPDID                string `xml:"product>dpd-id"`
	EMAProductCode       string `xml:"product>ema-product-code"`
	EMAProductNumber     string `xml:"product>ema-ma-number"`
	StartedMarketing     string `xml:"product>started-marketing-on"`
	EndedMarketing       string `xml:"product>ended-marketing-on"`
	DosageForm           string `xml:"product>dosage-form"`
	Strength             string `xml:"product>strength"`
	Route                string `xml:"product>route"`
	FDAApplicationNumber string `xml:"product>fda-application-number"`
	Generic              bool   `xml:"product>generic"`
	OverTheCounter       bool   `xml:"product>over-the-counter"`
	Approved             bool   `xml:"product>approved"`
	Country              string `xml:"product>country"`
	Source               string `xml:"product>source"`
}

type Group struct {
	Name string `xml:"group"`
}

type Reference struct {
	Articles []Article `xml:"articles"`
	Books    []Book    `xml:"textbooks"`
	Links    []Link    `xml:"links"`
}

type Article struct {
	PubMedID string `xml:"article>pubmed-id"`
	Citation string `xml:"article>citation"`
}

type Book struct {
	ISBN     string `xml:"textbook>isbn"`
	Citation string `xml:"textbook>citation"`
}

type Link struct {
	Title string `xml:"link>title"`
	URL   string `xml:"link>url"`
}
