package main

import (
	"encoding/xml"
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
		if count%500 == 0 {
			break
		}
	}
	// for _, d := range drugs {
	// 	if len(d.Carriers) > 0 {
	// 		fmt.Println("---------------------------")
	// 		s, _ := json.MarshalIndent(d.Carriers, "", "\t")
	// 		fmt.Println(string(s))
	// 	}
	// }
}

type Drug struct {
	ID                     string               `xml:"drugbank-id"`
	Name                   string               `xml:"name"`
	Description            string               `xml:"description"`
	CAS                    string               `xml:"cas-number"`
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
	ATCCodes               []ATCCode            `xml:"atc-codes"`
	FDALabel               string               `xml:"fda-label"`
	MSDS                   string               `xml:"msds"`
	Patents                []Patent             `xml:"patents"`
	DrugInteractions       []DrugInteraction    `xml:"drug-interactions"`
	Sequences              []Sequence           `xml:"sequences>sequence"`
	ExperimentalProperties []Property           `xml:"experimental-properties>property"`
	ExternalIdentifiers    []ExternalIdentifier `xml:"external-identifiers>external-identifier"`
	ExternalLinks          []ExternalLink       `xml:"external-links>external-link"`
	Pathways               []Pathway            `xml:"pathways>pathway"`
	Targets                []Target             `xml:"targets>target"`
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

type AdverseReaction struct {
	ProteinName     string `xml:"protein-name"`
	GeneSymbol      string `xml:"gene-symbol"`
	UNIPROTID       string `xml:"uniprot-id"`
	Allele          string `xml:"allele"` // TODO check
	Adversereaction string `xml:"adverse-reaction"`
	Description     string `xml:"description"`
	PubmedID        string `xml:"pubmed-id"`
}

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

type Salt struct {
	ID        string `xml:"drugbank-id"`
	Name      string `xml:"name"`
	UNII      string `xml:"unii"`
	CASNumber string `xml:"cas-number"`
	InchiKey  string `xml:"inchikey"`
}

type Brand struct {
	Name    string `xml:"name"`
	Company string `xml:"company"`
}

// <carriers/>
// <transporters/>

//     <transmembrane-regions/>   Polypeptide

type Target struct {
	ID          string      `xml:"id"`
	Organism    string      `xml:"organism"`
	Actions     []string    `xml:"actions>action"`
	References  []Reference `xml:"references"`
	KnownAction string      `xml:"known-action"`
	Polypeptide Polypeptide `xml:"polypeptide"`
}

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

type Pfam struct {
	Identifier string `xml:"identifier"`
	Name       string `xml:"name"`
}

type GoClassifier struct {
	Category    string `xml:"category"`
	Description string `xml:"description"`
}

type Pathway struct {
	SMPDBID  string        `xml:"smpdb-id"`
	Name     string        `xml:"name"`
	Category string        `xml:"category"`
	Drugs    []PathwayDrug `xml:"drugs>drug"`
	Enzymes  []Enzyme      `xml:"enzymes"`
}

type PathwayDrug struct {
	ID   string `xml:"drugbank-id"`
	Name string `xml:"name"`
}

type Enzyme struct {
	UNIPROTID string `xml:"uniprot-id"`
}

type ExternalLink struct {
	Resource string `xml:"resource"`
	URL      string `xml:"url"`
}

type ExternalIdentifier struct {
	Resource   string `xml:"resource"`
	Identifier string `xml:"identifier"`
}

type Property struct {
	Kind   string `xml:"kind"`
	Value  string `xml:"value"`
	Source string `xml:"source"`
}
type Sequence struct {
	Format   string `xml:"format,attr"`
	Sequence string `xml:",chardata"`
}

type DrugInteraction struct {
	ID          string `xml:"drug-interaction>drugbank-id"`
	Name        string `xml:"drug-interaction>name"`
	Description string `xml:"drug-interaction>description"`
}

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
		Code   string `xml:"code,attr"`
		Levels []struct {
			Code        string `xml:"code,attr"`
			Description string `xml:",chardata"`
		} `xml:",any"` // TODO parse into a new struct (Level)
	} `xml:"atc-code"`
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
