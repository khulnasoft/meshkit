package generators

import (
	"fmt"

	"github.com/khulnasoft/meshkit/generators/artifacthub"
	"github.com/khulnasoft/meshkit/generators/github"
	"github.com/khulnasoft/meshkit/models"
	"github.com/khulnasoft/meshkit/utils"
)

const (
	artifactHub = "artifacthub"
	gitHub      = "github"
)

func NewGenerator(registrant, url, packageName string) (models.PackageManager, error) {
	registrant = utils.ReplaceSpacesAndConvertToLowercase(registrant)
	switch registrant {
	case artifactHub:
		return artifacthub.ArtifactHubPackageManager{
			PackageName: packageName,
			SourceURL:   url,
		}, nil
	case gitHub:
		return github.GitHubPackageManager{
			PackageName: packageName,
			SourceURL:   url,
		}, nil
	}
	return nil, ErrUnsupportedRegistrant(fmt.Errorf("generator not implemented for the registrant %s", registrant))
}
