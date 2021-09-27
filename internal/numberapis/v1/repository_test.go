package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	numv1 "github.com/tuanden0/generate_number/proto/gen/go/v1/number"
)

func Test_repoManager_Get(t *testing.T) {

	// Init
	minA := int64(-10)
	maxB := int64(10)
	minC := int64(-20)
	maxD := int64(30)

	in := &numv1.NumberServiceGetRequest{
		NumA: minA,
		NumB: maxB,
		NumC: minC,
		NumD: maxD,
	}

	r := NewRepoManager()

	// Execute
	res, err := r.Get(context.Background(), in)

	// Assertion
	assert.Nil(t, err)
	assert.NotNil(t, res)

	if (res < int64(minA) || res < int64(minC)) && (res > int64(maxB)) || res > int64(maxD) {
		t.Errorf(
			"generate number failed, got %v not in range %v, %v, %v, %v",
			res, minA, maxB, minC, maxD)
	}
}
