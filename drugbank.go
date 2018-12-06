// Package main parses the drugbank dataset xml file
package main

// TODO
//
// <transporters/> -> complex structure containing a polypeptide.
// 					  may it's not worth the hassle given the info
// <transmembrane-regions/> -> it's within a Polypeptide

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docopt/docopt-go"

	"github.com/schollz/progressbar"
)

var version = "0.1"

func main() {
	usage := `Drugbank parser.

	Usage:
		drugbank parse <path> <outputdir>
		drugbank process <path> <outputdir> <host> [--password=<password> | --user=<user>]
		drugbank -h | --help
		drugbank --version
	
	Options:
		--password=<password>		Password for Tigergraph instance.
		--user=<user> 			Username for Tigergraph instance.
		-h --help     			Show this screen.
		--version    	 		Show version.`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], version)
	defer TimeTrack("main", time.Now())

	if p, _ := arguments.Bool("parse"); p {
		path, _ := arguments.String("<path>")
		outputdir, _ := arguments.String("<outputdir>")
		fmt.Printf("Parsing %s to %s\n", path, outputdir)
		parse(path, outputdir)
		fmt.Println("Done.")
		os.Exit(0)
	}

	if p, _ := arguments.Bool("process"); p {
		path, _ := arguments.String("<path>")
		outputdir, _ := arguments.String("<outputdir>")
		host, _ := arguments.String("<host>")
		fmt.Printf("Parsing %s to %s...\n", path, outputdir)
		// parse(path, outputdir)
		fmt.Println("Done parsing")
		fmt.Printf("Uploading data to %s...\n", host)
		upload(outputdir)
		os.Exit(0)
	}
}

func upload(directory string) {
	defer TimeTrack("upload", time.Now())
	files, _ := filepath.Glob(filepath.Join(directory, "*.json"))

	fmt.Println("Uploading nodes...")
	for _, file := range files {
		contents, _ := ioutil.ReadFile(file)
		fmt.Println(string(contents))

		break
	}
	fmt.Println("Uploading edges...")
	// for _, file := range files {

	// }
}

func parse(path, outputdir string) {
	defer TimeTrack("parse", time.Now())
	xmlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	numberOfDrugs := getDrugsNumber(xmlFile)
	bar := progressbar.New(numberOfDrugs)
	// drugs := []Drug{}
	xmlFile.Seek(0, 0)
	count := 0

	var (
		jsonDrugs             [][]byte
		jsonManufacturers     [][]byte
		jsonDrugManufacturers [][]byte
		jsonProducts          [][]byte
		jsonDrugProducts      [][]byte
		jsonReactions         [][]byte
		seenReaction          map[string]bool
		jsonAdverseReactions  [][]byte
		jsonSNPEffects        [][]byte
		jsonGroups            [][]byte
		jsonBooks             [][]byte
		jsonArticles          [][]byte
		jsonLinks             [][]byte
		jsonClassifications   [][]byte
		jsonSynonyms          [][]byte
		jsonMixtures          [][]byte
		jsonPackagers         [][]byte
		jsonPrices            [][]byte
		jsonCategories        [][]byte
		jsonOrganisms         [][]byte
		jsonATCCodes          [][]byte
		jsonATCLevels         [][]byte
		jsonDosages           [][]byte
		jsonPatents           [][]byte
		jsonDrugInteractions  [][]byte
		jsonFoodInteractions  [][]byte
		jsonProperties        [][]byte
		jsonExtID             [][]byte
		jsonExtLinks          [][]byte
	)

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
				// drugs = append(drugs, d)

				// DRUG
				jsonDrug, _ := json.Marshal(d)
				jsonDrugs = append(jsonDrugs, jsonDrug)

				// CLASSIFICATION
				jsonClassification, _ := json.Marshal(struct {
					ID string `json:"drugbank-id"`
					Classification
				}{
					d.ID,
					d.Classification,
				})
				jsonClassifications = append(jsonClassifications, jsonClassification)

				// MANUFACTURERS
				for _, manufacturer := range d.Manufacturers {
					if manufacturer.Name == "" {
						continue
					}
					jsonManufacturer, _ := json.Marshal(manufacturer)
					drugManufacturer := struct {
						DrugID         string `json:"drugbank-id"`
						ManufacturerID string `json:"manufacturer-id"`
					}{
						d.ID,
						manufacturer.Name,
					}
					jsonManufacturers = append(jsonManufacturers, jsonManufacturer)

					jsonDrugManufacturer, _ := json.Marshal(drugManufacturer)
					jsonDrugManufacturers = append(jsonDrugManufacturers, jsonDrugManufacturer)
				}

				// PRODUCTS
				for _, product := range d.Products {
					jsonProduct, _ := json.Marshal(product)
					drugProduct := struct {
						DrugID    string `json:"drugbank-id"`
						ProductID string `json:"name"`
					}{
						d.ID,
						product.Name,
					}

					jsonDrugProduct, _ := json.Marshal(drugProduct)

					jsonProducts = append(jsonProducts, jsonProduct)
					jsonDrugProducts = append(jsonDrugProducts, jsonDrugProduct)
				}

				// REACTIONS
				for _, reaction := range d.Reactions {
					_, ok := seenReaction[reaction.Sequence]
					if ok {
						continue
					}
					jsonReaction, _ := json.Marshal(struct {
						Sequence  string `json:"sequence"`
						LeftID    string `json:"left-id"`
						LeftName  string `json:"left-name"`
						RightID   string `json:"right-id"`
						RightName string `json:"right-name"`
					}{
						reaction.Sequence,
						reaction.Left.ID,
						reaction.Left.Name,
						reaction.Right.ID,
						reaction.Right.Name,
					})
					jsonReactions = append(jsonReactions, jsonReaction)
				}

				// ADVERSE REACTIONS
				for _, reaction := range d.AdverseReactions {
					if reaction.UNIPROTID == "" {
						continue
					}
					jsonAdverseReaction, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						AdverseReaction
					}{
						d.ID,
						reaction,
					})
					jsonAdverseReactions = append(jsonAdverseReactions, jsonAdverseReaction)
				}

				// SNP EFFECTS
				for _, effect := range d.SNPEffects {
					if effect.UNIPROTID == "" {
						continue
					}
					jsonEffect, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						SNPEffect
					}{
						d.ID,
						effect,
					})
					jsonSNPEffects = append(jsonSNPEffects, jsonEffect)
				}

				// GROUPS
				for _, group := range d.Groups {
					if group.Name == "" {
						continue
					}
					jsonGroup, _ := json.Marshal(struct {
						ID   string `json:"drugbank-id"`
						Name string `json:"name"`
					}{
						d.ID,
						group.Name,
					})
					jsonGroups = append(jsonGroups, jsonGroup)
				}

				// REFERENCES
				// BOOKS
				for _, book := range d.References.Books {
					if book.ISBN == "" {
						continue
					}
					jsonBook, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Book
					}{
						d.ID,
						book,
					})
					jsonBooks = append(jsonBooks, jsonBook)
				}

				// LINKS
				for _, link := range d.References.Links {
					if link.URL == "" {
						continue
					}
					jsonLink, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Link
					}{
						d.ID,
						link,
					})
					jsonLinks = append(jsonLinks, jsonLink)
				}

				// PAPERS
				for _, paper := range d.References.Articles {
					if paper.PubMedID == "" {
						continue
					}
					jsonArticle, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Article
					}{
						d.ID,
						paper,
					})
					jsonArticles = append(jsonArticles, jsonArticle)
				}

				// SYNONYMS
				for _, syn := range d.Synonyms {
					if syn.Synonym == "" {
						continue
					}
					jsonSynonym, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Synonym
					}{
						d.ID,
						syn,
					})
					jsonSynonyms = append(jsonSynonyms, jsonSynonym)
				}

				// MIXTURES
				for _, mix := range d.Mixtures {
					if mix.Name == "" {
						continue
					}
					jsonMixture, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Mixture
					}{
						d.ID,
						mix,
					})
					jsonMixtures = append(jsonMixtures, jsonMixture)
				}

				// PACKAGERS
				for _, pack := range d.Packagers {
					if pack.Name == "" {
						continue
					}
					jsonPackager, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Packager
					}{
						d.ID,
						pack,
					})
					jsonPackagers = append(jsonPackagers, jsonPackager)
				}

				// PRICES
				for _, price := range d.Prices {
					if price.Details.Amount == 0.0 {
						continue
					}
					jsonPrice, _ := json.Marshal(struct {
						DrugID      string  `json:"drugbank-id"`
						Description string  `json:"description"`
						Amount      float64 `json:"cost"`
						Currency    string  `json:"currency"`
						Unit        string  `json:"sale-unit"`
					}{
						d.ID,
						price.Description,
						price.Details.Amount,
						price.Details.Currency,
						price.Unit,
					})
					jsonPrices = append(jsonPrices, jsonPrice)
				}

				// CATEGORY
				for _, cat := range d.Categories {
					if cat.Category == "" {
						continue
					}
					jsonCategory, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Category
					}{
						d.ID,
						cat,
					})
					jsonCategories = append(jsonCategories, jsonCategory)
				}

				// AFFECTED ORGANISMS
				for _, org := range d.AffectedOrganisms {
					if org.Description == "" {
						continue
					}
					jsonOrganism, _ := json.Marshal(struct {
						DrugID   string `json:"drugbank-id"`
						Organism string `json:"organism"`
					}{
						d.ID,
						org.Description,
					})
					jsonOrganisms = append(jsonOrganisms, jsonOrganism)
				}

				// ATC CODES
				for _, code := range d.ATCCodes {
					jsonCode, _ := json.Marshal(struct {
						ATCCode string `json:"atc-code"`
						DrugID  string `json:"drugbank-id"`
					}{
						code.Code.Code,
						d.ID,
					})

					jsonATCCodes = append(jsonATCCodes, jsonCode)

					for _, level := range code.Code.Levels {

						jsonLevel, _ := json.Marshal(struct {
							ATCCode      string `json:"atc-code"`
							ATCLevelCode string `json:"atc-level"`
							Description  string `json:"description"`
						}{
							code.Code.Code,
							level.Code,
							level.Description,
						})

						jsonATCLevels = append(jsonATCLevels, jsonLevel)
					}
				}

				// DOSAGE
				for _, dosage := range d.Dosages {
					if dosage.Form == "" {
						continue
					}
					jsonDosage, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Dosage
					}{
						d.ID,
						dosage,
					})
					jsonDosages = append(jsonDosages, jsonDosage)
				}

				// PATENT
				for _, patent := range d.Patents {
					if patent.Number == "" {
						continue
					}
					jsonPatent, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Patent
					}{
						d.ID,
						patent,
					})
					jsonPatents = append(jsonPatents, jsonPatent)
				}

				// DRUG INTERACTION
				for _, interaction := range d.DrugInteractions {
					if interaction.ID == "" {
						continue
					}
					jsonInteraction, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						DrugInteraction
					}{
						d.ID,
						interaction,
					})
					jsonDrugInteractions = append(jsonDrugInteractions, jsonInteraction)
				}

				// FOOD INTERACTION
				for _, interaction := range d.FoodInteractions {
					if interaction == "" {
						continue
					}
					jsonInteraction, _ := json.Marshal(struct {
						DrugID      string `json:"drugbank-id"`
						Interaction string `json:"interaction"`
					}{
						d.ID,
						interaction,
					})
					jsonFoodInteractions = append(jsonFoodInteractions, jsonInteraction)
				}

				// PROPERTIES
				for _, property := range d.ExperimentalProperties {
					if property.Value == "" {
						continue
					}
					jsonProperty, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						Property
					}{
						d.ID,
						property,
					})
					jsonProperties = append(jsonProperties, jsonProperty)
				}

				// EXTERNAL LINK
				for _, link := range d.ExternalLinks {
					if link.URL == "" {
						continue
					}
					jsonLink, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						ExternalLink
					}{
						d.ID,
						link,
					})
					jsonExtLinks = append(jsonExtLinks, jsonLink)
				}

				// EXTERNAL IDENTIFIERS
				for _, id := range d.ExternalIdentifiers {
					if id.Identifier == "" {
						continue
					}
					jsonID, _ := json.Marshal(struct {
						DrugID string `json:"drugbank-id"`
						ExternalIdentifier
					}{
						d.ID,
						id,
					})
					jsonExtID = append(jsonExtID, jsonID)
				}
			}
		}
		count++
		bar.Add(1)
	}
	fmt.Println()

	err = os.MkdirAll(outputdir, 0770)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(filepath.Join(outputdir, "drugs.json"), bytes.Join(jsonDrugs, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "classifications.json"), bytes.Join(jsonClassifications, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "manufacturers.json"), bytes.Join(jsonManufacturers, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "drugs-manufacturers-join.json"), bytes.Join(jsonDrugManufacturers, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "products.json"), bytes.Join(jsonProducts, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "drugs-products-join.json"), bytes.Join(jsonDrugProducts, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "reactions.json"), bytes.Join(jsonReactions, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "adverse-reactions.json"), bytes.Join(jsonReactions, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "snp-effects.json"), bytes.Join(jsonSNPEffects, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "groups.json"), bytes.Join(jsonGroups, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "articles.json"), bytes.Join(jsonArticles, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "books.json"), bytes.Join(jsonBooks, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "links.json"), bytes.Join(jsonLinks, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "synonyms.json"), bytes.Join(jsonSynonyms, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "mixtures.json"), bytes.Join(jsonMixtures, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "packagers.json"), bytes.Join(jsonPackagers, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "prices.json"), bytes.Join(jsonPrices, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "categories.json"), bytes.Join(jsonCategories, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "organisms.json"), bytes.Join(jsonOrganisms, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "atc_codes.json"), bytes.Join(jsonATCCodes, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "atc_levels.json"), bytes.Join(jsonATCLevels, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "dosages.json"), bytes.Join(jsonDosages, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "patents.json"), bytes.Join(jsonPatents, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "drug_interactions.json"), bytes.Join(jsonDrugInteractions, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "food_interactions.json"), bytes.Join(jsonFoodInteractions, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "experimental_properties.json"), bytes.Join(jsonProperties, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "external_links.json"), bytes.Join(jsonExtLinks, []byte("\n")), 0644)
	ioutil.WriteFile(filepath.Join(outputdir, "external_identifiers.json"), bytes.Join(jsonExtID, []byte("\n")), 0644)
}

// getDrugsNumber counts the number of opening drug tags
// in the xml file. Th tags that contain attributes are top-level
// tags for drugs
func getDrugsNumber(file io.Reader) int {
	defer TimeTrack("getDrugsNumber", time.Now())
	counter := 0
	pattern := "<drug type="
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			counter++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return counter
}

// TimeTrack tracks the execution time of a function
func TimeTrack(name string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %.2f seconds\n", name, elapsed.Seconds())
}

// Drug represents a drug and all its related information
// More detailed info available at https://www.drugbank.ca/documentation#drug-cards
type Drug struct {
	ID                     string               `xml:"drugbank-id" json:"drugbank-id"`
	DrugRecordCreatedOn    string               `xml:"created,attr" json:"record-creation"`
	DrugRecordUpdatedOn    string               `xml:"updated,attr" json:"record-update"`
	DrugType               string               `xml:"type,attr" json:"drug-type"`
	Name                   string               `xml:"name" json:"name"`
	Description            string               `xml:"description" json:"description"`
	CAS                    string               `xml:"cas-number" json:"cas-number"` // Chemical Abstract Service identification number
	UNII                   string               `xml:"unii" json:"unii"`
	State                  string               `xml:"state" json:"state"`
	Groups                 []Group              `xml:"groups" json:"-"`
	References             Reference            `xml:"general-references" json:"-"`
	Indication             string               `xml:"indication" json:"indication"`
	Pharmacodynamics       string               `xml:"pharmacodynamics" json:"pharmacodynamycs"`
	MechanismOfAction      string               `xml:"mechanism-of-action" json:"mechanism-of-action"`
	Toxicity               string               `xml:"toxicity" json:"toxicity"`
	Metabolism             string               `xml:"metabolism" json:"metabolism"`
	Absorption             string               `xml:"absorption" json:"absorption"`
	HalfLife               string               `xml:"half-life" json:"half-life"`
	RouteOfElimination     string               `xml:"route-of-elimination" json:"route-of-elimination"`
	VolumeOfDistribution   string               `xml:"volume-of-distribution" json:"volume-of-distribution"`
	Clearance              string               `xml:"clearance" json:"clearance"`
	Classification         Classification       `xml:"classification" json:"-"`
	Synonyms               []Synonym            `xml:"synonyms" json:"-"`
	Products               []Product            `xml:"products" json:"-"`
	Mixtures               []Mixture            `xml:"mixtures" json:"-"`
	Packagers              []Packager           `xml:"packagers" json:"-"`
	Manufacturers          []Manufacturer       `xml:"manufacturers" json:"-"`
	Prices                 []Price              `xml:"prices" json:"-"`
	Categories             []Category           `xml:"categories" json:"-"`
	AffectedOrganisms      []Organism           `xml:"affected-organisms" json:"-"`
	Dosages                []Dosage             `xml:"dosages" json:"-"`
	ATCCodes               []ATCCode            `xml:"atc-codes" json:"-"` // WHO drug classification system (ATC) identifiers
	FDALabel               string               `xml:"fda-label" json:"fda-label"`
	MSDS                   string               `xml:"msds" json:"msds"`
	Patents                []Patent             `xml:"patents" json:"-"`
	DrugInteractions       []DrugInteraction    `xml:"drug-interactions" json:"-"`
	Sequences              []Sequence           `xml:"sequences>sequence" json:"-"`
	ExperimentalProperties []Property           `xml:"experimental-properties>property" json:"-"`
	ExternalIdentifiers    []ExternalIdentifier `xml:"external-identifiers>external-identifier" json:"-"`
	ExternalLinks          []ExternalLink       `xml:"external-links>external-link" json:"-"`
	Targets                []Target             `xml:"targets>target" json:"-"`
	Pathways               []Pathway            `xml:"pathways>pathway" json:"-"`
	SynthesysReference     string               `xml:"synthesis-reference" json:"synthesis-reference"`
	ProteinBinding         string               `xml:"protein-binding" json:"protein-binding"`
	Salts                  []Salt               `xml:"salts>salt" json:"-"`
	InternationalBrands    []Brand              `xml:"internation-brands>international-brand" json:"-"`
	AHFSCodes              []string             `xml:"ahfs-code>ahfs-code" json:"-"`
	PDBEntries             []string             `xml:"pdb-entries>pdb-entry" json:"-"`
	FoodInteractions       []string             `xml:"food-interactions>food-interaction" json:"-"`
	Reactions              []Reaction           `xml:"reactions>reaction" json:"-"`
	SNPEffects             []SNPEffect          `xml:"snp-effects>effect" json:"-"`
	AdverseReactions       []AdverseReaction    `xml:"snp-adverse-drug-reaction>reaction" json:"-"`
	Carriers               []Carrier            `xml:"carriers>carrier" json:"-"`
}

// AdverseReaction represents a possible adverse reaction a drug may cause
type AdverseReaction struct {
	ProteinName     string `xml:"protein-name" json:"protein-name"`
	GeneSymbol      string `xml:"gene-symbol" json:"gene-symbol"`
	UNIPROTID       string `xml:"uniprot-id" json:"uniprot-id"`
	Allele          string `xml:"allele" json:"allele"` // TODO check
	AdverseReaction string `xml:"adverse-reaction" json:"adverse-reaction"`
	Description     string `xml:"description" json:"description"`
	PubmedID        string `xml:"pubmed-id" json:"pubmed-id"`
}

// Article represents a scientific paper regarding a drug
type Article struct {
	PubMedID string `xml:"article>pubmed-id" json:"pubmed-id"`
	Citation string `xml:"article>citation" json:"citation"`
}

// ATCCode represents the WHO drug classification system (ATC) identifiers
type ATCCode struct {
	Code struct {
		Code   string `xml:"code,attr"  json:"code"`
		Levels []struct {
			Code        string `xml:"code,attr"`
			Description string `xml:",chardata"`
		} `xml:",any"`
	} `xml:"atc-code"`
}

// Book represents a textbook regarding a drug
type Book struct {
	ISBN     string `xml:"textbook>isbn" json:"isbn"`
	Citation string `xml:"textbook>citation" json:"citation"`
}

// Brand identifies brands for mixtures or brand names
type Brand struct {
	Name    string `xml:"name" json:"name"`
	Company string `xml:"company" json:"company"`
}

// Carrier represents a secreted protein which binds to drugs,
// carrying them to cell transporters, where they are moved into the cell.
// Drug carriers may be used in drug design to increase the
// effectiveness of drug delivery to the target sites of pharmacological actions.
// Targets, Enzymes, Carriers, And Transporters may switch roles depending
// on the drug to which they bind. Some drugs specifically target transporters,
// and in this case a transporter can also be the target
// (for example: Procaine targeting the Sodium-dependent dopamine transporter).
type Carrier struct {
	Position    string      `xml:"position,attr"`
	ID          string      `xml:"id"`
	Name        string      `xml:"name"`
	Organism    string      `xml:"organism"`
	Actions     []string    `xml:"actions>action"`
	References  []Reference `xml:"references"`
	KnownAction string      `xml:"known-action"`
	Polypeptide Polypeptide `xml:"polypeptide"`
}

// Category represents a category of sub-division
type Category struct {
	Category string `xml:"category>category" json:"category"`
	MeshID   string `xml:"category>mesh-id" json:"mesh-id"`
}

// Classification describes the class of a substance
type Classification struct {
	Description string `xml:"description" json:"description"`
	Parent      string `xml:"direct-parent" json:"direct-parent"`
	Kingdom     string `xml:"kingdom" json:"kingdom"`
	Superclass  string `xml:"superclass" json:"superclass"`
	Class       string `xml:"class" json:"class"`
	Subclass    string `xml:"subclass" json:"subclass"`
}

// Dosage describes the dosage in which a drug is
// to be administered and the route it should take.
type Dosage struct {
	Form     string `xml:"dosage>form" json:"form"`
	Route    string `xml:"dosage>route" json:"route"`
	Strength string `xml:"dosage>strength" json:"strength"` // TODO
}

// DrugInteraction represents a possible interaction between to drugs
type DrugInteraction struct {
	ID          string `xml:"drug-interaction>drugbank-id" json:"reagent-id"`
	Name        string `xml:"drug-interaction>name" json:"name"`
	Description string `xml:"drug-interaction>description" json:"description"`
}

// Enzyme contains the enzyme ID on UNIPROT
type Enzyme struct {
	UNIPROTID string `xml:"uniprot-id" json:"uniprot0id"`
}

// ExternalIdentifier is an identifier to
// link a drug to external resources
type ExternalIdentifier struct {
	Resource   string `xml:"resource" json:"resource"`
	Identifier string `xml:"identifier" json:"identifier"`
}

// ExternalLink is a link to an external resource
type ExternalLink struct {
	Resource string `xml:"resource" json:"resource"`
	URL      string `xml:"url" json:"url"`
}

// GoClassifier represents Gene ontology classification
// including function, cellular process and location
type GoClassifier struct {
	Category    string `xml:"category"`
	Description string `xml:"description"`
}

// Group describes a category
type Group struct {
	Name string `xml:"group" json:"name"`
}

// Link is th elink to a resource containing information regarding a drug
type Link struct {
	Title string `xml:"link>title" json:"title"`
	URL   string `xml:"link>url" json:"url"`
}

// Manufacturer describes the manufacturer of a mixture
type Manufacturer struct {
	Name string `xml:"manufacturer" json:"name"`
	URL  string `xml:"url,attr" json:"url"`
}

// Mixture describes a mixture in which a drug can be found
type Mixture struct {
	Name        string `xml:"mixture>name" json:"name"`
	Ingredients string `xml:"mixture>ingredients" json:"ingredients"`
}

// Organism describes an organism affected by a drug
type Organism struct {
	Description string `xml:"affected-organism"`
}

// Packager describes a packager of the drug
type Packager struct {
	Name string `xml:"packager>name" json:"name"`
	URL  string `xml:"packager>url" json:"url"`
}

// Patent represents a Patent related to the drug
type Patent struct {
	Number    string `xml:"patent>number" json:"number"`
	Country   string `xml:"patent>country" json:"country"`
	Approved  string `xml:"patent>approved" json:"approved"`
	Expires   string `xml:"patent>expires" json:"expiration"`
	Pediatric bool   `xml:"patent>pediatric-extension" json:"pediatric"`
}

// Pathway represents  processes (from SMPD) that the given molecule is involved in
// A protein, macromolecule, nucleic acid, or small molecule
// to which a given drug binds, resulting in an alteration of the
// normal function of the bound molecule and a desirable therapeutic effect.
// Drug targets are most commonly proteins such as enzymes,
// ion channels, and receptors.
type Pathway struct {
	SMPDBID  string        `xml:"smpdb-id"`
	Name     string        `xml:"name"`
	Category string        `xml:"category"`
	Drugs    []PathwayDrug `xml:"drugs>drug"`
	Enzymes  []Enzyme      `xml:"enzymes"`
}

// PathwayDrug identifies drugs involved with pathways
type PathwayDrug struct {
	ID   string `xml:"drugbank-id"`
	Name string `xml:"name"`
}

// Pfam represents names and ID numbers of PFAM domains
type Pfam struct {
	Identifier string `xml:"identifier"`
	Name       string `xml:"name"`
}

// Polypeptide represents a single polypeptide and its relative details
type Polypeptide struct {
	ID                 string `xml:"id,attr"`
	Source             string `xml:"source,attr"`
	Name               string `xml:"name"`
	GeneralFunction    string `xml:"general-function"`
	SpecificFunction   string `xml:"specific-function"`
	GeneName           string `xml:"gene-name"`
	Locus              string `xml:"locus"`
	CellularLocation   string `xml:"cellular-location"`
	SignalRegion       string `xml:"signal-regions"`
	TheoreticalPi      string `xml:"theoretical-pi"`
	MolecularWeight    string `xml:"molecular-weight"`
	ChromosomeLocation string `xml:"chromosome-location"`
	OrganismTaxonomy   struct {
		TaxonomyID string `xml:"ncbi-taxonomy-id,attr"`
		Organism   string `xml:",chardata"`
	} `xml:"organism"`
	ExternalIdentifiers []ExternalIdentifier `xml:"external-identifiers>external-identifier"`
	Synonyms            []string             `xml:"synonyms>synonym"`
	AminoAcidSequence   struct {
		Format   string `xml:"format,attr"`
		Sequence string `xml:",chardata"`
	} `xml:"amino-acid-sequence"`
	GeneSequence struct {
		Format   string `xml:"format,attr"`
		Sequence string `xml:",chardata"`
	} `xml:"gene-sequence"`
	Pfams         []Pfam         `xml:"pfams>pfam"`
	GoClassifiers []GoClassifier `xml:"go-classifiers>go-classifier"`
}

// Price details the cost and currency of a medication
type Price struct {
	Description string `xml:"price>description"`
	Details     struct {
		Amount   float64 `xml:",chardata"`
		Currency string  `xml:"currency,attr"` // TODO
	} `xml:"price>cost"`
	Unit string `xml:"price>unit"`
}

// Product represents a product in which a drug can be found.
type Product struct {
	Name                 string `xml:"product>name" json:"name"`
	Labeller             string `xml:"product>labeller" json:"labeller"`
	NDCID                string `xml:"product>ndc-id" json:"ncd-id"`
	NDCProductCode       string `xml:"product>ndc-product-code" json:"ncd-product-code"`
	DPDID                string `xml:"product>dpd-id" json:"dpd-id"`
	EMAProductCode       string `xml:"product>ema-product-code" json:"ema-product-code"`
	EMAProductNumber     string `xml:"product>ema-ma-number" json:"ema-product-number"`
	StartedMarketing     string `xml:"product>started-marketing-on" json:"started-marketing-on"`
	EndedMarketing       string `xml:"product>ended-marketing-on" json:"ended-marketing-on"`
	DosageForm           string `xml:"product>dosage-form" json:"dosage-form"`
	Strength             string `xml:"product>strength" json:"strngth"`
	Route                string `xml:"product>route" json:"route"`
	FDAApplicationNumber string `xml:"product>fda-application-number" json:"fda-application-number"`
	Generic              bool   `xml:"product>generic" json:"generic"`
	OverTheCounter       bool   `xml:"product>over-the-counter" json:"over-the-counter"`
	Approved             bool   `xml:"product>approved" json:"approved"`
	Country              string `xml:"product>country" json:"country"`
	Source               string `xml:"product>source" json:"source"`
}

// Property represents a property of a drug as recorded in the source
type Property struct {
	Kind   string `xml:"kind" json:"kind"`
	Value  string `xml:"value" json:"value"`
	Source string `xml:"source" json:"source"`
}

// Reaction describes a reaction a specific drug can undergo with
// another reagent
type Reaction struct {
	Sequence string `xml:"sequence"`
	Left     struct {
		ID   string `xml:"drugbank-id"`
		Name string `xml:"name"`
	} `xml:"left-element"`
	Right struct {
		ID   string `xml:"drugbank-id"`
		Name string `xml:"name"`
	} `xml:"right-element"`
	Enzymes []string `xml:"uniprot-id"`
}

// Reference contains information on publications involving a drug
type Reference struct {
	Articles []Article `xml:"articles"`
	Books    []Book    `xml:"textbooks"`
	Links    []Link    `xml:"links"`
}

// Salt represents a salt in which a drug can present itself
type Salt struct {
	ID        string `xml:"drugbank-id"`
	Name      string `xml:"name"`
	UNII      string `xml:"unii"`
	CASNumber string `xml:"cas-number"`
	InchiKey  string `xml:"inchikey"`
}

// Sequence represents a sequence of aminoacids
// and the format in which it is represented
type Sequence struct {
	Format   string `xml:"format,attr"`
	Sequence string `xml:",chardata"`
}

// SNPEffect identifies possible nucleotide mutations
// a drug might cause.
// SNP -> Single Nucleotide Polymorphism
type SNPEffect struct {
	ProteinName    string `xml:"protein-name" json:"protein-name"`
	GeneSymbol     string `xml:"gene-symbol" json:"gene-symbol"`
	RSID           string `xml:"rs-id" json:"rs-id"`
	UNIPROTID      string `xml:"uniprot-id" json:"uniprot-id"`
	Allele         string `xml:"allele" json:"allele"` // TODO check
	DefiningChange string `xml:"defining-change" json:"defining-change"`
	Description    string `xml:"description" json:"description"`
	PubmedID       string `xml:"pubmed-id" json:"pubmed-id"`
}

// Synonym describes a synonym of a drug
type Synonym struct {
	Language string `xml:"language,attr" json:"language"`
	Coder    string `xml:"coder,attr" json:"coder"`
	Synonym  string `xml:"synonym" json:"synonym"`
}

// Target represents a protein, macromolecule, nucleic acid,
// or small molecule to which a given drug binds,
// resulting in an alteration of the normal function of the
// bound molecule and a desirable therapeutic effect.
// Drug targets are most commonly proteins such as enzymes,
// ion channels, and receptors.
type Target struct {
	ID          string      `xml:"id"`
	Organism    string      `xml:"organism"`
	Actions     []string    `xml:"actions>action"`
	References  []Reference `xml:"references"`
	KnownAction string      `xml:"known-action"`
	Polypeptide Polypeptide `xml:"polypeptide"`
}
