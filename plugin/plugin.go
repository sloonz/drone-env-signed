// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/environ"

	"github.com/gbrlsnchs/jwt/v3"
)

func New(secret string, algorithm jwt.Algorithm) environ.Plugin {
	return &plugin{
		secret:    secret,
		algorithm: algorithm,
	}
}

type plugin struct {
	secret    string
	algorithm jwt.Algorithm
}

type repoPayload struct {
	jwt.Payload
	drone.Repo
}

type buildPayload struct {
	jwt.Payload
	drone.Build
}

func (p *plugin) List(ctx context.Context, req *environ.Request) (map[string]string, error) {
	now := jwt.NumericDate(time.Now())

	repo := repoPayload{
		Payload: jwt.Payload{
			IssuedAt: now,
		},
		Repo: req.Repo,
	}
	signedRepo, err := jwt.Sign(repo, p.algorithm)
	if err != nil {
		return nil, err
	}

	build := buildPayload{
		Payload: jwt.Payload{
			IssuedAt: now,
		},
		Build: req.Build,
	}
	signedBuild, err := jwt.Sign(build, p.algorithm)
	if err != nil {
		return nil, err
	}

	env := map[string]string{
		"DRONE_SIGNED_REPO":  string(signedRepo),
		"DRONE_SIGNED_BUILD": string(signedBuild),
	}

	return env, nil
}
