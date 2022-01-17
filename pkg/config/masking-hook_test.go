package config

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/sirupsen/logrus"
)

func TestMaskingHook(t *testing.T) {
	spec.Run(t, "MaskingHook", testMaskingHook, spec.Report(report.Terminal{}))
}

func testMaskingHook(t *testing.T, when spec.G, it spec.S) {
	var h *MaskingHook

	it.Before(func() {
		RegisterTestingT(t)
		h = &MaskingHook{
			MaskedValue: "***",
		}
	})

	it("does not add empty strings to the masker", func() {
		h.AddToMaskList("")
		h.AddToMaskList("  ")
		h.AddToMaskList("\t")
		h.AddToMaskList("\n")

		Expect(len(h.toMask)).To(Equal(0))
	})

	it("adds valid strings to the masker", func() {
		h.AddToMaskList("secret")
		h.AddToMaskList("value")

		Expect(len(h.toMask)).To(Equal(2))
		Expect(h.toMask[0]).To(Equal("secret"))
		Expect(h.toMask[1]).To(Equal("value"))
	})

	it("does not modify messages when there is nothing to mask", func() {
		entry := &logrus.Entry{
			Message: "a secret message",
		}

		h.Fire(entry)

		Expect(entry.Message).To(Equal("a secret message"))
	})

	it("removes all instances of secrets from messages", func() {
		entry := &logrus.Entry{
			Message: "a secret message that is a secret",
		}

		h.toMask = append(h.toMask, "secret")
		h.Fire(entry)

		Expect(entry.Message).To(Equal("a *** message that is a ***"))
	})

	it("removes all instances of secrets from fields", func() {
		entry := &logrus.Entry{}
		entry = entry.WithField("a-field", "a secret message that is a secret")

		h.toMask = append(h.toMask, "secret")
		h.Fire(entry)

		Expect(entry.Data["a-field"]).To(Equal("a *** message that is a ***"))
	})

	// it("returns an error if the manifest file does not exist", func() {
	// 	p.configuration.Paths.DeploymentManifestFile = "nothing.yaml"
	// 	err := p.Execute(fakeContextProvider, nil, nil)
	// 	assert.NotNil(t, err)
	// })

	// it("returns an error if the manifest file is not yaml", func() {
	// 	afero.WriteFile(config.Appfs, testManifestFilePath, []byte("not yaml"), 0755)
	// 	err := p.Execute(fakeContextProvider, nil, nil)
	// 	assert.NotNil(t, err)
	// })
}
