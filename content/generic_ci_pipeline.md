Title: Creating generic CI pipelines in BuildKite (and elsewhere?)
Date: 2018-03-13
Category: infrastructure
Status: Draft

BuildKite allows creation of [dynamic pipelines](https://buildkite.com/docs/pipelines/defining-steps#dynamic-pipelines)
which basically means that at build time, based on the information contained in the build pipeline or the project
you are building, you can generate the pipeline steps to be relevant only to the project you are building. In this post,
I will describe one such strategy where I setup dynamic pipelines using environment variables specified for a
project's CI pipeline. I then combined it with some bash scripts with the final result being that a project's
repository doesn't need to have any information at all regarding the CI setup. This to me means, you don't have to
maintain CI related files across your repositories - no duplicate files/commands.