package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/model"
	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/util"
)

type EntityRepository interface {
	GetEntityById(entityID string) (*model.Entity, error)
}

type EntityRepositoryImpl struct {
	SparqlQueryString string
	SparqlServiceUrl  string
}

type entityIdResponse struct {
	Results results `json:"results"`
}

type results struct {
	Bindings []binding `json:"bindings"`
}

type binding struct {
	Type          field `json:"type"`
	Title         field `json:"title"`
	TitleLanguage field `json:"titleLanguage"`
	Organization  field `json:"responsibleOrganization"`
}

func (b *binding) toEntity(entityId string) *model.Entity {
	if b == nil {
		return nil
	}

	return &model.Entity{
		EntityId:     entityId,
		Type:         model.ParseEntityType(&b.Type.Value),
		Title:        b.Title.Value,
		Organization: b.Organization.Value,
	}
}

type field struct {
	Value string `json:"value"`
}

const sparqlFormatQuery = `
PREFIX dcat: <http://www.w3.org/ns/dcat#>
PREFIX modelldcatno: <https://data.norge.no/vocabulary/modelldcatno#>
PREFIX dct: <http://purl.org/dc/terms/>
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
PREFIX skosxl: <http://www.w3.org/2008/05/skos-xl#>
PREFIX skos:  <http://www.w3.org/2004/02/skos/core#>
PREFIX rov:   <http://www.w3.org/ns/regorg#>
PREFIX cpsv: <http://purl.org/vocab/cpsv#>
PREFIX cv: <http://data.europa.eu/m8g/>

SELECT DISTINCT  ?type ?title ?responsibleOrganization ?titleLanguage
WHERE {
    ?record dct:identifier "%s" .
    ?record foaf:primaryTopic ?entity .
    ?entity a ?type .

    # Get title of responsible organization
    OPTIONAL {
        ?entity dct:publisher ?publisher.
        ?publisher dct:identifier ?publisheridentification .
        ?publishernode dct:identifier ?publisheridentification .
        ?publishernode rov:legalName ?publishertitle .
        }

    # If public service
    OPTIONAL {
        ?entity cv:hasCompetentAuthority ?authnode.
        ?authnode rov:legalName ?authtitle.
        }

    bind( IF(?type = cpsv:PublicService, ?authtitle, ?publishertitle) as ?responsibleOrganization )


    # Get title
    OPTIONAL { ?entity dct:title ?dcttitle . }

    # If concept
    OPTIONAL {
        ?entity skosxl:prefLabel ?titleNode .
        ?titleNode skosxl:literalForm ?skostitle .
        }

    bind( IF(?type = skos:Concept, ?skostitle, ?dcttitle) as ?title )
    bind(lang(?title) as ?titleLanguage)
}
`

func (entityRepository *EntityRepositoryImpl) GetEntityById(entityID string) (*model.Entity, error) {
	var parsedRepsonse entityIdResponse
	params := map[string]string{
		"query": fmt.Sprintf(sparqlFormatQuery, entityID),
	}

	rawReponse, err := util.Request(util.RequestOptions{
		Method:          http.MethodGet,
		EndpointUrl:     entityRepository.SparqlServiceUrl,
		QueryParameters: &params,
	})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*rawReponse, &parsedRepsonse)
	if err != nil {
		log.Println("Error on entityIdResponse unmarshal.\n[ERROR] -", err)
		return nil, err
	}

	switch len(parsedRepsonse.Results.Bindings) {
	case 0:
		log.Printf("No entity with id %s found.\n", entityID)
		return nil, errors.New("entity not found")
	case 1:
		return parsedRepsonse.Results.Bindings[0].toEntity(entityID), nil
	default:
		entity, err := getPreferredEntity(parsedRepsonse.Results.Bindings)
		return entity.toEntity(entityID), err
	}
}

func getPreferredEntity(bindings []binding) (*binding, error) {
	var nb, nn, en, any *binding

	if len(bindings) == 0 {
		return nil, errors.New("no entity in results")
	}

	for _, b := range bindings {
		switch b.TitleLanguage.Value {
		case "nb":
			nb = &b
		case "nn":
			nn = &b
		case "en":
			en = &b
		default:
			any = &b
		}
	}

	switch {
	case nb != nil:
		return nb, nil
	case nn != nil:
		return nn, nil
	case en != nil:
		return en, nil
	default:
		return any, nil
	}
}

var CurrentEntityRepository EntityRepository
