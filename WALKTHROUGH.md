# Walkthrough of config and general usage

Presumably the binary has been installed or built.  It can be installed via brew, 
see the top level [README](https://github.com/Maahsome/gitlab-tool/blob/main/README.md) file.  

## Configuration

For our example, we will say that we have a couple of different servers, and top
level gitlab groups that we normally access.  

The gitlab-tool uses `access tokens` to interface with the gitlab api, and so we 
will need to have an ENV variable that contains the access token for each server
we will be configuring

### Our Example Servers

- gitlab.com 
  - my-company-group
- git.example.com
  - my-team-group
  - random-pipeline-jobs

So, the gitlab.com group would be accessed as `https://gitlab.com/my-company-group`
and the git.example.com would be accessed on `https://git.example.com/my-team-group`
and `https://git.example.com/random-pipeline-jobs` respectively. 

### gitlab.com configuration

We will choose to set GITLAB_TOKEN_PUBLIC as our ENV variable for this server.
We will also need a place to start our gitlab directory structure on our local
machine.  In our example we will use `~/src/my-company-group`

The command then to configure this server is:

```bash
gitlab-tool config group --gitlab-host gitlab.com --group my-company-group --directory ~/src/my-company-group --env-var GITLAB_TOKEN_PUBLIC
# once added
gitlab-tool config list
```

### git.example.com configuration

Similarly, we will use GITLAB_TOKEN_PRIVATE as our ENV variable, and `~/src/my-team-group` 
and `~/src/random-pipeline-jobs` as our local directories from which to start 
our top level group storage.

The two commands then to configure this server and two top level group are:

```bash
gitlab-tool config group --gitlab-host git.example.com --group my-team-group --directory ~/src/my-team-group --env-var GITLAB_TOKEN_PRIVATE
gitlab-tool config group --gitlab-host git.example.com --group random-pipeline-jobs --directory ~/src/random-pipeline-jobs --env-var GITLAB_TOKEN_PRIVATE
# once added
gitlab-tool config list
```

The list output should look similar to this:

```plaintext
DIRECTORY                                    	ENVVAR   	          GROUP           	    HOST
/Users/<you>/src/my-company-group        	    GITLAB_TOKEN_PUBLIC	  my-company-group      gitlab.com
/Users/<you>/src/my-team-group          	    GITLAB_TOKEN_PRIVATE  my-team-group    	    git.example.com
/Users/<you>/src/random-pipeline-jobs      	    GITLAB_TOKEN_PRIVATE  random-pipeline-jobs  git.example.com
```

## Usage

### Listing and Cloning

The tool is written to associate the directory you are in with the gitlab top 
level group, and group/project (repo) structure below.  So we will start with a 
simple list/ls function

```bash
gitlab-tool ls
```

The output looks like:

```plaintext
STAT 	PATH
gc---	provisioning
ge---	tf-modules
pe--0	bootstrap
pc--0	tools
```

The stat block is defined as such:

```plaintext
-     : group or project (g/p)
 -    : exists on disk (e exists, c created)
  -   : dirty (- clean, + dirty)
   -  : has worktrees (-/w)
    - : number of MRs (for project/repo)
```

For groups, any that do not already have a local directory that matches the name
will be created, thus maintaining a one to one relationship between the directory
structure and the group structure.

For projects/repo the `c` represents that there is no local directory that matches
the name, and can be cloned.  An `e` indicates that the project/repo has a local
directory.

In the example output above, we can clone the `tools` project/repo.

```bash
gitlab-tool clone
```

Running the above will prompt with a list of projects/repos that haven't been
cloned, selecting one will clone the project/repo.

### Getting MRs

If there are Open MRs, we can list them with:

```bash
gitlab-tool get mr
```

```plaintext
IID	TITLE             	STATE 	AUTHOR       	CREATED            	DIFF
15 	Added CHANGELOG.md	opened	You          	2022-07-19 16:29:03	<bash:gitlab-tool get diff -p 2112 -m 15>
```

The README contains information on how to setup Smart Selection in (iterm2) to allow
COMMAND-clicking on the <bash:...> links.  I should probably add a config to turn 
this field off.

You can get the diff of the MR with the following command:

```bash
gitlab-tool get mr -m 15
```

Since we are IN the directory that is the project/repo, the project ID, represented
by the `-p 2112` is automatically determined.

There is currently a hard-coded reference to `/usr/local/bin/delta` to produce the
MR diff.  TODO: add a config to define the diff tool to be called.

### Pipelines

We can list the pipelines for the project.

```bash
gitlab-tool get pipeline
```

```plaintext
ID     	PROJECT ID	STATUS 	JOBS
1470043	6853      	success	<bash:gitlab-tool get jobs -p 2112 -l 1470043>
1470041	6853      	success	<bash:gitlab-tool get jobs -p 2112 -l 1470041>
1469972	6853      	success	<bash:gitlab-tool get jobs -p 2112 -l 1469972>
```

And then the jobs for a specific pipeline

```bash
gitlab-tool get jobs -l 1470043
```

```plaintext
ID        	STATUS 	NAME         	        TRACE
26501120	success	apply_dev_region_east4	<bash:gitlab-tool get trace -p 2112 -j 26501120>
```

And of course, the best part, getting the LOG without having to open a browser window...

```bash
gitlab-tool get trace -j 26501120
```

### CICD Variables

If you do any pipeline work in gitlab, and if you have CICD variables defined at 
group levels, you understand the annoyance of trying to track down where they are
and pulling them in to do some local testing.

We can list the variables, starting at the project, and working back up to the 
top level group.

```bash
gl get variables -o json | jq -r '.[] | [.key,.source] | @tsv' | column -t
```

```plaintext
VARIABLE_ONE                             tools
TF_VAR_terraform_var                     tools
VARIABLE_TWO                             my-company-group
```

This one is quite MacOS specific, because it uses `pbcopy` to capture the output.
It is also specific to terraform, reading varibles that start with TF_VAR and 
setting `export` commands.  These auto-feed terraform variables and is a common
technique for gitlab jobs that run terraform.

This basically outputs `export TF_VAR_terraform_var="my secret value"`, and since 
we want to run these, we capture it using `pbcopy` to put it on the clipboard,
then we can paste it back into the shell and set the ENV varialbes.

```bash
gl get variables | grep TF_VAR | awk '{ print "export "$1"=\x27"$2$3$4$5$6$7"\x27"}' | pbcopy
```

### TODO: add more content!

More content to come

