/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package target

import (
	"context"

	"github.com/gardener/gardenctl-v2/pkg/config"
)

var Merge = merge

// ResolvePatternOverlay exposes managerImpl.resolvePatternOverlay for tests.
func ResolvePatternOverlay(ctx context.Context, m Manager, tf TargetFlags, persistedTarget Target, value string) (TargetFlags, error) {
	return m.(*managerImpl).resolvePatternOverlay(ctx, tf, persistedTarget, value)
}

// ResolveAccessLevelForWrite exposes managerImpl.resolveAccessLevelForWrite for tests.
func ResolveAccessLevelForWrite(ctx context.Context, m Manager, t Target) (config.KubeconfigAccessLevel, error) {
	return m.(*managerImpl).resolveAccessLevelForWrite(ctx, t)
}

// EffectiveAccessLevel exposes managerImpl.effectiveAccessLevel for tests.
func EffectiveAccessLevel(m Manager, t Target, scope AccessScope) config.KubeconfigAccessLevel {
	return m.(*managerImpl).effectiveAccessLevel(t, scope)
}
