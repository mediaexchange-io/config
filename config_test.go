package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/mediaexchange-io/assert"
)

type basicList []interface{}

type basicMap struct {
	_integer int     `json:"integer"`
	_float   float64 `json:"float"`
	_string  string  `json:"string"`
	_boolean bool    `json:"boolean"`
}

type testMap struct {
	*basicMap
	_sublist basicList `json:"list"`
	_submap  basicMap  `json:"map"`
}

type testConfig struct {
	_map testMap `json:"map"`
}

func TestFromJsonFile(t *testing.T) {
	var k = new(testConfig)
	var err = FromFile("config_test.json", k)
	assert.With(t).
		That(err).
		IsOk()
}

func TestFromYamlFile(t *testing.T) {
	var k = new(testConfig)
	var err = FromFile("config_test.yaml", k)
	assert.With(t).
		That(err).
		IsOk()
}

func TestWithEnvOverride(t *testing.T) {
	expected := 98765
	_ = os.Setenv("NUM", strconv.Itoa(expected))

	// Nested structure to test.
	type testStruct struct {
		Str    string `json:"str"`
		Num    int    `json:"num" env:"NUM"`
		Nested struct {
			Str string `json:"str"`
			Num int    `json:"num" env:"NUM"`
		}
	}

	jsonString := "{\"str\":\"foobar\",\"num\":-1111,\"nested\":{\"str\":\"deadbeef\",\"num\":-2222}}"
	ts := new(testStruct)
	err := fromJson([]byte(jsonString), ts)

	Assert := assert.With(t)
	Assert.
		That(err).
		IsOk()
	Assert.
		That(ts.Num).
		IsEqualTo(98765)
	Assert.
		That(ts.Nested.Num).
		IsEqualTo(98765)
}

func validate(k *testConfig, t *testing.T) {
	Assert := assert.With(t)
	Assert.
		That(k._map._boolean).
		IsEqualTo(true)
	Assert.
		That(k._map._float).
		IsEqualTo(3.14159)
	Assert.
		That(k._map._integer).
		IsEqualTo(1234)
	Assert.
		That(k._map._string).
		IsEqualTo("string")
	Assert.
		That(k._map._sublist[0]).
		IsEqualTo(true)
	Assert.
		That(k._map._sublist[1]).
		IsEqualTo("string")
	Assert.
		That(k._map._sublist[2]).
		IsEqualTo(3.14159)
	Assert.
		That(k._map._sublist[3]).
		IsEqualTo(1234)
	Assert.
		That(k._map._submap._boolean).
		IsEqualTo(true)
	Assert.
		That(k._map._submap._float).
		IsEqualTo(3.14159)
	Assert.
		That(k._map._submap._integer).
		IsEqualTo(1234)
	Assert.
		That(k._map._submap._string).
		IsEqualTo("string")
}

// Example structure that matches config_example.yaml.
// NOTE: if the field names match the names used in the YAML file, the json
// struct tags are not necessary.
type AppConfig struct {
	Server struct {
		Port int16 `env:"PORT"`
	}
	Database struct {
		Driver   string
		Hostname string
		Port     int16
		Username string
		Password string
		Name     string
	}
}

// Reads the AppConfig data from config_example.yaml and prints it.
func ExampleFromFile() {
	os.Setenv("PORT", "9000")
	var c = new(AppConfig)
	if err := FromFile("config_example.yaml", c); err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("%+v\n", *c)
	// OUTPUT: {Server:{Port:9000} Database:{Driver:postgres Hostname:localhost Port:5432 Username:postgres Password:dummy Name:my_database}}
}
