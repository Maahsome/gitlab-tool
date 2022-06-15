# Pipeline Documentation

## Build-ChangeLog

This pipeline runs on a push/merge to main, it finds the last RELEASE, uses that to find MRs since the last TAG (release) and parses all of the MR `descriptions` for changelog info, and builds a v0.0.0-NEXT.md file, which it commits back to `changelog/v0.0.0-NEXT.md`.

## TAGGED Release Pipeline

A new release is triggered by tagging the repository.  The pipeline runs a goreleaser process, using v0.0.0-NEXT.md as the ChangeLog file, modifying the header line to match the TAG version.  It then renames and commits back the v0.0.0-NEXT.md as `<TAG>.md` in the `changelog` direcory.

vRELEASE.md is used to pass along the ChangeLog from Job to Job, it needs to remain in the `.gitignore` file so `goreleaser` will run, as it does not tolerate dirty git repositories.

The TAG is passed along in line 1 of `vRELEASE.md`
