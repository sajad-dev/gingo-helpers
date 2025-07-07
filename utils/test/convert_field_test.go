package test

import (
	"encoding/json"
	"testing"

	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

type OtherTest struct {
	Image string `json:"image"`
}

type Education struct {
	Name string `json:"name"`
}

type Validation struct {
	Name       string      `json:"name"`
	Educations []Education `json:"educations"`
}

type Table struct {
	Name       string         `json:"name"`
	Educations datatypes.JSON `json:"educations"`
	Image      string         `json:"image"`
}

func TestConvertValidationToTable(t *testing.T) {
	ed := []Education{{Name: "hihi"}}
	table := Table{}
	err := utils.ConvertValidationToTable(&Validation{Name: "haha", Educations: ed}, &OtherTest{Image: "Gili Gili Gili"}, &table)
	assert.NoError(t, err)

	json, _ := json.Marshal(ed)
	assert.Equal(t, table, Table{Name: "haha", Educations: json, Image: "Gili Gili Gili"})
}
