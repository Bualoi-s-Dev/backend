package testing_utils

import (
	"os"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func SetupGodog(file string, InitializeScenario func(*godog.ScenarioContext)) *godog.TestSuite {
	filePath := "./testing/features/" + file
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{filePath},
		Output: colors.Colored(os.Stdout),
	}
	testSuite := godog.TestSuite{
		Name:                "api",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}
	return &testSuite
}
