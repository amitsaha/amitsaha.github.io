Title: Using Travis CI to publish to GitHub pages with custom domain
Date: 2018-01-02 22:00
Category: software
Status: draft

As of today, this blog is automatically being published via [Travis CI](travis-ci.org). 
When I push a new commit to my [GitHub repository](https://github.com/amitsaha/amitsaha.github.io/)
it triggers a new build in Travis CI. The build completes and the the git repository is then
updated with the generated output (mostly HTML with some static CSS). This is how I set it all
up.

## Blog repository setup

I use [pelican](http://docs.getpelican.com/en/stable/) as my blog engine. The "source" code for my
blog lives at the [amitsaha.github.io](https://github.com/amitsaha/amitsaha.github.io/)
repository's `site` branch. Besides the content (markdown and restructured text files) and
pelican specific files, the important files related to publishing are:

- `Dockerfile`
- `.travis.yml`
- `Makefile`

The `Makefile` has a number of targets, but only the `build` target is currently being used:

```
build:
	$(PELICAN) $(INPUTDIR) -o $(OUTPUTDIR) -s $(PUBLISHCONF) $(PELICANOPTS)
	cp 404.md $(OUTPUTDIR)/
```

The first commane generates the site and places the generated files in the `output` sub-directory. In addition
we also copy the `404.md` file to the `output` directory to serve a 
[custom 404](https://help.github.com/articles/creating-a-custom-404-page-for-your-github-pages-site/) page.

The contents of the `output` sub-directory is what we copy to the `master` branch.

To summarize, my blog has two branches:

- `site`: "source" for the blog and other files necessary for generating the HTML for the blog
- `master`: The generated HTML files live in this branch

The generation step is done via Travis and the generated files are pushed to the `master` branch.

Before we move onto setting up Travis


## Generating the blog




