from pelican import signals
from pelican.generators import Generator
from pelican.readers import MarkdownReader
from functools import partial

import os
import logging

logger = logging.getLogger(__name__)


class NotesGenerator(Generator):
    """Generate notes sub-directory for articles which are notes"""

    def __init__(self, *args, **kwargs):
        super(NotesGenerator, self).__init__(*args, **kwargs)

    def generate_context(self):
        pass

    def generate_output(self, writer=None):

        # we don't use the writer passed as argument here
        # since we write our own files
        write = partial(writer.write_file,
                        relative_urls=self.settings['RELATIVE_URLS'])
        notes_path = os.path.join(self.output_path, 'notes')
        if not os.path.exists(notes_path):
            try:
                os.mkdir(notes_path)
            except OSError:
                logger.error("Couldn't create the notes output folder in " +
                             notes_path)

        mdreader = MarkdownReader(self.settings)
        for article in self.context['articles']:
            text, meta = mdreader.read(article.source_path)
            if meta['status'] == 'Note':
                write(os.path.join(notes_path, article.save_as), self.get_template(article.template),
                      self.context, article=article, category=article.category,
                      override_output=hasattr(article, 'override_save_as'),
                      blog=True)


def get_generators(generators):
    return NotesGenerator


def register():
    signals.get_generators.connect(get_generators)
