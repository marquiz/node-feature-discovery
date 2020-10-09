---
name: New Release
about: Propose a new release
title: Release v0.x.0
assignees: marquiz

---

## Release Checklist
<!--
Please do not remove items from the checklist
-->
- [ ] All [OWNERS](https://github.com/kubernetes-sigs/node-feature-discovery/blob/master/OWNERS) must LGTM the release proposal
- [ ] Verify that the changelog in this issue is up-to-date
- [ ] For major releases (v0.$MAJ.0), an OWNER creates a release branch with
      `git branch release-0.$MAJ master`
- [ ] Prepare `release-0.$MAJ` release branch
  - [ ] Update the deployment templates to use the new tagged container image:
        `sed s"|image: .*|image: k8s.gcr.io/nfd/node-feature-discovery:$VERSION|" -i *yaml.template`
  - [ ] Update quick start instructions in README.md to use $VERSION
  - [ ] Update version configuration in `docs/_config.yml`:
        - set `version: $VERSION`
        - add $VERSION to `versions:` list
  - [ ] An OWNER runs
       `git tag -s $VERSION`
        and inserts the changelog into the tag description.
  - [ ] An OWNER pushes the release branch with
        `git push release-0.$MAJ`
        This will trigger build of the documentation and publish it at https://kubernetes-sigs.github.io/node-feature-discovery/0.$MAJ/
- [ ] An OWNER pushes the tag with
      `git push $VERSION`
      This will trigger prow to build and publish a staging container image
      `gcr.io/k8s-staging-nfd/node-feature-discovery:$VERSION`
- [ ] Submit a PR against [k8s.io](https://github.com/kubernetes/k8s.io), updating `k8s.gcr.io/images/k8s-staging-nfd/images.yaml` to promote the container image to production
- [ ] Wait for the PR to be merged and verify that the image (`k8s.gcr.io/nfd/node-feature-discovery:$VERSION`) is available.
- [ ] Write the change log into the [Github release info](https://github.com/kubernetes-sigs/node-feature-discovery/releases).
- [ ] Add a link to the tagged release in this issue.
- [ ] Update `index.html` in `gh-pages` branch to point to the latest release
- [ ] Send an announcement email to `kubernetes-dev@googlegroups.com` with the subject `[ANNOUNCE] node-feature-discovery $VERSION is released`
- [ ] Add a link to the release announcement in this issue
- [ ] Close this issue


## Changelog
<!--
Describe changes since the last release here.
-->
