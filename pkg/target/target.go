/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package target

import (
	"errors"

	gardencore "github.com/gardener/gardener/pkg/apis/core"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clientgarden "github.com/gardener/gardenctl-v2/internal/client/garden"
	"github.com/gardener/gardenctl-v2/pkg/config"
)

/*
	Targets are either

	garden -> seed -> shoot
	garden -> project -> shoot
*/

// Target represents the Kubernetes cluster/namespace which should
// be the target for user operations in gardenctl. It works similar
// to the context defined in a kubeconfig.
type Target interface {
	// GardenName returns the currently targeted garden cluster name.
	GardenName() string
	// ProjectName returns the currently targeted project name.
	ProjectName() string
	// SeedName returns the currently targeted seed cluster name.
	SeedName() string
	// ShootName returns the currently targeted shoot cluster name.
	ShootName() string
	// ControlPlane returns true if shoot control plane is targeted.
	ControlPlane() bool
	// AccessLevel returns the kubeconfig access level resolved for this target
	// when it was written (empty when gardenlogin's own default applies).
	AccessLevel() config.KubeconfigAccessLevel
	// WithGardenName returns a copy of the target with the garden name updated.
	// The returned target can be invalid.
	WithGardenName(name string) Target
	// WithProjectName returns a copy of the target with the project name updated.
	// The returned target can be invalid.
	WithProjectName(name string) Target
	// WithSeedName returns a copy of the target with the seed name updated.
	// The returned target can be invalid.
	WithSeedName(name string) Target
	// WithShootName returns a copy of the target with the shoot name updated.
	// The returned target can be invalid.
	WithShootName(name string) Target
	// WithControlPlane returns a copy of the target with the control plane flag updated.
	// The returned target can be invalid.
	WithControlPlane(controlPlane bool) Target
	// WithAccessLevel returns a copy of the target with the access level updated.
	// The returned target can be invalid.
	WithAccessLevel(level config.KubeconfigAccessLevel) Target
	// Validate checks for semantical correctness of the target, without
	// actually connecting to the targeted clusters.
	Validate() error
	// IsEmpty returns true if all values of the target are empty
	IsEmpty() bool
	// AsListOption returns the target as list option
	AsListOption() client.ListOption

	// DeepCopy returns a deep copy of the target
	DeepCopy() Target
}

type targetImpl struct {
	Garden           string                       `json:"garden,omitempty"`
	Project          string                       `json:"project,omitempty"`
	Seed             string                       `json:"seed,omitempty"`
	Shoot            string                       `json:"shoot,omitempty"`
	ControlPlaneFlag bool                         `json:"controlPlane,omitempty"`
	AccessLevelValue config.KubeconfigAccessLevel `json:"accessLevel,omitempty"`
}

var _ Target = &targetImpl{}

// NewTarget returns a new target. This function does not perform any validation,
// so the returned target can be invalid.
func NewTarget(gardenName, projectName, seedName, shootName string) Target {
	return &targetImpl{
		Garden:           gardenName,
		Project:          projectName,
		Seed:             seedName,
		Shoot:            shootName,
		ControlPlaneFlag: false,
	}
}

// copy returns a copy of the target (all fields are value types).
func (t *targetImpl) copy() *targetImpl {
	c := *t
	return &c
}

// Validate checks that the target is not malformed and all required
// steps are configured correctly.
func (t *targetImpl) Validate() error {
	if len(t.Project) > 0 && len(t.Seed) > 0 && len(t.Shoot) == 0 {
		return errors.New("seed and project must not be configured at the same time without a shoot")
	}

	if t.ControlPlaneFlag && t.Shoot == "" {
		return errors.New("control-plane requires a shoot target")
	}

	if t.ControlPlaneFlag && t.Garden == "" {
		return errors.New("control-plane requires a garden target")
	}

	if t.Garden != "" {
		if err := config.ValidateGardenName(t.Garden); err != nil {
			return err
		}
	}

	if t.Shoot != "" {
		if errs := apivalidation.NameIsDNSLabel(t.Shoot, false); len(errs) > 0 {
			return errors.New("shoot name is invalid: " + errs[0])
		}
	}

	if t.Project != "" {
		if errs := apivalidation.NameIsDNSSubdomain(t.Project, false); len(errs) > 0 {
			return errors.New("project name is invalid: " + errs[0])
		}
	}

	if t.Seed != "" {
		if errs := apivalidation.NameIsDNSLabel(t.Seed, false); len(errs) > 0 {
			return errors.New("seed name is invalid: " + errs[0])
		}
	}

	if err := t.AccessLevelValue.Validate(); err != nil {
		return err
	}

	return nil
}

// GardenName returns the currently targeted garden cluster name.
func (t *targetImpl) GardenName() string {
	return t.Garden
}

// ProjectName returns the currently targeted project name.
func (t *targetImpl) ProjectName() string {
	return t.Project
}

// SeedName returns the currently targeted seed cluster name.
func (t *targetImpl) SeedName() string {
	return t.Seed
}

// ShootName returns the currently targeted shoot cluster name.
func (t *targetImpl) ShootName() string {
	return t.Shoot
}

// ControlPlane returns true if shoot control plane is targeted.
func (t *targetImpl) ControlPlane() bool {
	return t.ControlPlaneFlag
}

// AccessLevel returns the kubeconfig access level resolved for this target.
func (t *targetImpl) AccessLevel() config.KubeconfigAccessLevel {
	return t.AccessLevelValue
}

// WithGardenName returns a copy of the target with the garden name updated.
// The returned target can be invalid.
func (t *targetImpl) WithGardenName(name string) Target {
	c := t.copy()
	c.Garden = name

	return c
}

// WithProjectName returns a copy of the target with the project name updated.
// The returned target can be invalid.
func (t *targetImpl) WithProjectName(name string) Target {
	c := t.copy()
	c.Project = name

	return c
}

// WithSeedName returns a copy of the target with the seed name updated.
// The returned target can be invalid.
func (t *targetImpl) WithSeedName(name string) Target {
	c := t.copy()
	c.Seed = name

	return c
}

// WithShootName returns a copy of the target with the shoot name updated.
// The returned target can be invalid.
func (t *targetImpl) WithShootName(name string) Target {
	c := t.copy()
	c.Shoot = name

	return c
}

// WithControlPlane returns a copy of the target with the control plane flag updated.
// The returned target can be invalid.
func (t *targetImpl) WithControlPlane(controlPlane bool) Target {
	c := t.copy()
	c.ControlPlaneFlag = controlPlane

	return c
}

// WithAccessLevel returns a copy of the target with the access level updated.
// The returned target can be invalid.
func (t *targetImpl) WithAccessLevel(level config.KubeconfigAccessLevel) Target {
	c := t.copy()
	c.AccessLevelValue = level

	return c
}

func (t *targetImpl) IsEmpty() bool {
	return t.Garden == "" && t.Project == "" && t.Seed == "" && t.Shoot == ""
}

func (t *targetImpl) AsListOption() client.ListOption {
	opt := clientgarden.ProjectFilter{}

	if t.ShootName() != "" {
		opt["metadata.name"] = t.ShootName()
	}

	if t.ProjectName() != "" {
		opt["project"] = t.ProjectName()
	} else if t.SeedName() != "" {
		opt[gardencore.ShootSeedName] = t.SeedName()
	}

	return opt
}

func (t *targetImpl) DeepCopy() Target {
	return t.copy()
}
