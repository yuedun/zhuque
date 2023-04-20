#!/bin/bash

sh build.sh && ./zhuque $@ 2>&1|tee zhuque.log 