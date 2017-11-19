package main

import (
	"errors"
	"strings"

	"github.com/containers/image/transports/alltransports"
	"github.com/containers/image/types"
	"github.com/urfave/cli"
)

func contextFromGlobalOptions(c *cli.Context, flagPrefix string) (*types.SystemContext, error) {
	ctx := &types.SystemContext{
		RegistriesDirPath:  c.GlobalString("registries.d"),
		ArchitectureChoice: c.GlobalString("override-arch"),
		OSChoice:           c.GlobalString("override-os"),
		DockerCertPath:     c.String(flagPrefix + "cert-dir"),
		// DEPRECATED: keep this here for backward compatibility, but override
		// them if per subcommand flags are provided (see below).
		DockerInsecureSkipTLSVerify: !c.GlobalBoolT("tls-verify"),
		OSTreeTmpDirPath:            c.String(flagPrefix + "ostree-tmp-dir"),
		OCISharedBlobDirPath:        c.String(flagPrefix + "shared-blob-dir"),
		DirForceCompress:            c.Bool(flagPrefix + "compress"),
	}
	if c.IsSet(flagPrefix + "tls-verify") {
		ctx.DockerInsecureSkipTLSVerify = !c.BoolT(flagPrefix + "tls-verify")
	}
	if c.IsSet(flagPrefix + "creds") {
		var err error
		ctx.DockerAuthConfig, err = getDockerAuth(c.String(flagPrefix + "creds"))
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func parseCreds(creds string) (string, string, error) {
	if creds == "" {
		return "", "", errors.New("credentials can't be empty")
	}
	up := strings.SplitN(creds, ":", 2)
	if len(up) == 1 {
		return up[0], "", nil
	}
	if up[0] == "" {
		return "", "", errors.New("username can't be empty")
	}
	return up[0], up[1], nil
}

func getDockerAuth(creds string) (*types.DockerAuthConfig, error) {
	username, password, err := parseCreds(creds)
	if err != nil {
		return nil, err
	}
	return &types.DockerAuthConfig{
		Username: username,
		Password: password,
	}, nil
}

// parseImage converts image URL-like string to an initialized handler for that image.
// The caller must call .Close() on the returned ImageCloser.
func parseImage(c *cli.Context) (types.ImageCloser, error) {
	imgName := c.Args().First()
	ref, err := alltransports.ParseImageName(imgName)
	if err != nil {
		return nil, err
	}
	ctx, err := contextFromGlobalOptions(c, "")
	if err != nil {
		return nil, err
	}
	return ref.NewImage(ctx)
}

// parseImageSource converts image URL-like string to an ImageSource.
// The caller must call .Close() on the returned ImageSource.
func parseImageSource(c *cli.Context, name string) (types.ImageSource, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}
	ctx, err := contextFromGlobalOptions(c, "")
	if err != nil {
		return nil, err
	}
	return ref.NewImageSource(ctx)
}
