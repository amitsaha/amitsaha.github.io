from pelican import signals
from pelican.generators import CachingGenerator
from pelican.contents import Content
from pelican.utils import process_translations
from itertools import chain

import os
import logging

logger = logging.getLogger(__name__)


class Note(Content):
    mandatory_properties = ('title',)
    allowed_statuses = ('published')
    default_status = 'published'
    default_template = 'page'

    def is_valid(self):
        # TODO
        return True


class NotesGenerator(CachingGenerator):
    """Generate notes sub-directory for articles which are notes"""

    def __init__(self, *args, **kwargs):
        self.notes = []
        self.hidden_translations = []
        super(NotesGenerator, self).__init__(*args, **kwargs)
        # XXX TODO
        signals.page_generator_init.send(self)

    def generate_context(self):
        all_notes = []
        for f in self.get_files(
                self.settings['NOTE_PATHS']):
            note = self.get_cached_data(f, None)
            if note is None:
                try:
                    note = self.readers.read_file(
                        base_path=self.path, path=f, content_class=Note,
                        context=self.context,
                        preread_signal=signals.page_generator_preread,
                        preread_sender=self,
                        context_signal=signals.page_generator_context,
                        context_sender=self)
                except Exception as e:
                    logger.error(
                        'Could not process %s\n%s', f, e,
                        exc_info=self.settings.get('DEBUG', False))
                    self._add_failed_source_path(f)
                    continue

                if not note.is_valid():
                    self._add_failed_source_path(f)
                    continue

                self.cache_data(f, note)

            if note.status == "published":
                all_notes.append(note)
            self.add_source_path(note)

        self.notes, self.translations = process_translations(
            all_notes,
            order_by=self.settings['PAGE_ORDER_BY'])

        self._update_context(('notes',))

        self.save_cache()
        self.readers.save_cache()
        signals.page_generator_finalized.send(self)

    def generate_output(self, writer):
        notes_path = os.path.join(self.output_path, 'notes')
        if not os.path.exists(notes_path):
            try:
                os.mkdir(notes_path)
            except OSError:
                logger.error("Couldn't create the notes output folder in " +
                             notes_path)

        for note in chain(self.translations, self.notes):
            writer.write_file(
                note.save_as, self.get_template(note.template),
                self.context, page=note,
                relative_urls=self.settings['RELATIVE_URLS'],
                override_output=hasattr(note, 'override_save_as'))
            writer.write_file('notes/notes.html', self.get_template('notes_home_template'),
                             self.context, page='<a href="{0}">note.title</a>'.format(note.save_as),
                             relative_urls=self.settings['RELATIVE_URLS'],
                             override_output=hasattr(note, 'override_save_as'),
                            )
        signals.page_writer_finalized.send(self, writer=writer)

def get_generators(generators):
    return NotesGenerator


def register():
    signals.get_generators.connect(get_generators)
