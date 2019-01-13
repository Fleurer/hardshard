#!/usr/bin/env python

import sys
from os import path

sys.path.insert(0, path.expanduser('~/night/pymysql'))

import pymysql

connection = pymysql.connect(
    host='localhost',
    port=4001,
    user='uuuuu',
    password='passwd',
    db='db233',
    charset='utf8mb4',
    cursorclass=pymysql.cursors.DictCursor)
