package main

import (
	"fmt"
	"io"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"

	autotag "github.com/rupor-github/fork-autotag"
)

// Options holds the CLI args
type Options struct {
	JustVersion         bool   `short:"n" description:"Just output the next version, don't autotag"`
	Verbose             bool   `short:"v" description:"Enable verbose logging"`
	Branch              string `short:"b" long:"branch" description:"Git branch to scan (if not set picked up from the repository)" default:""`
	RepoPath            string `short:"r" long:"repo" description:"Path to the repo" default:"./" `
	PreReleaseName      string `short:"p" long:"pre-release-name" description:"create a pre-release tag"`
	PreReleaseTimestamp string `short:"T" long:"pre-release-timestamp" description:"create a pre-release tag and append a timestamp (can be: datetime|epoch)"`
	PreReleaseAttempt   int    `short:"A" long:"pre-release-attempt" description:"create a pre-release tag and append attempt number to it with '.' character (can be positive number or 0)" default:"-1"`
	BuildMetadata       string `short:"m" long:"build-metadata" description:"optional SemVer build metadata to append to the version with '+' character"`
	Scheme              string `short:"s" long:"scheme" description:"The commit message scheme to use (can be: autotag|conventional)" default:"autotag"`
	NoVersionPrefix     bool   `short:"e" long:"empty-version-prefix" description:"Do not prepend v to version tag"`
	Push                bool   `short:"P" long:"push-to-origin" description:"Push new tag to remote"`
	Remote              string `short:"R" long:"remote" description:"Remote to push tag to" default:"origin"`
	NoGitStatusCheck    bool   `short:"C" long:"no-git-status-check" description:"Bypass checking worktree status"`
}

var opts Options

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !flags.WroteHelp(err) {
			log.Println(err)
		}
		os.Exit(1)
	}
}

func main() {
	log.SetOutput(io.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	r, err := autotag.NewRepo(autotag.GitRepoConfig{
		RepoPath:                  opts.RepoPath,
		Branch:                    opts.Branch,
		PreReleaseName:            opts.PreReleaseName,
		PreReleaseTimestampLayout: opts.PreReleaseTimestamp,
		PreReleaseAttempt:         opts.PreReleaseAttempt,
		BuildMetadata:             opts.BuildMetadata,
		Scheme:                    opts.Scheme,
		Prefix:                    !opts.NoVersionPrefix,
		Remote:                    opts.Remote,
		Push:                      opts.Push,
		Check:                     !opts.NoGitStatusCheck,
	})
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Error initializing: " + err.Error())
		os.Exit(1)
	}

	// Tag unless asked otherwise
	if !opts.JustVersion {
		err = r.AutoTag()
		if err != nil {
			log.SetOutput(os.Stderr)
			log.Println("Error auto updating version: " + err.Error())
			os.Exit(1)
		}
	}

	if opts.NoVersionPrefix {
		fmt.Println(r.LatestVersion())
	} else {
		fmt.Println("v" + r.LatestVersion())
	}

	// TODO:(jnelson) Add -major -minor -patch flags for force bumps Fri Sep 11 10:04:20 2015
	os.Exit(0)
}
