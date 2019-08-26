package typirelease

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v27/github"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/git"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
	"golang.org/x/oauth2"
)

// ReleaseDistribution to release the distribution
func ReleaseDistribution(rel typictx.Release, force bool) (binaries, changeLogs []string, err error) {
	if len(rel.GoOS) < 0 {
		err = errors.New("Missing 'GoOS' in Typical Context")
		return
	}
	if len(rel.GoArch) < 0 {
		err = errors.New("Missing 'GoArch' in Typical Context")
		return
	}
	git.Fetch()
	defer git.Fetch()
	status := git.Status()
	if !force && status != "" {
		err = fmt.Errorf("Please commit changes first:\n%s", status)
		return
	}
	latestTag := git.LatestTag()
	if !force && latestTag == rel.Tag() {
		log.Infof("%s already released", latestTag)
		return
	}
	changeLogs = git.Logs(latestTag)
	if !force && len(changeLogs) < 1 {
		log.Info("No change to be released")
		return
	}
	for _, changeLog := range changeLogs {
		log.Infof("Change Log: %s", changeLog)
	}
	mainPackage := typienv.AppMainPackage()
	for _, os1 := range rel.GoOS {
		for _, arch := range rel.GoArch {
			binary := rel.ReleaseBinary(os1, arch)
			binaryPath := fmt.Sprintf("%s/%s", typienv.Release(), binary)
			log.Infof("Create release binary for %s/%s: %s", os1, arch, binaryPath)
			// TODO: support cgo
			err = bash.GoBuild(binaryPath, mainPackage, "GOOS="+os1, "GOARCH="+arch)
			if err != nil {
				return
			}
			binaries = append(binaries, binary)
		}
	}
	return
}

// GithubRelease for github release
func GithubRelease(binaries, changeLogs []string, rel typictx.Release) (err error) {
	if rel.Github == nil {
		return
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Environment 'GITHUB_TOKEN' is missing")
	}
	owner := rel.Github.Owner
	repo := rel.Github.RepoName
	ctx0 := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx0, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	releaser := githubReleaser{rel}
	if releaser.IsReleased(ctx0, client.Repositories) {
		log.Infof("Release for %s/%s (%s) already exist", owner, repo, rel.Tag())
		return
	}
	log.Info("Generate release note")
	var releaseNote strings.Builder
	for _, changelog := range changeLogs {
		if !ignoring(changelog) {
			releaseNote.WriteString(changelog)
			releaseNote.WriteString("\n")
		}
	}
	log.Infof("Create github release for %s/%s", owner, repo)
	var release *github.RepositoryRelease
	release, err = releaser.CreateRelease(ctx0, client.Repositories, releaseNote.String())
	if err != nil {
		return
	}
	for _, binary := range binaries {
		log.Infof("Upload asset: %s", binary)
		err = releaser.Upload(ctx0, client.Repositories, *release.ID, binary)
		if err != nil {
			return
		}
	}
	return
}

func cleanMessage(message string) string {
	iCoAuthor := strings.Index(message, "Co-Authored-By")
	if iCoAuthor > 0 {
		message = message[0:strings.Index(message, "Co-Authored-By")]
	}
	message = strings.TrimSpace(message)
	return message
}

func ignoring(changelog string) bool {
	message := cleanMessage(changelog[7:])
	lowerMessage := strings.ToLower(message)
	return strings.HasPrefix(lowerMessage, "merge") ||
		strings.HasPrefix(lowerMessage, "bump") ||
		strings.HasPrefix(lowerMessage, "revision") ||
		strings.HasPrefix(lowerMessage, "generate") ||
		strings.HasPrefix(lowerMessage, "wip")
}
