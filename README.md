# bytes
для работы с битами

Пакет для кодирования и декодирования слов.

Основная идея:
 - мы берем поступающее слово, различных классов (суммы, обычные слова, аббревиатуры), приводим к универсальному виду для однозначной идентификации слова
 - формируем контрольные хэши слова, для лучшей обработки группы символов
 - фильтруем различные опечатки и пересортицу символов кириллицы и латиницы