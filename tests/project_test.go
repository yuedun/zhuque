package tests

import (
	"encoding/json"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestUpdateProject(t *testing.T) {
	var configStr = `{
		"apps" : [{
			"script": "index.js",
			"watch": "."
		}],
		
		"deploy" : {
			"production" : {
				"user" : "node",
				"host" : ["212.83.163.1", "212.83.163.2", "212.83.163.3"],
				"ref"  : "origin/master",
				"repo" : "git@github.com:repo.git",
				"path" : "/var/www/production",
				"pre-setup" : "echo 'commands or local script path to be run on the host before the setup process starts'",
				"post-setup": "echo 'commands or a script path to be run on the host after cloning the repo'",
				"post-deploy" : "pm2 startOrRestart ecosystem.json --env production",
				"pre-deploy-local" : "echo 'This is a local executed command'"
			  }
		}
	}`
	var config map[string]interface{}
	// var deployConfig project.DeployConfig
	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		t.Error(err)
	}
	production := config["deploy"].(map[string]interface{})["production"].(map[string]interface{})

	t.Log(production["pre-setup"])
}

func TestFileExists(t *testing.T) {
	i, err := os.Stat("d:/workspace/testDeploy")
	t.Log(i, err)
}
