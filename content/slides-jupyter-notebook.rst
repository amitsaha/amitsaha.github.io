:Title: Presentation slides with Jupyter Notebook
:Date: 2016-05-31 11:00
:Category: Python


I presented at the PyCon 2016 Education Summit on "Doing Math with Python" day before yesterday and a lightning talk yesterday. This is the first time, I prepared a `slide deck <doingmathwithpython.github.io/pycon-us-2016>`__ using Jupyter Notebook + Reveal.js. I was pleased with the content creation process and the end result. So, here is what worked for me.

Basically, there are two steps:

- Create content in a Jupyter Notebook
- Use ``jupyter-nbconvert`` to convert the notebook to a HTML/JS based slide deck

Finally, I hosted the slide deck on GitHub pages.

Please note that, you will need the latest version of ``jupyter-notebook`` (4.2.0).

Quick start to publish your slides using GitHub pages
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

I have copied my *real* slides to a repo along with a static copy of "reveal.js" (3.3.0) and a Bash script to automate it all. A good starting point then would  be to `fork <https://github.com/amitsaha/jupyter-notebook-slides#fork-destination-box>`__ my GitHub repo:

.. code::

  $ git clone git@github.com:<your-user-name>/jupyter-notebook-slides.git
  $ cd jupyter-notebook-slides
  $ # Start notebook server
  $ # Edit your slides

  $ ./publish.sh

Your slides should now be live at http://<user-github-username>.github.io/jupyter-notebook-slides/

publish.sh script
~~~~~~~~~~~~~~~~~

The ``publish.sh`` script is as follows:

.. code::

   #!/bin/bash

   $ jupyter-nbconvert --to slides slides.ipynb --reveal-prefix=reveal.js
   $ mv slides.slides.html  index.html
   $ mkdir -p /tmp/workspace
   $ cp -r * /tmp/workspace/
   $ git add -A .
   $ git commit -m "Update"
   $ git checkout -B gh-pages
   $ cp -r /tmp/workspace/* .
   $ git add -A .
   $ git commit -m "new version"
   $ git push origin master gh-pages
   $ git checkout master
   $ rm -rf /tmp/workspace


The part of the script above where I update the `gh-pages` branch is of a quality where it does its job.

Viewing your slides locally
~~~~~~~~~~~~~~~~~~~~~~~~~~~

To view your slides locally:

.. code::

   $ jupyter-nbconvert --to slides slides.ipynb --reveal-prefix=reveal.js --post serve

Thanks for your work, Damián Avila and everyone else!
