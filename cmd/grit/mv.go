package main

import (
	"errors"
	"os"
	"path"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/jmalloc/grit"
	"github.com/jmalloc/grit/index"
	"github.com/urfave/cli"
)

func mv(cfg grit.Config, idx *index.Index, c *cli.Context) error {
	src, err := dirFromArg(c, 0)
	if err != nil {
		return err
	}

	rem, ok, err := chooseRemote(cfg, c, src, func(_ *config.RemoteConfig, ep *transport.Endpoint) string {
		return " --> " + grit.EndpointToDir(cfg.Clone.Root, ep)
	})

	if err != nil {
		return err
	} else if !ok {
		return errSilentFailure
	}

	ep, _, err := grit.EndpointFromRemote(rem)
	if err != nil {
		return err
	}

	dst := grit.EndpointToDir(cfg.Clone.Root, ep)

	return moveClone(cfg, idx, c, src, dst)
}

func moveClone(cfg grit.Config, idx *index.Index, c *cli.Context, src, dst string) error {
	writeln(c, dst)

	if src == dst {
		return nil
	}

	if wd, _ := os.Getwd(); wd == src {
		exec(c, "cd", dst)
	}

	_, err := os.Stat(dst)
	if err == nil {
		return errors.New("destination directory already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}

	if err := idx.Add(dst); err != nil {
		return err
	}

	return idx.Remove(src)
}
