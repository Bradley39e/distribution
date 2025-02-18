package testutil

import (
	"fmt"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/internal/dcontext"
	"github.com/distribution/distribution/v3/manifest/manifestlist"
	"github.com/distribution/distribution/v3/manifest/schema2"
	"github.com/opencontainers/go-digest"
)

// MakeManifestList constructs a manifest list out of a list of manifest digests
func MakeManifestList(blobstatter distribution.BlobStatter, manifestDigests []digest.Digest) (*manifestlist.DeserializedManifestList, error) {
	ctx := dcontext.Background()

	manifestDescriptors := make([]manifestlist.ManifestDescriptor, 0, len(manifestDigests))
	for _, manifestDigest := range manifestDigests {
		descriptor, err := blobstatter.Stat(ctx, manifestDigest)
		if err != nil {
			return nil, err
		}
		platformSpec := manifestlist.PlatformSpec{
			Architecture: "atari2600",
			OS:           "CP/M",
			Variant:      "ternary",
			Features:     []string{"VLIW", "superscalaroutoforderdevnull"},
		}
		manifestDescriptor := manifestlist.ManifestDescriptor{
			Descriptor: descriptor,
			Platform:   platformSpec,
		}
		manifestDescriptors = append(manifestDescriptors, manifestDescriptor)
	}

	return manifestlist.FromDescriptors(manifestDescriptors)
}

// MakeSchema2Manifest constructs a schema 2 manifest from a given list of digests and returns
// the digest of the manifest
func MakeSchema2Manifest(repository distribution.Repository, digests []digest.Digest) (distribution.Manifest, error) {
	ctx := dcontext.Background()
	blobStore := repository.Blobs(ctx)

	var configJSON []byte

	d, err := blobStore.Put(ctx, schema2.MediaTypeImageConfig, configJSON)
	if err != nil {
		return nil, fmt.Errorf("unexpected error storing content in blobstore: %v", err)
	}
	builder := schema2.NewManifestBuilder(d, configJSON)
	for _, digest := range digests {
		builder.AppendReference(distribution.Descriptor{Digest: digest})
	}

	mfst, err := builder.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("unexpected error generating manifest: %v", err)
	}

	return mfst, nil
}
