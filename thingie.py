#!/usr/bin/env python

import Adafruit_CharLCD as LCD
import time, os, signal, requests, sys
from collections import OrderedDict

lcd = LCD.Adafruit_CharLCDPlate()
exchange_api = 'https://api.bitcoinaverage.com/exchanges/USD'

def signal_handler(signum, frame):
    lcd.clear()
    sys.exit(0)


signal.signal(signal.SIGINT, signal_handler)

while True:
    r = requests.get(exchange_api)
    prices = r.json()
    prices.pop('timestamp')
    sortedprices = sorted(prices.items(), key=lambda t: t[1]['volume_percent'], reverse=True)
    for exchange_index in range(0,6):
        if exchange_index<len(sortedprices):
            print sortedprices[exchange_index]
            lcd.clear()
            lcd.message("{:<8.8} ${:.2f}\n${:.2f}/{:.2f}".format(
                sortedprices[exchange_index][1]['display_name'][:8],
                sortedprices[exchange_index][1]['rates']['last'],
                sortedprices[exchange_index][1]['rates']['bid'],
                sortedprices[exchange_index][1]['rates']['ask']))
        time.sleep(10)
