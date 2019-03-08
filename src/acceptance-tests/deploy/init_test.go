package deploy_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/consul-release/src/acceptance-tests/testing/helpers"
	"github.com/pivotal-cf-experimental/bosh-test/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	config               helpers.Config
	boshClient           bosh.Client
	consulReleaseVersion string
)

func TestDeploy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "deploy")
}

var _ = BeforeSuite(func() {
	configPath, err := helpers.ConfigPath()
	Expect(err).NotTo(HaveOccurred())

	config, err = helpers.LoadConfig(configPath)
	Expect(err).NotTo(HaveOccurred())

	consulReleaseVersion = helpers.ConsulReleaseVersion()
	boshClient = bosh.NewClient(bosh.Config{
		URL:              config.BOSH.Target,
		Username:         config.BOSH.Username,
		Password:         config.BOSH.Password,
		AllowInsecureSSL: true,
		UAA: true,
	})
})

func lockedDeployments() ([]string, error) {
	var lockNames []string
	locks, err := boshClient.Locks()
	if err != nil {
		return []string{}, err
	}
	for _, lock := range locks {
		lockNames = append(lockNames, lock.Resource[0])
	}
	return lockNames, nil
}
