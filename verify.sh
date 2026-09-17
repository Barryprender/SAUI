#!/usr/bin/env bash
# Every gate that decides whether this repository is shippable, in one place.
# CI runs this file, not a copy of it, so a green run locally means the same
# thing as a green run on GitHub.
#
#   ./verify.sh               run every gate
#   ./verify.sh --update-sbom regenerate sbom.json, then run every gate
set -euo pipefail

# Pinned so the gates' own output cannot change without a commit. govulncheck
# is pinned too; only its vulnerability database is expected to move.
TEMPL_VERSION=v0.3.1001
CYCLONEDX_VERSION=v1.9.0
GOVULNCHECK_VERSION=v1.8.0

export PATH="$(go env GOPATH)/bin:$PATH"

fail() { echo "FAIL: $*" >&2; exit 1; }
step() { echo; echo "== $* =="; }

# Drop the fields that change on every run without the software changing:
# the generation timestamp, the main module's git-derived pseudo-version, and
# the hashes cyclonedx-gomod records of its own binary. That last one is not
# cosmetic: the generator is installed with `go install`, so its binary differs
# between a Windows developer machine and the Linux runner, and the gate could
# never pass on both at once. Those hashes describe the machine that ran the
# tool, not the software being described.
#
# Dependency versions and hashes are left alone — those are the point. The
# eight-space anchor below matches only the tools entry; component hashes sit
# one level shallower, at six spaces, and are untouched.
normalise_sbom() {
  sed -e '/"timestamp"/d' \
      -e '/^        "hashes": \[$/,/^        \],$/d' \
      -e '/"metadata"/,/"components"/ s/v0\.0\.0-[0-9]\{14\}-[0-9a-f]\{12\}/v0.0.0-devel/g' \
      -e 's#golang/saui@v0\.0\.0-[0-9]\{14\}-[0-9a-f]\{12\}#golang/saui@v0.0.0-devel#g' \
      "$1"
}

# The SBOM must describe the container that ships, not the developer's machine.
# Build constraints steer Go's module selection, so an SBOM generated on
# windows/amd64 is an SBOM of a binary nobody deploys. The constraints go on
# the generator only — the toolchain itself must still build for this host.
generate_sbom() {
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    cyclonedx-gomod app -json -main . -noserial -output "$1" . >/dev/null
}

step "install pinned tooling"
go install "github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}"
go install "github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@${CYCLONEDX_VERSION}"
go install "golang.org/x/vuln/cmd/govulncheck@${GOVULNCHECK_VERSION}"
echo "ok"

step "gofmt"
# _templ.go files are generated; templ owns their formatting, not gofmt.
unformatted=$(gofmt -l . | grep -v '_templ\.go$' || true)
[ -z "$unformatted" ] || fail "not gofmt'd:"$'\n'"$unformatted"
echo "ok"

step "go vet"
go vet ./...
echo "ok"

step "templ generation is current"
# The silent lie in this repository: edit a .templ file, forget to regenerate,
# and the whole suite passes against the previous version of the markup.
templ generate >/dev/null
# Unstaged drift only: regenerating must not change what is already recorded.
# git status would also flag a newly added file whose generated form is correct.
stale=$(git diff --name-only -- '*_templ.go')
[ -z "$stale" ] || fail "regenerating changed these, so the recorded markup is behind its source:"$'\n'"$stale"$'\n'"review the regenerated files and include them in the commit"
echo "ok"

step "tests"
go test ./... -count=1 \
  -coverpkg=./actions/...,./handlers/...,./middleware/...,./projections/...,./statestore/... \
  -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -1

step "govulncheck"
govulncheck ./...

step "SBOM is current"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
if [ "${1:-}" = "--update-sbom" ]; then
  generate_sbom "$tmp/raw.json"
  normalise_sbom "$tmp/raw.json" > sbom.json
  echo "sbom.json regenerated"
else
  [ -f sbom.json ] || fail "sbom.json is missing; run: ./verify.sh --update-sbom"
  generate_sbom "$tmp/raw.json"
  normalise_sbom "$tmp/raw.json" > "$tmp/fresh.json"
  diff -u sbom.json "$tmp/fresh.json" \
    || fail "sbom.json is stale; run: ./verify.sh --update-sbom"
  echo "ok"
fi

echo
echo "all gates passed"
