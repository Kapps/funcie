#!/bin/bash

# Ensure we're in the root directory of the repository
if [ ! -f VERSION ]; then
    echo "ERROR: VERSION file not found. Please run this script from the root directory of the repository."
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "ERROR: There are uncommitted changes in the repository. Commit or stash them before running this script."
    exit 1
fi

CURRENT_VERSION=$(cat VERSION)

# Split the version number into its components
IFS='.' read -ra ADDR <<< "$CURRENT_VERSION"

MAJOR=${ADDR[0]}
MINOR=${ADDR[1]}
PATCH=${ADDR[2]}

# Determine the type of version bump
BUMP_TYPE=$1

# Bump the version based on the specified type
case $BUMP_TYPE in
    major)
        NEW_VERSION="$((MAJOR + 1)).0.0"
        ;;
    minor)
        NEW_VERSION="${MAJOR}.$((MINOR + 1)).0"
        ;;
    patch)
        NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH + 1))"
        ;;
    *)
        echo "ERROR: Invalid argument. Please specify 'major', 'minor', or 'patch' as the bump type."
        exit 1
        ;;
esac

echo $NEW_VERSION > VERSION

git add VERSION
git commit -m "Bump version to $NEW_VERSION"
git tag $NEW_VERSION


read -p "Version bumped from $CURRENT_VERSION to $NEW_VERSION. Do you want to push the changes? [y/N] " confirm
case $confirm in
    [yY][eE][sS]|[yY])
        # Push changes to the repository
        git push origin main
        git push origin $NEW_VERSION
        echo "Version bumped to $NEW_VERSION and pushed to the repository."
        ;;
    *)
        echo "Push aborted. The commit and tag are still available locally. Don't forget to push the tag along with the commit."
        exit 1
        ;;
esac
