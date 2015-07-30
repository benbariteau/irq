import sqlite3
import time


def read_quotes_file(path):
    with open(path) as f:
        raw_text = f.read()
    return [
        unicode(quote.strip())
        for quote in raw_text.decode("utf-8").split('\n%') if quote.strip()
    ]


def create_table(cursor, path):
    with open(path) as f:
        schema = f.read()
    cursor.execute(schema)


def insert_quotes(cursor, quotes, is_nishbot=False, is_offensive=False):
    cursor.executemany("""
        INSERT INTO quote
        (text, score, time_created, is_offensive, is_nishbot)
        VALUES (?, 0, ?, ?, ?);
    """, [
        (quote, int(time.time()), is_offensive, is_nishbot)
        for quote in quotes
    ])


if __name__ == '__main__':
    conn = sqlite3.connect('quotes.db')
    cursor = conn.cursor()

    create_table(cursor, 'schema/quote.sql')

    insert_quotes(
        cursor,
        read_quotes_file('batch/backfill_data/quotes'),
    )
    insert_quotes(
        cursor,
        read_quotes_file('batch/backfill_data/offensive'),
        is_offensive=True,
    )
    insert_quotes(
        cursor,
        read_quotes_file('batch/backfill_data/nishbot'),
        is_nishbot=True,
    )

    conn.commit()

    conn.close()
