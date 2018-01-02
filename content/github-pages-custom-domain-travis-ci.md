Title: Using Travis CI to publish to GitHub pages with custom domain
Date: 2018-01-02 22:00
Category: software
Status: draft

As of today, this blog is automatically being published via [Travis CI](travis-ci.org). 
When I push a new commit to my [GitHub repository](https://github.com/amitsaha/amitsaha.github.io/)
it triggers a new build in Travis CI. The build completes and the the git repository is then
updated with the generated output (mostly HTML with some static CSS). This is how I set it all
up.

**My setup**

I use [pelican](http://docs.getpelican.com/en/stable/) as my blog engine. The "source" code for my
blog lives at the [amitsaha.github.io](https://github.com/amitsaha/amitsaha.github.io/)
repository's `site` branch. Besides the content (markdown and restructured text files) and
pelican specific files, the important files related to publishing are:

- `Dockerfile`
- `.travis.yml`
- `Makefile`
