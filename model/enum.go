package model

import (
	"log"
	"regexp"
)

type EntityType int

const (
	Undefined EntityType = iota
	Dataset
	DataService
	Concept
	InformationModel
	PublicService
	Event
)

func ParseEntityType(str *string) EntityType {
	if str == nil {
		return Undefined
	}

	dataset_match, err := regexp.Match("(?i)Dataset", []byte(*str))
	if err != nil {
		log.Println("error parsing dataset EntityType")
		return Undefined
	}

	dataService_match, err := regexp.Match("(?i)DataService", []byte(*str))
	if err != nil {
		log.Println("error parsing DataService EntityType")
		return Undefined
	}

	concept_match, err := regexp.Match("(?i)Concept", []byte(*str))
	if err != nil {
		log.Println("error parsing Concept EntityType")
		return Undefined
	}

	informationModel_match, err := regexp.Match("(?i)InformationModel", []byte(*str))
	if err != nil {
		log.Println("error parsing InformationModel EntityType")
		return Undefined
	}

	publicService_match, err := regexp.Match("(?i)PublicService", []byte(*str))
	if err != nil {
		log.Println("error parsing PublicService EntityType")
		return Undefined
	}

	event_match, err := regexp.Match("(?i)Event", []byte(*str))
	if err != nil {
		log.Println("error parsing Event EntityType")
		return Undefined
	}

	if dataset_match {
		return Dataset
	}
	if dataService_match {
		return DataService
	}
	if concept_match {
		return Concept
	}
	if informationModel_match {
		return InformationModel
	}
	if publicService_match {
		return PublicService
	}
	if event_match {
		return Event
	}
	return Undefined
}

func (e EntityType) ToPath() *string {
	var path string
	switch e {
	case Dataset:
		path = "/datasets/"
	case DataService:
		path = "/dataservices/"
	case Concept:
		path = "/concepts/"
	case InformationModel:
		path = "/informationmodels/"
	case PublicService:
		path = "/public-services-and-events/"
	case Event:
		path = "/events/"
	default:
		return nil
	}
	return &path
}

func (e EntityType) String() string {
	switch e {
	case Dataset:
		return "Dataset"
	case DataService:
		return "Data Service"
	case Concept:
		return "Concept"
	case InformationModel:
		return "Information Model"
	case PublicService:
		return "Service"
	case Event:
		return "Event"
	default:
		return ""
	}
}

func (e EntityType) StringNb() string {
	switch e {
	case Dataset:
		return "Datasett"
	case DataService:
		return "Datatjeneste"
	case Concept:
		return "Begrep"
	case InformationModel:
		return "Informasjonsmodell"
	case PublicService:
		return "Tjeneste"
	case Event:
		return "Hendelse"
	default:
		return ""
	}
}

func (e EntityType) StringNbPlural() string {
	switch e {
	case Dataset:
		return "datasettet"
	case DataService:
		return "datatjenesten"
	case Concept:
		return "begrepet"
	case InformationModel:
		return "informasjonsmodellen"
	case PublicService:
		return "tjenesten"
	case Event:
		return "hendelsen"
	default:
		return ""
	}
}
