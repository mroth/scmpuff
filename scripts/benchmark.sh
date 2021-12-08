#!/usr/bin/env bash
#
# Quick script to generate a big git status change to compare performance of
# scmpuff with scm_breeze.
#
# For this just populates a fake git repo so I can do some manual timing on it.
#

DIR=/tmp/scmpuff/bench

# make bench directory
rm -rf ${DIR}
mkdir -p ${DIR}

# make it a git repo
cd ${DIR}
git init .

k=25

# make 25 files in git history (no changes)
for ((i=0;i<k;i++)); do
  F=$(mktemp foo.XXXX)
  echo "X" > ${F}
  git add ${F}
done
git commit -m "added some files"

# make 25 files in git history with staged changes
for ((i=0;i<k;i++)); do
  F=$(mktemp foo.XXXX)
  echo "X" >> ${F}
  git add ${F}
  git commit -m "added a file" ${F}
  echo "Y" >> ${F}
  git add ${F}
done

# make 25 files in git history with unstaged changes
for ((i=0;i<k;i++)); do
  F=$(mktemp foo.XXXX)
  echo "X" >> ${F}
  git add ${F}
  git commit -m "added a file" ${F}
  echo "Y" >> ${F}
done

# make 25 new files to be added
for ((i=0;i<k;i++)); do
  F=$(mktemp foo.XXXX)
  echo "X" >> ${F}
  git add ${F}
done

# make 25 untracked files
for ((i=0;i<k;i++)); do
  F=$(mktemp foo.XXXX)
  echo "X" >> ${F}
done
