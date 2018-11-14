package main

// TODO
//
// <transporters/> -> complex structure containing a polypeptide.
// 					  may it's not worth the hassle given the info
// <transmembrane-regions/> -> it's within a Polypeptide

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar"
)

func main() {
	defer TimeTrack("main", time.Now())
	xmlFile, err := os.Open("drugbank.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	numberOfDrugs := getDrugsNumber(xmlFile)
	bar := progressbar.New(numberOfDrugs)
	drugs := []Drug{}
	xmlFile.Seek(0, 0)
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
		bar.Add(1)
		if count%500 == 0 {
			break
		}
	}
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
	ID                     string               `xml:"drugbank-id"`
	DrugRecordCreatedOn    string               `xml:"created,attr"`
	DrugRecordUpdatedOn    string               `xml:"updated,attr"`
	DrugType               string               `xml:"type,attr"`
	Name                   string               `xml:"name"`
	Description            string               `xml:"description"`
	CAS                    string               `xml:"cas-number"` // Chemical Abstract Service identification number
	UNII                   string               `xml:"unii"`
	State                  string               `xml:"state"`
	Groups                 []Group              `xml:"groups"`
	References             Reference            `xml:"general-references"`
	Indication             string               `xml:"indication"`
	Pharmacodynamics       string               `xml:"pharmacodynamics"`
	MechanismOfAction      string               `xml:"mechanism-of-action"`
	Toxicity               string               `xml:"toxicity"`
	Metabolism             string               `xml:"metabolism"`
	Absorption             string               `xml:"absorption"`
	HalfLife               string               `xml:"half-life"`
	RouteOfElimination     string               `xml:"route-of-elimination"`
	VolumeOfDistribution   string               `xml:"volume-of-distribution"`
	Clearance              string               `xml:"clearance"`
	Classification         Classification       `xml:"classification"`
	Synonyms               []Synonym            `xml:"synonyms"`
	Products               []Product            `xml:"products"`
	Mixtures               []Mixture            `xml:"mixtures"`
	Packagers              []Packager           `xml:"packagers"`
	Manufacturers          []Manufacturer       `xml:"manufacturers"`
	Prices                 []Price              `xml:"prices"`
	Categories             []Category           `xml:"categories"`
	AffectedOrganisms      []Organism           `xml:"affected-organisms"`
	Dosages                []Dosage             `xml:"dosages"`
	ATCCodes               []ATCCode            `xml:"atc-codes"` // WHO drug classification system (ATC) identifiers
	FDALabel               string               `xml:"fda-label"`
	MSDS                   string               `xml:"msds"`
	Patents                []Patent             `xml:"patents"`
	DrugInteractions       []DrugInteraction    `xml:"drug-interactions"`
	Sequences              []Sequence           `xml:"sequences>sequence"`
	ExperimentalProperties []Property           `xml:"experimental-properties>property"`
	ExternalIdentifiers    []ExternalIdentifier `xml:"external-identifiers>external-identifier"`
	ExternalLinks          []ExternalLink       `xml:"external-links>external-link"`
	Targets                []Target             `xml:"targets>target"`
	Pathways               []Pathway            `xml:"pathways>pathway"`
	SynthesysReference     string               `xml:"synthesis-reference"`
	ProteinBinding         string               `xml:"protein-binding"`
	Salts                  []Salt               `xml:"salts>salt"`
	InternationalBrands    []Brand              `xml:"internation-brands>international-brand"`
	AHFSCodes              []string             `xml:"ahfs-code>ahfs-code"`
	PDBEntries             []string             `xml:"pdb-entries>pdb-entry"`
	FoodInteractions       []string             `xml:"food-interactions>food-interaction"`
	Reactions              []Reaction           `xml:"reactions>reaction"`
	SNPEffects             []SNPEffect          `xml:"snp-effects>effect"`
	AdverseReactions       []AdverseReaction    `xml:"snp-adverse-drug-reaction>reaction"`
	Carriers               []Carrier            `xml:"carriers>carrier"`
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

// AdverseReaction represents a possible adverse reaction a drug may cause
type AdverseReaction struct {
	ProteinName     string `xml:"protein-name"`
	GeneSymbol      string `xml:"gene-symbol"`
	UNIPROTID       string `xml:"uniprot-id"`
	Allele          string `xml:"allele"` // TODO check
	Adversereaction string `xml:"adverse-reaction"`
	Description     string `xml:"description"`
	PubmedID        string `xml:"pubmed-id"`
}

// SNPEffect identifies possible nucleotide mutations
// a drug might cause.
// SNP -> Single Nucleotide Polymorphism
type SNPEffect struct {
	ProteinName    string `xml:"protein-name"`
	GeneSymbol     string `xml:"gene-symbol"`
	RSID           string `xml:"rs-id"`
	UNIPROTID      string `xml:"uniprot-id"`
	Allele         string `xml:"allele"` // TODO check
	DefiningChange string `xml:"defining-change"`
	Description    string `xml:"description"`
	PubmedID       string `xml:"pubmed-id"`
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

// Salt represents a salt in which a drug can present itself
type Salt struct {
	ID        string `xml:"drugbank-id"`
	Name      string `xml:"name"`
	UNII      string `xml:"unii"`
	CASNumber string `xml:"cas-number"`
	InchiKey  string `xml:"inchikey"`
}

// Brand identifies brands for mixtures or brand names
type Brand struct {
	Name    string `xml:"name"`
	Company string `xml:"company"`
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

// Pfam represents names and ID numbers of PFAM domains
type Pfam struct {
	Identifier string `xml:"identifier"`
	Name       string `xml:"name"`
}

// GoClassifier represents Gene ontology classification
// including function, cellular process and location
type GoClassifier struct {
	Category    string `xml:"category"`
	Description string `xml:"description"`
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

// Enzyme contains the enzyme ID on UNIPROT
type Enzyme struct {
	UNIPROTID string `xml:"uniprot-id"`
}

// ExternalLink is a link to an external resource
type ExternalLink struct {
	Resource string `xml:"resource"`
	URL      string `xml:"url"`
}

// ExternalIdentifier is an identifier to
// link a drug to external resources
type ExternalIdentifier struct {
	Resource   string `xml:"resource"`
	Identifier string `xml:"identifier"`
}

// Property represents a property of a drug as recorded in the source
type Property struct {
	Kind   string `xml:"kind"`
	Value  string `xml:"value"`
	Source string `xml:"source"`
}

// Sequence represents a sequence of aminoacids
// and the format in which it is represented
type Sequence struct {
	Format   string `xml:"format,attr"`
	Sequence string `xml:",chardata"`
}

// DrugInteraction represents a possible interaction between to drugs
type DrugInteraction struct {
	ID          string `xml:"drug-interaction>drugbank-id"`
	Name        string `xml:"drug-interaction>name"`
	Description string `xml:"drug-interaction>description"`
}

// Patent represents a Patent related to the drug
type Patent struct {
	Number    string `xml:"patent>number"`
	Country   string `xml:"patent>country"`
	Approved  string `xml:"patent>approved"`
	Expires   string `xml:"patent>expires"`
	Pediatric bool   `xml:"patent>pediatric-extension"`
}

// Dosage describes the dosage in which a drug is
// to be administered and the route it should take.
type Dosage struct {
	Form     string `xml:"dosage>form"`
	Route    string `xml:"dosage>route"`
	Strength string `xml:"dosage>strength"` // TODO
}

// ATCCode represents the WHO drug classification system (ATC) identifiers
type ATCCode struct {
	Code struct {
		Code   string `xml:"code,attr"`
		Levels []struct {
			Code        string `xml:"code,attr"`
			Description string `xml:",chardata"`
		} `xml:",any"`
	} `xml:"atc-code"`
}

// Mixture describes a mixture in which a drug can be found
type Mixture struct {
	Name        string `xml:"mixture>name"`
	Ingredients string `xml:"mixture>ingredients"`
}

// Packager describes a packager of the drug
type Packager struct {
	Name string `xml:"packager>name"`
	URL  string `xml:"packager>url"`
}

// Manufacturer describes the manufacturer of a mixture
type Manufacturer struct {
	Name string `xml:"manufacturer"`
	URL  string `xml:"url,attr"`
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

// Category represents a category of sub-division
type Category struct {
	Category string `xml:"category>category"`
	MeshID   string `xml:"category>mesh-id"`
}

// Organism describes an organism affected by a drug
type Organism struct {
	Description string `xml:"affected-organism"`
}

// Synonym describes a synonym of a drug
type Synonym struct {
	Language string `xml:"language,attr"`
	Coder    string `xml:"coder,attr"`
	Synonym  string `xml:"synonym"`
}

// Classification describes the class of a substance
type Classification struct {
	Description string `xml:"description"`
	Parent      string `xml:"direct-parent"`
	Kingdom     string `xml:"kingdom"`
	Superclass  string `xml:"superclass"`
	Class       string `xml:"class"`
	Subclass    string `xml:"subclass"`
}

// Product represents a product in which a drug can be found.
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

// Group describes a category
type Group struct {
	Name string `xml:"group"`
}

// Reference contains information on publications involving a drug
type Reference struct {
	Articles []Article `xml:"articles"`
	Books    []Book    `xml:"textbooks"`
	Links    []Link    `xml:"links"`
}

// Article represents a scientific paper regarding a drug
type Article struct {
	PubMedID string `xml:"article>pubmed-id"`
	Citation string `xml:"article>citation"`
}

// Book represents a textbook regarding a drug
type Book struct {
	ISBN     string `xml:"textbook>isbn"`
	Citation string `xml:"textbook>citation"`
}

// Link is th elink to a resource containing information regarding a drug
type Link struct {
	Title string `xml:"link>title"`
	URL   string `xml:"link>url"`
}
